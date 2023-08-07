package main

// import (
// 	"edetector_API/config"
// 	"edetector_API/pkg/elastic"
// 	"fmt"
// 	// "github.com/gin-gonic/gin"
// )

// func main() {
// 	if config.LoadConfig() == nil {
// 		fmt.Println("Error loading config file")
// 		return
// 	}
// 	err := elastic.SetElkClient()
// 	if err != nil {
// 		fmt.Println(err.Error())
// 	}
// 	indices, err := elastic.GetAllIndices()
// 	if err != nil {
// 		fmt.Println(err.Error())
// 	}
// 	fmt.Println(indices)
// }