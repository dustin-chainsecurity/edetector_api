package query

import (
	"edetector_API/pkg/mariadb"
)

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