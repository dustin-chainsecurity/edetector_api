package searchEvidence

import (
    "context"
	"edetector_API/internal/channel"
	"edetector_API/pkg/logger"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func WebSocket(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        logger.Error(err.Error())
        return
    }

    ctx, cancel := context.WithCancel(context.Background())
    messageChannel := make(chan []byte)

    go readMessages(ctx, conn, cancel, messageChannel)
    go handleSignals(ctx, conn, cancel, messageChannel)
    go heartbeat(ctx, conn, cancel, messageChannel)
    go writeMessages(ctx, conn, cancel, messageChannel)
}

func readMessages(ctx context.Context, conn *websocket.Conn, cancel context.CancelFunc, messageChannel chan []byte) {
    for {
        select {
        case <- ctx.Done():
            return
        default:
            _, p, err := conn.ReadMessage()
            if err != nil {
                logger.Error(err.Error())
                fmt.Println("Error reading message from ws: " + err.Error())
                cancel()
                return
            }

			response := []byte("message: " + string(p) + " received")
			select {
			case <-ctx.Done():
				return
			default: 
                messageChannel <- response
			}
        }
    }
}

func writeMessages(ctx context.Context, conn *websocket.Conn, cancel context.CancelFunc, messageChannel <-chan []byte) {
	for {
		select {
		case <-ctx.Done():
			return
		case message := <-messageChannel:
			err := conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				logger.Error(err.Error())
				fmt.Println("Error sending message to ws: " + err.Error())
				cancel()
				return
			}
		}
	}
}

func handleSignals(ctx context.Context, conn *websocket.Conn, cancel context.CancelFunc, messageChannel chan []byte) {
    for {
        select {
        case <- ctx.Done():
            return
        case signal := <-channel.SignalChannel:
			select {
			case <-ctx.Done():
				return
			default:
                messageChannel <- []byte(signal)
			}
        }
    }
}

func heartbeat(ctx context.Context, conn *websocket.Conn, cancel context.CancelFunc, messageChannel chan []byte) {
    for {
        select {
        case <- ctx.Done():
            return
        default:
            select {
            case <- ctx.Done():
                return
            default:
                messageChannel <- []byte("connection check")
                time.Sleep(30 * time.Second)
            }
        }
    }
}