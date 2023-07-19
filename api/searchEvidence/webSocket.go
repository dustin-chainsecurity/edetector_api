package searchEvidence

import (
    "context"
	"edetector_API/internal/channel"
    "edetector_API/internal/token"
	"edetector_API/pkg/logger"
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
    userId  int
	conn   *websocket.Conn
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        return true
	},
}

var users = make(map[*websocket.Conn]bool)
var mutex sync.Mutex

func WebSocket(global_ctx context.Context, w http.ResponseWriter, r *http.Request) {
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
        conn:   conn,
    }

    go client.readMessage(messageChannel)
    go handleSignal(global_ctx, conn, messageChannel) // update task -> refresh
    go heartbeat(global_ctx, conn, messageChannel)
    go broadcast(global_ctx, conn, messageChannel)
}

func (u *User) readMessage(ch chan WsResponse) {
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
            logger.Error("Error receiving message from ws: " + err.Error())
            break
        }
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
            logger.Info("token incorrect")
            break
        }
        ch <- res
    }
}

func broadcast(ctx context.Context, conn *websocket.Conn, ch <-chan WsResponse) {
	for {
		select {
		case <-ctx.Done():
			return
		case message := <-ch:
            for conn := range users {
                err := conn.WriteJSON(message)
                if err != nil {
                    logger.Error("Broadcast error: " + err.Error())
                }
            }
		}
	}
}

func handleSignal(ctx context.Context, conn *websocket.Conn, ch chan WsResponse) {
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
                ch <- msg
            }
        }
    }
}