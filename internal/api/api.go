package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"edetector_API/pkg/mariadb"
	"edetector_API/pkg/logger"
	login "edetector_API/internal/login"
)

func API_init() {

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
	router.POST("/member/login", login.Init)


	router.Run(":5000")
}