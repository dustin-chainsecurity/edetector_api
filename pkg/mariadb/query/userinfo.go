package query

import (
	"database/sql"
	"edetector_API/pkg/mariadb"
)

func CheckUser(username string, password string) (int, error) {
	userId := -1
	query := "SELECT id FROM user WHERE username = ? AND password = MD5(?)"
	err := mariadb.DB.QueryRow(query, username, password).Scan(&userId)
	if err != nil {
		if err == sql.ErrNoRows { // Username not exist or wrong password
			return userId, nil
		} else {
			return userId, err
		}
	}
	return userId, nil
}

func UpdateToken(token string, userId int) error {
	query := "UPDATE user_info SET token = ?, token_time = CURRENT_TIMESTAMP WHERE id = ?"
	_, err := mariadb.DB.Exec(query, token, userId)
	if err != nil {
		return err
	}
	return nil
}