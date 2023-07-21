package websocket

import (
	"context"
	"edetector_API/internal/channel"
	"edetector_API/internal/token"
	"edetector_API/pkg/logger"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var logTime string = time.Now().Format("2006-01-02 - 15:04:05")

type WsRequest struct {
    Authorization string  `json:"authorization"`
    Message       string  `json:"message"`
}

type WsResponse struct {
    IsSuccess   bool      `json:"isSuccess"`
    DeviceId    []string  `json:"deviceId"`
    Message     string    `json:"message"`
}

type User struct {
    userId   int
	conn    *websocket.Conn
    channel  chan WsResponse
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        return true
	},
}

var users = make(map[*websocket.Conn]bool)
var mutex, broadcastMutex sync.Mutex

func webSocket(global_ctx context.Context, w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        logger.Error(err.Error())
        return
    }
    messageChannel := make(chan WsResponse)

    mutex.Lock()
    users[conn] = true
    mutex.Unlock()
    client := &User{
        conn: conn,
        channel: messageChannel,
    }
    fmt.Println("[WS]  " + logTime + " new client from " + conn.RemoteAddr().String())

    go client.readMessage()
    go handleUpdateTask(global_ctx, conn, client.channel) // update task -> refresh
    go heartbeat(global_ctx, conn, client.channel)
    go broadcast(global_ctx, conn, client.channel)
}

func (u *User) readMessage() {
    defer func() {
        mutex.Lock()
        delete(users, u.conn)
        mutex.Unlock()
        u.conn.Close()
    }()

    for {
        var req WsRequest
        err := u.conn.ReadJSON(&req)
        if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
                fmt.Println("[WS]  " + logTime + " client from " + u.conn.RemoteAddr().String() + " disconnected")
			} else {
                logger.Error("Error receiving message from ws: " + err.Error())
			}
            break
        }
        fmt.Println("[WS]  Request content: ", req)
        res := WsResponse {
            IsSuccess: true,
            DeviceId:  []string{},
            Message:   "authorized",
        }
        t := req.Authorization
        u.userId, err = token.Verify(t)
        if err != nil {
            logger.Error(err.Error())
            break
        }
        if u.userId == -1 { // incorrect token
            break
        }
        u.channel <- res
    }
}

func broadcast(ctx context.Context, conn *websocket.Conn, ch <-chan WsResponse) {
	for {
		select {
		case <-ctx.Done():
			return
		case message := <-ch:
            broadcastMutex.Lock()
            for conn := range users {
                err := conn.WriteJSON(message)
                if err != nil {
                    logger.Error("Broadcast error: " + err.Error())
                }
                fmt.Println("[WS]  " + logTime + " broadcast message to " + conn.RemoteAddr().String())
            }
            broadcastMutex.Unlock()
		}
	}
}

func handleUpdateTask(ctx context.Context, conn *websocket.Conn, ch chan WsResponse) {
    for {
        select {
        case <- ctx.Done():
            return
        case signal := <-channel.SignalChannel:
			select {
			case <-ctx.Done():
				return
			default:
                msg := WsResponse {
                    IsSuccess: true,
                    DeviceId:  signal,
                    Message:   "refresh devices required",
                }
                fmt.Println("[WS]  " + logTime + " update device " + signal[0] + " required")
                ch <- msg
			}
        }
    }
}

func heartbeat(ctx context.Context, conn *websocket.Conn, ch chan WsResponse) {
    for {
        select {
        case <- ctx.Done():
            return
        default:
            select {
            case <- ctx.Done():
                return
            default:
                time.Sleep(30 * time.Second)
                msg := WsResponse {
                    IsSuccess: true,
                    DeviceId:  []string{},
                    Message:   "connection check",
                }
                fmt.Println("[WS]  " + logTime + " connection check")
                ch <- msg
            }
        }
    }
}