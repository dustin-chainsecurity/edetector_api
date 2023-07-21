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
    ctx      context.Context
    cancel   context.CancelFunc
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
    ctx, cancel := context.WithCancel(context.Background())

    mutex.Lock()
    users[conn] = true
    mutex.Unlock()
    user := &User{
        conn:    conn,
        channel: messageChannel,
        ctx:     ctx,
        cancel:  cancel,
    }
    fmt.Println("[WS]  " + time.Now().Format("2006-01-02 - 15:04:05") + " new client from " + conn.RemoteAddr().String())

    go user.readMessage()
    go user.heartbeat()
    go handleUpdateTask(global_ctx, conn, user.channel) // update task -> refresh
    go broadcast(global_ctx, conn, user.channel)
}

func (u *User) readMessage() {
    defer func() {
        mutex.Lock()
        delete(users, u.conn)
        mutex.Unlock()
        u.conn.Close()
        u.cancel()
    }()

    for {
        var req WsRequest
        err := u.conn.ReadJSON(&req)
        if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
                fmt.Println("[WS]  " + time.Now().Format("2006-01-02 - 15:04:05") + " client from " + u.conn.RemoteAddr().String() + " disconnected")
			} else {
                logger.Error("Error receiving message from ws: " + err.Error())
			}
            break
        }
        fmt.Println("[WS]  " + time.Now().Format("2006-01-02 - 15:04:05") + " Request content: ", req)
        // verify request
        t := req.Authorization
        u.userId, err = token.Verify(t)
        if err != nil {
            logger.Error(err.Error())
            break
        }
        if u.userId == -1 { // incorrect token
            break
        }
        // response
        res := WsResponse {
            IsSuccess: true,
            DeviceId:  []string{},
            Message:   "authorized",
        }
        broadcastMutex.Lock()
        err = u.conn.WriteJSON(res)
        if err != nil {
            logger.Error("Write authorized message error: " + err.Error())
            broadcastMutex.Unlock()
            break
        }
        fmt.Println("[WS]  " + time.Now().Format("2006-01-02 - 15:04:05") + " write message to " + u.conn.RemoteAddr().String())
        broadcastMutex.Unlock()
    }
}

func (u *User)heartbeat() {
    for {
        select {
        case <- u.ctx.Done():
            return
        default:
            time.Sleep(30 * time.Second)
            select {
            case <- u.ctx.Done():
                return
            default:
                msg := WsResponse {
                    IsSuccess: true,
                    DeviceId:  []string{},
                    Message:   "connection check",
                }
                broadcastMutex.Lock()
                err := u.conn.WriteJSON(msg)
                if err != nil {
                    logger.Error("Heartbeat error: " + err.Error())
                }
                fmt.Println("[WS]  " + time.Now().Format("2006-01-02 - 15:04:05") + " connection check to " + u.conn.RemoteAddr().String())
                broadcastMutex.Unlock()
            }
        }
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
                fmt.Println("[WS]  " + time.Now().Format("2006-01-02 - 15:04:05") + " broadcast message to " + conn.RemoteAddr().String())
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
                fmt.Println("[WS]  " + time.Now().Format("2006-01-02 - 15:04:05") + " update device " + signal[0] + " required")
                ch <- msg
			}
        }
    }
}