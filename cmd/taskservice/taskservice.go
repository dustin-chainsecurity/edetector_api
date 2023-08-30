package main

import (
	"edetector_API/internal/taskservice"
)

var version string

func main() {
	taskservice.Start(version)
}