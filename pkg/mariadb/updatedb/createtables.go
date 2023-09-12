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
	department VARCHAR(45) NOT NULL,
	email VARCHAR(45) NOT NULL,
	permission VARCHAR(255),
	status TINYINT(1) NOT NULL,
	token VARCHAR(45) NOT NULL,
	token_time TIMESTAMP,
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
	client_id VARCHAR(45) NOT NULL,
	group_id INT NOT NULL,
	PRIMARY KEY (client_id, group_id),
	FOREIGN KEY (client_id) REFERENCES client(client_id) ON DELETE CASCADE,
	FOREIGN KEY (group_id) REFERENCES group_info(group_id)
);
`

var ClientTaskStatusSQL = `
CREATE TABLE IF NOT EXISTS client_task_status (
	client_id VARCHAR(45) NOT NULL,
	scan_schedule VARCHAR(45),
	scan_finish_time TIMESTAMP,
	collect_schedule VARCHAR(45),
	collect_finish_time TIMESTAMP,
	file_schedule VARCHAR(45),
	file_finish_time TIMESTAMP,
	image_finish_time TIMESTAMP,
	FOREIGN KEY (client_id) REFERENCES client(client_id) ON DELETE CASCADE
);
`

var TaskSQL = `
CREATE TABLE IF NOT EXISTS task (
	task_id VARCHAR(45) NOT NULL,
	client_id VARCHAR(45) NOT NULL,
	type VARCHAR(45) NOT NULL,
	status INT NOT NULL,
	progress INT NOT NULL,
	timestamp TIMESTAMP NOT NULL,
	PRIMARY KEY (task_id),
	FOREIGN KEY (client_id) REFERENCES client(client_id)
);
`

var TemplateSQL = `
CREATE TABLE IF NOT EXISTS analysis_template (
	template_id VARCHAR(45) NOT NULL,
	template_name VARCHAR(45) NOT NULL,
	work VARCHAR(45) NOT NULL,
	keyword_type VARCHAR(45),
	keyword VARCHAR(255),
	history_and_bookmark VARCHAR(45) NOT NULL,
	cookie_and_cache VARCHAR(45) NOT NULL,
	connection_history VARCHAR(45) NOT NULL,
	process_history VARCHAR(45) NOT NULL,
	vanishing_history VARCHAR(45) NOT NULL,
	recent_opening VARCHAR(45) NOT NULL,
	usb_history VARCHAR(45) NOT NULL,
	email_history VARCHAR(45) NOT NULL,
	PRIMARY KEY (template_id)
);
`
var WhiteListSQL = `
CREATE TABLE IF NOT EXISTS white_list (
	id INT AUTO_INCREMENT PRIMARY KEY,
	setup_user VARCHAR(45) NOT NULL,
	filename VARCHAR(255) NOT NULL,
	md5 VARCHAR(255) NOT NULL,
	signature VARCHAR(10) NOT NULL,
	path VARCHAR(255) NOT NULL,
	reason VARCHAR(255) NOT NULL
);
`

var BlackListSQL = `
CREATE TABLE IF NOT EXISTS black_list (
	id INT AUTO_INCREMENT PRIMARY KEY,
	setup_user VARCHAR(45) NOT NULL,
	filename VARCHAR(255) NOT NULL,
	md5 VARCHAR(255) NOT NULL,
	signature VARCHAR(10) NOT NULL,
	path VARCHAR(255) NOT NULL,
	reason VARCHAR(255) NOT NULL
);
`

var HackListSQL = `
CREATE TABLE IF NOT EXISTS hack_list (
	id INT AUTO_INCREMENT PRIMARY KEY,
	process_name VARCHAR(255) NOT NULL,
	cmd VARCHAR(255) NOT NULL,
	path VARCHAR(255) NOT NULL,
	adding_point INT NOT NULL
);
`

var LogSQL = `
CREATE TABLE IF NOT EXISTS log (
	id INT AUTO_INCREMENT PRIMARY KEY,
	timestamp TIMESTAMP NOT NULL,
	level VARCHAR(20) NOT NULL,
	service VARCHAR(20) NOT NULL,
	content VARCHAR(255) NOT NULL
);
`

var KeyImageSQL = `
CREATE TABLE IF NOT EXISTS key_image (
	id INT AUTO_INCREMENT PRIMARY KEY,
	type VARCHAR(20) NOT NULL,
	apptype VARCHAR(45) NOT NULL,
	path VARCHAR(255) NOT NULL,
	keyword VARCHAR(255) NOT NULL
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