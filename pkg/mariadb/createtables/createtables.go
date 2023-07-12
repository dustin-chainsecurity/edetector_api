package main

import (
	"database/sql"
	"edetector_API/pkg/mariadb"
	"edetector_API/config"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

var UserSQL = `
CREATE TABLE IF NOT EXISTS user (
	id INT AUTO_INCREMENT PRIMARY KEY,
	username VARCHAR(45) NOT NULL,
	password VARCHAR(45) NOT NULL
);
`

var UserInfoSQL = `
CREATE TABLE IF NOT EXISTS user_info (
	id INT NOT NULL,
	token VARCHAR(200) NOT NULL,
	token_time TIMESTAMP,
	email VARCHAR(45) NOT NULL,
	FOREIGN KEY (id) REFERENCES user(id)
);
`

var GroupInfoSQL = `
CREATE TABLE IF NOT EXISTS group_info (
	group_id INT AUTO_INCREMENT PRIMARY KEY,
	group_name VARCHAR(45) NOT NULL,
	description VARCHAR(1000) NOT NULL,
	range_begin VARCHAR(45) NOT NULL,
	range_end VARCHAR(45) NOT NULL
);
`

var ClientGroupSQL = `
CREATE TABLE IF NOT EXISTS client_group (
	client_id varchar(45) NOT NULL,
	group_id INT NOT NULL,
	PRIMARY KEY (client_id, group_id),
	FOREIGN KEY (client_id) REFERENCES client(client_id),
	FOREIGN KEY (group_id) REFERENCES group_info(group_id)
);
`

var ClientTaskStatusSQL = `
CREATE TABLE IF NOT EXISTS client_task_status (
	client_id varchar(45) NOT NULL,
	scan_schedule varchar(45),
	scan_finish_time TIMESTAMP,
	collect_schedule varchar(45),
	collect_finish_time TIMESTAMP,
	file_schedule varchar(45),
	file_finish_time TIMESTAMP,
	image_finish_time TIMESTAMP,
	FOREIGN KEY (client_id) REFERENCES client(client_id)
);
`

var TaskSQL = `
CREATE TABLE IF NOT EXISTS task (
	task_id varchar(45) NOT NULL,
	client_id varchar(45) NOT NULL,
	type varchar(45) NOT NULL,
	status INT NOT NULL,
	timestamp TIMESTAMP NOT NULL,
	PRIMARY KEY (task_id),
	FOREIGN KEY (client_id) REFERENCES client(client_id)
);
`

func main() {

	var err error

	// Load configuration
	if config.LoadConfig() == nil {
		fmt.Println("Error loading config file")
		return
	}

	// Connect to MariaDB
	if err = mariadb.Connect_init(); err != nil {
		fmt.Println("Error connecting to mariadb: " + err.Error())
	}
	defer mariadb.DB.Close()

	// Create table
	if err = createTable(TaskSQL); err != nil {
		fmt.Println("Failed to create table: " + err.Error())
		return
	}
}

func createTable(stmt string) error {
	_, err := mariadb.DB.Exec(stmt)
	if err != nil {
		return err
	}
	fmt.Println("table added")
	return nil
}