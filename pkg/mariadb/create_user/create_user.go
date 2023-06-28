package main

import (
	"edetector_API/pkg/mariadb"
	"edetector_API/pkg/mariadb/method"
	"edetector_API/pkg/logger"
	"fmt"
)

func main() {

	// Connect to MariaDB
	if err := mariadb.Connect_init(); err != nil {
		logger.Error("Error connecting to mariadb: " + err.Error())
	} else {
		fmt.Println("Connected to MariaDB")
	}
	defer mariadb.DB.Close()

	// Create user table
	if err := createUserTable(); err != nil {
		logger.Error("Failed to create user table: " + err.Error())
		return
	} else {
		fmt.Println("user table added")
	}

	// Create user_info table
	if err := createUserInfoTable(); err != nil {
		logger.Error("Failed to create user_info table: " + err.Error())
		return
	} else {
		fmt.Println("user_info table added")
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
	_, err := method.Exec(createTableSQL)
	return err
}

func createUserInfoTable() error {
	createTableSQL := `
		CREATE TABLE IF NOT EXISTS user_info (
			id int FOREIGN KEY REFERENCES user(id),
			token VARCHAR(45) NOT NULL,
			email VARCHAR(45) NOT NULL
		);
	`
	_, err := method.Exec(createTableSQL)
	return err
}