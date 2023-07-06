package api

import (
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
        // Handle error
        fmt.Println("error hoho")
        return
    }

    for {
        messageType, p, err := conn.ReadMessage()
        if err != nil {
            // Handle error
            fmt.Println("error heyhey")
            break
        }

        // Echo received message back to the client
        err = conn.WriteMessage(messageType, p)
        if err != nil {
            // Handle error
            fmt.Println("error haha")
            break
        }
    }
}