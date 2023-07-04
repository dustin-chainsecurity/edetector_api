package main

import (
	"database/sql"
	"edetector_API/pkg/logger"
	"edetector_API/pkg/mariadb"
	"edetector_API/config"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func main() {

	var err error

	// Load configuration
	if config.LoadConfig() == nil {
		fmt.Println("Error loading config file")
		return
	}

	// Connect to MariaDB
	if err = mariadb.Connect_init(); err != nil {
		logger.Error("Error connecting to mariadb: " + err.Error())
	}
	defer mariadb.DB.Close()

	// Create user table
	if err = createUserTable(); err != nil {
		logger.Error("Failed to create user table: " + err.Error())
		return
	}

	// Create user info table
	if err = createUserInfoTable(); err != nil {
		logger.Error("Failed to create user info table: " + err.Error())
		return
	}
}

func createUserTable() error {
	createTableSQL := `
		CREATE TABLE IF NOT EXISTS user (
			id INT AUTO_INCREMENT PRIMARY KEY,
			username VARCHAR(45) NOT NULL,
			password VARCHAR(45) NOT NULL
		);
	`
	_, err := mariadb.DB.Exec(createTableSQL)
	if err != nil {
		return err
	}
	fmt.Println("user table added")
	return nil
}

func createUserInfoTable() error {
	createTableSQL := `
		CREATE TABLE IF NOT EXISTS user_info (
			id INT NOT NULL,
			token VARCHAR(200) NOT NULL,
			email VARCHAR(45) NOT NULL,
			FOREIGN KEY (id) REFERENCES user(id)
		);
	`
	_, err := mariadb.DB.Exec(createTableSQL)
	if err != nil {
		return err
	}
	fmt.Println("user_info table added")
	return nil
}
