package main

import (
	"edetector_API/config"
	"edetector_API/pkg/logger"
	"edetector_API/pkg/mariadb"
	"edetector_API/pkg/mariadb/updatedb"
	"edetector_API/pkg/redis"
	"fmt"
)

var err error

func updatedb_init(LOG_PATH string, HOSTNAME string, APP string) {
	// Load configuration
	if config.LoadConfig() == nil {
		fmt.Println("Error loading config file")
		return
	}
	// Init Logger
	logger.InitLogger(config.Viper.GetString(LOG_PATH), HOSTNAME, APP)
	logger.Log.Info("Logger enabled, log file: " + config.Viper.GetString(LOG_PATH))
	// Connect to Redis
	if db := redis.Redis_init(); db == nil {
		logger.Error("Error connecting to redis")
		return
	}
	// Connect to MariaDB
	if err = mariadb.Connect_init(); err != nil {
		logger.Error("Error connecting to mariadb: " + err.Error())
		return
	}
}

func main() {
	updatedb_init("TEST_LOG_FILE", "TEST_HOSTNAME", "TEST_APP")
	updatedb.CreateTables()
	// updatedb.CreateEvents()
	// updatedb.InsertData()
}