package main

import (
	"edetector_API/api"
	"edetector_API/pkg/mariadb/updatedb"
)

func main() {
	api.API_init("TEST_LOG_FILE")
	updatedb.CreateTables()
	// updatedb.CreateEvents()
	// updatedb.InsertData()
}