package main

import (
	"fmt"
)

func main() {
	slice1 := [][]string{{"a", "b"}, {"c", "d"}}
	slice2 := [][]string{{"a", "b"}, {"c", "d"}}
	fmt.Println(append(slice1, slice2...))
}