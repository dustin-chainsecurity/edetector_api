package main

import (
	"edetector_API/config"
	"fmt"
)

func main() {
	if config.LoadConfig() == nil {
		fmt.Println("Error loading config file")
		return
	}
	fmt.Println(config.Viper.GetString("USER"))
}