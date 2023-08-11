package main

import (
	"edetector_API/internal/template"
	"edetector_API/pkg/mariadb/query"
	"fmt"
)

func main() {
	// test process raw template code
	raw := query.RawTemplate{
		HistoryAndBookmark: "0100110100111",
		CookieAndCache: "1011001",
		ConnectionHistory: "011",
		ProcessHistory: "11110100110011",
		VanishingHistory: "0011",
		RecentOpening: "01011",
		USBHistory: "011",
		EmailHistory: "01",
	}
	template, err := template.ProcessRawTemplate(raw)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(template)

}