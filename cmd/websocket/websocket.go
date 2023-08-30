package main

import (
	"edetector_API/internal/websocket"
)

var version string

func main() {
	websocket.Main(version)
}