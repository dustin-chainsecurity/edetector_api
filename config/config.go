package config

import (
	"log"

	"github.com/spf13/viper"
)

var Viper *viper.Viper

// Config is a struct that holds the configuration for the application
func LoadConfig() *viper.Viper {
	vp := viper.New()
	vp.SetConfigFile("./config/app.env")
	vp.SetConfigType("env")
	// vp.AutomaticEnv()
	if err := vp.ReadInConfig(); err != nil {
		log.Printf("Env file not found")
		// for _, env := range os.Environ() {
		// 	log.Printf("ENV: %s", env)
		// }
	} else {
		log.Println("Config file loaded successfully")
	}
	Viper = vp
	// log.Printf("Viper configuration:")
	// for _, key := range Viper.AllKeys() {
	// 	log.Printf("%s: %s", key, Viper.GetString(key))
	// }
	return vp
}
