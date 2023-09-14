package query

import (
	"edetector_API/pkg/mariadb"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type List struct {
	ID         string  `json:"id"`
	Filename   string  `json:"filename"`
	MD5        string  `json:"md5"`
	Signature  string  `json:"signature"`
	Path 	   string  `json:"path"`
	SetupUser  string  `json:"setupUser"`
	Reason     string  `json:"reason"`
}

type HackList struct {
	ID            string  `json:"id"`
	ProcessName   string  `json:"processName"`
	CMD 		  string  `json:"cmd"`
	Path 		  string  `json:"path"`
	AddingPoint   int     `json:"addingPoint"`
}

func GetList(table string) ([]List, error) {
	var list []List
	query := fmt.Sprintf("SELECT * FROM %s", table)
	rows, err := mariadb.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var l List
		err := rows.Scan(&l.ID, &l.Filename, &l.MD5, &l.Signature, &l.Path, &l.SetupUser, &l.Reason)
		if err != nil {
			return nil, err
		}
		list = append(list, l)
	}
	return list, nil
}

func GetHackList() ([]HackList, error) {
	var list []HackList
	query := "SELECT * FROM hack_list"
	rows, err := mariadb.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var l HackList
		err := rows.Scan(&l.ID, &l.ProcessName, &l.CMD, &l.Path, &l.AddingPoint)
		if err != nil {
			return nil, err
		}
		list = append(list, l)
	}
	return list, nil
}

func AddList(table string, l List) (string, error) {
	id := uuid.NewString()
	path := strings.ReplaceAll(l.Path, "\\", "\\\\")
	query := fmt.Sprintf("INSERT INTO %s (id, filename, md5, signature, path, setup_user, reason) VALUES (?, ?, ?, ?, ?, ?, ?)", table)
	_, err := mariadb.DB.Exec(query, id, l.Filename, l.MD5, l.Signature, path, l.SetupUser, l.Reason)
	if err != nil {
		return "", err
	}
	return id, nil
}

func AddHackList(l HackList) (string, error) {
	id := uuid.NewString()
	cmd := strings.ReplaceAll(l.CMD, "\\", "\\\\")
	path := strings.ReplaceAll(l.Path, "\\", "\\\\")
	query := "INSERT INTO hack_list (id, process_name, cmd, path, adding_point) VALUES (?, ?, ?, ?, ?)"
	_, err := mariadb.DB.Exec(query, id, l.ProcessName, cmd, path, l.AddingPoint)
	if err != nil {
		return "", err
	}
	return id, nil
}

func DeleteList(table string, ids []string) error {
	for _, id := range ids {
		query := fmt.Sprintf("DELETE FROM %s WHERE id = ?", table)
		_, err := mariadb.DB.Exec(query, id)
		if err != nil {
			return err
		}
	}
	return nil
}

func DeleteHackList(ids []string) error{
	for _, id := range ids {
		query := "DELETE FROM hack_list WHERE id = ?"
		_, err := mariadb.DB.Exec(query, id)
		if err != nil {
			return err
		}
	}
	return nil
}