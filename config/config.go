package config

import (
	"fmt"

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
	}
	Viper = vp
	return vp
}
