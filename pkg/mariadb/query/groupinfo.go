package query

import (
	"edetector_API/pkg/mariadb"
)

func LoadGroups(deviceId string) ([]string, error) {
	query := "SELECT group_id FROM client_group WHERE client_id = ?"
	rows, err := mariadb.DB.Query(query)
	if err != nil {
		return []string{}, err
	}
	defer rows.Close()
	
	groups := []string{}
	for rows.Next() {
		var groupId string
		err := rows.Scan(&groupId)
		if err != nil {
			return []string{}, err
        }
		groups = append(groups, groupId)
	}
	return groups, nil
}