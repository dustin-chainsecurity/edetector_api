package main

import (
	"edetector_API/api/setting"
	"fmt"
)

func main() {
	_, err := setting.GetSetting()
	if err != nil {
		fmt.Println(err)
		return
	}
}
