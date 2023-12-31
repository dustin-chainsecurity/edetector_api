package mariadb

import (
	"database/sql"
	"edetector_API/config"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func Connect_init() error {
	var err error
	dbUser := config.Viper.GetString("MARIADB_USER")
	dbPass := config.Viper.GetString("MARIADB_PASSWORD")
	dbHost := config.Viper.GetString("MARIADB_HOST")
	dbPort := config.Viper.GetInt("MARIADB_PORT")
	dbName := config.Viper.GetString("MARIADB_DATABASE")

	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	DB, err = sql.Open("mysql", connectionString)
	if err != nil {
		return err
	}
	DB.SetConnMaxLifetime(time.Minute * 5)
	return nil
}