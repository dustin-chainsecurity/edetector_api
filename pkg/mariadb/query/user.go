package query

import (
	"database/sql"
	"edetector_API/pkg/mariadb"
)

func CheckUsername(username string) (bool, error) {
	var exist bool
	query := "SELECT EXISTS(SELECT 1 FROM user WHERE username = ?)"
	err := mariadb.DB.QueryRow(query, username).Scan(&exist)
	if err != nil {
		return false, err
	}
	return exist, nil
}

func AddUser(username string, password string, token string, email string) error {
	query := "INSERT INTO user (username, password) VALUES (?, MD5(?))"
	_, err := mariadb.DB.Exec(query, username, password)
	if err != nil {
		return err
	}
	userID, err := GetUserId(username)
	if err != nil {
		return err
	}
	query = "INSERT INTO user_info (id, token, email) VALUES (?, ?, ?)"
	_, err = mariadb.DB.Exec(query, userID, token, email)
	if err != nil {
		return err
	}
	return nil
}

func CheckPassword(username string, password string) (int, error) {
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

func GetUsername(userId int) (string, error) {
	var username string
	query := "SELECT username FROM user WHERE id = ?"
	err := mariadb.DB.QueryRow(query, userId).Scan(&username)
	if err != nil {
		return "", err
	}
	return username, nil
}

func GetUserId(username string) (int, error) {
	var userId int
	query := "SELECT id FROM user WHERE username = ?"
	err := mariadb.DB.QueryRow(query, username).Scan(&userId)
	if err != nil {
		return -1, err
	}
	return userId, nil
}