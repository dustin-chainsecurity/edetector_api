package api

import (
    "edetector_API/internal/channel"
    "edetector_API/pkg/logger"
    "fmt"
    "net/http"

    "github.com/gorilla/websocket"	
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func webSocket(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        logger.Error(err.Error())
        return
    }

    go readMessages(conn)
    go handleSignals(conn)
}

func readMessages(conn *websocket.Conn) {
    for {
        messageType, p, err := conn.ReadMessage()
        if err != nil {
            logger.Error(err.Error())
            fmt.Println("Error reading message from ws: " + err.Error())
            break
        }
        err = conn.WriteMessage(messageType, []byte("message: " + string(p) + " received"))
        if err != nil {
            logger.Error(err.Error())
            fmt.Println("Error responding from ws: " + err.Error())
            break
        }
    }
}

func handleSignals(conn *websocket.Conn) {
    for {
        signal := <-channel.SignalChannel
        err := conn.WriteMessage(websocket.TextMessage, []byte(signal))
        if err != nil {
            logger.Error(err.Error())
            fmt.Println("Error sending refresh signal: " + err.Error())
            break
        }
    }
}