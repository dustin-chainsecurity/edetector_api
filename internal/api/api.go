package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"edetector_API/pkg/mariadb"
	"edetector_API/config"
	"edetector_API/pkg/logger"
	"edetector_API/internal/login"
)

func API_init() {

	// Init Logger
	logger.InitLogger(config.Viper.GetString("WORKER_LOG_FILE"))
	fmt.Println("logger is enabled please check all out info in log file: ", config.Viper.GetString("WORKER_LOG_FILE"))

	// Connect to MariaDB
	if err := mariadb.Connect_init(); err != nil {
		logger.Error("Error connecting to mariadb: " + err.Error())
	} else {
		fmt.Println("Connected to MariaDB")
	}

	// Routing
	router := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
    router.RedirectFixedPath = true
	router.Use(cors.New(corsConfig))

	// Login
	router.POST("/member/login", login.Handle)


	router.Run(":5000")
}