package main

import (
	"edetector_API/api"
	"edetector_API/pkg/mariadb/updatedb"
)

func main() {
	api.API_init()
	updatedb.CreateTables()
	// updatedb.CreateEvents()
	// updatedb.InsertData()
}