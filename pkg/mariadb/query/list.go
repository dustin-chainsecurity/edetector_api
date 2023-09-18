package query

import (
	"edetector_API/pkg/mariadb"
	"fmt"
)

type List struct {
	ID         int     `json:"id"`
	Filename   string  `json:"filename"`
	MD5        string  `json:"md5"`
	Signature  string  `json:"signature"`
	Path 	   string  `json:"path"`
	SetupUser  string  `json:"setupUser"`
	Reason     string  `json:"reason"`
}

type HackList struct {
	ID            int     `json:"id"`
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

func AddList(table string, l List) (int, error) {
	query := fmt.Sprintf("INSERT INTO %s (filename, md5, signature, path, setup_user, reason) VALUES (?, ?, ?, ?, ?, ?)", table)
	result, err := mariadb.DB.Exec(query, l.Filename, l.MD5, l.Signature, l.Path, l.SetupUser, l.Reason)
	if err != nil {
		return -1, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}
	return int(id), nil
}

func AddHackList(l HackList) (int, error) {
	query := "INSERT INTO hack_list (process_name, cmd, path, adding_point) VALUES (?, ?, ?, ?)"
	result, err := mariadb.DB.Exec(query, l.ProcessName, l.CMD, l.Path, l.AddingPoint)
	if err != nil {
		return -1, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}
	return int(id), nil
}

func UpdateList(table string, l List) error {
	query := fmt.Sprintf("UPDATE %s SET filename = ?, md5 = ?, signature = ?, path = ?, setup_user = ?, reason = ? WHERE id = ?", table)
	_, err := mariadb.DB.Exec(query, l.Filename, l.MD5, l.Signature, l.Path, l.SetupUser, l.Reason, l.ID)
	if err != nil {
		return err
	}
	return nil
}

func UpdateHackList(l HackList) error {
	query := "UPDATE hack_list SET process_name = ?, cmd = ?, path = ?, adding_point = ? WHERE id = ?"
	_, err := mariadb.DB.Exec(query, l.ProcessName, l.CMD, l.Path, l.AddingPoint, l.ID)
	if err != nil {
		return err
	}
	return nil
}

func DeleteList(table string, ids []int) error {
	for _, id := range ids {
		query := fmt.Sprintf("DELETE FROM %s WHERE id = ?", table)
		_, err := mariadb.DB.Exec(query, id)
		if err != nil {
			return err
		}
	}
	return nil
}

func DeleteHackList(ids []int) error{
	for _, id := range ids {
		query := "DELETE FROM hack_list WHERE id = ?"
		_, err := mariadb.DB.Exec(query, id)
		if err != nil {
			return err
		}
	}
	return nil
}