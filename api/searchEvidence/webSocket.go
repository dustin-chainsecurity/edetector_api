package searchEvidence

import (
    "context"
	"edetector_API/internal/channel"
    "edetector_API/internal/token"
	"edetector_API/pkg/logger"
	"fmt"
	"net/http"
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

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        return true
	},
}

func WebSocket(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        logger.Error(err.Error())
        return
    }

    ctx, cancel := context.WithCancel(context.Background())
    messageChannel := make(chan WsResponse)

    go readMessages(ctx, conn, cancel, messageChannel)
    go handleSignals(ctx, conn, cancel, messageChannel)
    go heartbeat(ctx, conn, cancel, messageChannel)
    go writeMessages(ctx, conn, cancel, messageChannel)
}

func readMessages(ctx context.Context, conn *websocket.Conn, cancel context.CancelFunc, messageChannel chan WsResponse) {
    for {
        select {
        case <- ctx.Done():
            return
        default:
            var req WsRequest
            err := conn.ReadJSON(&req)
            if err != nil {
                logger.Error(err.Error())
                fmt.Println("Error reading message from ws: " + err.Error())
                cancel()
                return
            }

            res := WsResponse {
                IsSuccess: false,
                DeviceId:  []string{},
                Message:   "unauthorized",
            }

            t := req.Authorization
            _, err = token.Verify(t)
            if err != nil {
                logger.Error(err.Error())
                cancel()
                return
            }

            res.IsSuccess = true
            res.Message = "message received"
			select {
			case <-ctx.Done():
				return
			default: 
                messageChannel <- res
			}
        }
    }
}

func writeMessages(ctx context.Context, conn *websocket.Conn, cancel context.CancelFunc, messageChannel <-chan WsResponse) {
	for {
		select {
		case <-ctx.Done():
			return
		case message := <-messageChannel:
			err := conn.WriteJSON(message)
			if err != nil {
				logger.Error(err.Error())
				fmt.Println("Error sending message to ws: " + err.Error())
				cancel()
				return
			}
		}
	}
}

func handleSignals(ctx context.Context, conn *websocket.Conn, cancel context.CancelFunc, messageChannel chan WsResponse) {
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
                    DeviceId:  []string{signal},
                    Message:   "refresh devices required",
                }
                messageChannel <- msg
			}
        }
    }
}

func heartbeat(ctx context.Context, conn *websocket.Conn, cancel context.CancelFunc, messageChannel chan WsResponse) {
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
                messageChannel <- msg
            }
        }
    }
}