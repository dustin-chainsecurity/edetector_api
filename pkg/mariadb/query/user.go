package query

import (
	"database/sql"
	"edetector_API/pkg/mariadb"
	"strings"
)

type User struct {
	UserID     int       `json:"userID"`
	Username   string    `json:"username"`
	Password   string    `json:"password"`
	Department string    `json:"department"`
	Email      string    `json:"email"`
	Permission []string  `json:"permission"`
	Status     bool      `json:"status"`
}

func GetUsers() ([]User, error) {
	var users []User
	query := `
	SELECT user.id, username, password, department, email, permission, status 
	FROM user LEFT JOIN user_info ON user.id = user_info.id`
	rows, err := mariadb.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var user User
		var permission string
		err := rows.Scan(&user.UserID, &user.Username, &user.Password, &user.Department, &user.Email, &permission, &user.Status)
		if err != nil {
			return nil, err
		}
		user.Permission = strings.Split(permission, "|")
		users = append(users, user)
	}
	return users, nil
}

func GetUserById(id int) (User, error) {
	var user User
	query := `
	SELECT user.id, username, password, department, email, permission, status 
	FROM user LEFT JOIN user_info ON user.id = user_info.id WHERE user.id = ?`
	var permission string
	err := mariadb.DB.QueryRow(query, id).Scan(&user.UserID, &user.Username, &user.Password, &user.Department, &user.Email, &permission, &user.Status)
	if err != nil {
		return user, err
	}
	user.Permission = strings.Split(permission, "|")
	return user, nil
}

func AddUser(user User) (int, error) {
	query := "INSERT INTO user (username, password) VALUES (?, MD5(?))"
	result, err := mariadb.DB.Exec(query, user.Username, user.Password)
	if err != nil {
		return -1, err
	}
	userId, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}
	query = "INSERT INTO user_info (id, department, email, permission, status) VALUES (?, ?, ?, ?, ?)"
	_, err = mariadb.DB.Exec(query, userId, user.Department, user.Email, strings.Join(user.Permission, "|"), user.Status)
	if err != nil {
		return -1, err
	}
	return int(userId), nil
}

func UpdateUser(user User) error {
	query := "UPDATE user SET username = ?, password = MD5(?) WHERE id = ?"
	_, err := mariadb.DB.Exec(query, user.Username, user.Password, user.UserID)
	if err != nil {
		return err
	}
	query = "UPDATE user_info SET department = ?, email = ?, permission = ?, status = ? WHERE id = ?"
	_, err = mariadb.DB.Exec(query, user.Department, user.Email, strings.Join(user.Permission, "|"), user.Status, user.UserID)
	if err != nil {
		return err
	}
	return nil
}

func DeleteUser(ids []int) error {
	for _, id := range ids {
		query := "DELETE FROM user WHERE id = ?"
		_, err := mariadb.DB.Exec(query, id)
		if err != nil {
			return err
		}
	}
	return nil
}

func CheckUserId(id int) (bool, error) {
	var exist bool
	query := "SELECT EXISTS(SELECT 1 FROM user WHERE id = ?)"
	err := mariadb.DB.QueryRow(query, id).Scan(&exist)
	if err != nil {
		return false, err
	}
	return exist, nil
}

func CheckUsername(username string) (bool, error) {
	var exist bool
	query := "SELECT EXISTS(SELECT 1 FROM user WHERE username = ?)"
	err := mariadb.DB.QueryRow(query, username).Scan(&exist)
	if err != nil {
		return false, err
	}
	return exist, nil
}

func CheckUsernameForUpdate(username string, userID int) (bool, error) {
	var exist bool
	query := "SELECT EXISTS(SELECT 1 FROM user WHERE username = ? AND id != ?)"
	err := mariadb.DB.QueryRow(query, username, userID).Scan(&exist)
	if err != nil {
		return false, err
	}
	return exist, nil
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