package main

import (
	"strings"
	"fmt"
)

func main() {
	partitions := strings.Split("11", ",")
	for _, part := range partitions {
		fmt.Println(part)
	}
}