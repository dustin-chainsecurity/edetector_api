package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

var Viper *viper.Viper

// Config is a struct that holds the configuration for the application
func LoadConfig() *viper.Viper {
	vp := viper.New()
	vp.SetConfigFile("app.env")
	vp.SetConfigType("env")
	vp.AutomaticEnv()
	if err := vp.ReadInConfig(); err != nil {
		fmt.Printf("Env file not found")
	} else {
		fmt.Println("Config file loaded successfully")
	}
	Viper = vp
	log.Printf("Viper configuration:")
	for _, key := range Viper.AllKeys() {
		log.Printf("%s: %s", key, Viper.GetString(key))
	}
	return vp
}
