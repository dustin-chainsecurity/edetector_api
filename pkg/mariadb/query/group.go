package query

import (
	"edetector_API/pkg/mariadb"
)

type GroupInfo struct {
	Name        string
	Description string
	RangeBegin  string
	RangeEnd	string
}

func AddGroup(g GroupInfo) (int, error) {
	query := `
		INSERT INTO group_info (group_name, description, range_begin, range_end)
		VALUES (?, ?, ?, ?);
	`
	_, err := mariadb.DB.Exec(query, g.Name, g.Description, g.RangeBegin, g.RangeEnd)
	if err != nil {
		return -1, err
	}
	query = `
		SELECT group_id FROM group_info WHERE group_name = ?;
	`
	var groupID int
	err = mariadb.DB.QueryRow(query, g.Name).Scan(&groupID)
	if err != nil {
		return -1, err
	}
	return groupID, nil
}

func RemoveGroup(groupID int) error {
	query := `
		DELETE FROM group_info WHERE group_id = ?;
	`
	_, err := mariadb.DB.Exec(query, groupID)
	if err != nil {
		return err
	}
	return nil
}

func LoadGroupInfo(groupID int) (GroupInfo, error) {
	query := `
		SELECT group_name, description, range_begin, range_end
		FROM group_info
		WHERE group_id = ?;
	`
	var g GroupInfo
	err := mariadb.DB.QueryRow(query, groupID).Scan(&g.Name, &g.Description, &g.RangeBegin, &g.RangeEnd)
	if err != nil {
		return GroupInfo{}, err
	}
	return g, nil
}

func LoadGroups(deviceId string) ([]string, error) {
	query := `
		SELECT group_info.group_name
		FROM client_group
		JOIN group_info ON group_info.group_id = client_group.group_id
		WHERE client_id = ?;
	`
	rows, err := mariadb.DB.Query(query, deviceId)
	if err != nil {
		return []string{}, err
	}
	defer rows.Close()
	
	groups := []string{}
	for rows.Next() {
		var groupName string
		err := rows.Scan(&groupName)
		if err != nil {
			return []string{}, err
        }
		groups = append(groups, groupName)
	}
	return groups, nil
}