package main

import (
	"time"
	"fmt"
)

func main() {
	layout := "2006-01-02 15:04:05"
	parsedTimestamp, err := time.Parse(layout, "2023-07-12 08:29:34")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(parsedTimestamp.Unix())
}