package updatedb

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
	FOREIGN KEY (id) REFERENCES user(id) ON DELETE CASCADE
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
	FOREIGN KEY (client_id) REFERENCES client(client_id) ON DELETE CASCADE,
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
	FOREIGN KEY (client_id) REFERENCES client(client_id) ON DELETE CASCADE
);
`

var TaskSQL = `
CREATE TABLE IF NOT EXISTS task (
	task_id varchar(45) NOT NULL,
	client_id varchar(45) NOT NULL,
	type varchar(45) NOT NULL,
	status INT NOT NULL,
	progress INT NOT NULL,
	timestamp TIMESTAMP NOT NULL,
	PRIMARY KEY (task_id),
	FOREIGN KEY (client_id) REFERENCES client(client_id)
);
`

var TemplateSQL = `
CREATE TABLE IF NOT EXISTS analysis_template (
	template_id varchar(45) NOT NULL,
	template_name varchar(45) NOT NULL,
	work varchar(45) NOT NULL,
	keyword_type varchar(45),
	keyword varchar(255),
	history_and_bookmark varchar(45) NOT NULL,
	cookie_and_cache varchar(45) NOT NULL,
	connection_history varchar(45) NOT NULL,
	process_history varchar(45) NOT NULL,
	vanishing_history varchar(45) NOT NULL,
	recent_opening varchar(45) NOT NULL,
	usb_history varchar(45) NOT NULL,
	email_history varchar(45) NOT NULL,
	PRIMARY KEY (template_id)
);
`

func CreateTables() {

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
	if err = createTable(TemplateSQL); err != nil {
		fmt.Println("Failed to create table: " + err.Error())
		return
	}
}

func createTable(stmt string) error {
	_, err := mariadb.DB.Exec(stmt)
	if err != nil {
		return err
	}
	fmt.Println("success")
	return nil
}