package main

import (
	"edetector_API/pkg/mariadb"
	"fmt"
)

var ClearTokenSQL = `
CREATE EVENT clear_tokens
ON SCHEDULE EVERY 1 DAY
DO
  BEGIN
	UPDATE user_info
	SET token = NULL
    WHERE token IS NOT NULL
	AND token_time < NOW() - INTERVAL 2 DAY;
  END
`

var ClearTaskSQL = `
CREATE EVENT clear_tasks
ON SCHEDULE EVERY 1 DAY
DO
  BEGIN
    DELETE FROM task
    WHERE timestamp < NOW() - INTERVAL 7 DAY
	AND status = 3;
  END
`

func CreateEvents() {
	_, err := mariadb.DB.Exec(ClearTaskSQL)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("create event success")
}