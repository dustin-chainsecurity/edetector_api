package query

import (
	"edetector_API/pkg/mariadb"
	"fmt"
	"strconv"
)

type GroupInfo struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	RangeBegin  string   `json:"rangeBegin"`
	RangeEnd    string   `json:"rangeEnd"`
	Devices     []string `json:"devices"`
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

func UpdateGroup(g GroupInfo) error {
	query := `
		UPDATE group_info
		SET group_name = ?, description = ?, range_begin = ?, range_end = ?
		WHERE group_id = ?;
	`
	_, err := mariadb.DB.Exec(query, g.Name, g.Description, g.RangeBegin, g.RangeEnd, g.ID)
	if err != nil {
		return err
	}
	return nil
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
	g.ID = groupID
	g.Devices, err = LoadGroupDevices(groupID)
	if err != nil {
		return GroupInfo{}, err
	}
	return g, nil
}

func LoadAllGroupInfo() ([]GroupInfo, error) {
	query := `
		SELECT group_id, group_name, description, range_begin, range_end
		FROM group_info;
	`
	rows, err := mariadb.DB.Query(query)
	if err != nil {
		return []GroupInfo{}, err
	}
	defer rows.Close()

	groups := []GroupInfo{}
	for rows.Next() {
		var g GroupInfo
		err := rows.Scan(&g.ID, &g.Name, &g.Description, &g.RangeBegin, &g.RangeEnd)
		if err != nil {
			return []GroupInfo{}, err
		}
		g.Devices, err = LoadGroupDevices(g.ID)
		if err != nil {
			return []GroupInfo{}, err
		}
		groups = append(groups, g)
	}
	return groups, nil
}

func LoadGroupDevices(groupID int) ([]string, error) {
	query := `
		SELECT client_id FROM client_group WHERE group_id = ?;
	`
	rows, err := mariadb.DB.Query(query, groupID)
	if err != nil {
		return []string{}, err
	}
	defer rows.Close()

	devices := []string{}
	for rows.Next() {
		var device string
		err := rows.Scan(&device)
		if err != nil {
			return []string{}, err
		}
		devices = append(devices, device)
	}
	return devices, nil
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

func AddGroupDevices(groups []int, devices []string) error {
	query := `
		INSERT INTO client_group (group_id, client_id)
		VALUES (?, ?);
	`
	for _, group := range groups {
		for _, device := range devices {
			if exist, err := InGroup(group, device); err != nil {
				return err
			} else if exist {
				continue
			}
			_, err := mariadb.DB.Exec(query, group, device)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func RemoveGroupDevices(groups []int, devices []string) error {
	query := `
		DELETE FROM client_group
		WHERE group_id = ? AND client_id = ?;
	`
	for _, group := range groups {
		for _, device := range devices {
			if exist, err := InGroup(group, device); err != nil {
				return err
			} else if exist {
				_, err := mariadb.DB.Exec(query, group, device)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func CheckGroupID(groupID int) (bool, error) {
	query := `
		SELECT EXISTS(SELECT group_id FROM group_info WHERE group_id = ?);
	`
	var exist bool
	err := mariadb.DB.QueryRow(query, groupID).Scan(&exist)
	if err != nil {
		return false, err
	} else if !exist {
		return false, fmt.Errorf("group id " + strconv.Itoa(groupID) + " not exist")
	}
	return exist, nil
}

func CheckGroupName(groupName string) (bool, error) {
	query := `
		SELECT EXISTS(SELECT group_id FROM group_info WHERE group_name = ?);
	`
	var exist bool
	err := mariadb.DB.QueryRow(query, groupName).Scan(&exist)
	if err != nil {
		return false, err
	}
	return exist, nil
}

func CheckGroupNameForUpdate(groupID int, groupName string) (error) {
	query := `
		SELECT EXISTS(SELECT group_id FROM group_info WHERE group_name = ? AND group_id != ?);
	`
	var exist bool
	err := mariadb.DB.QueryRow(query, groupName, groupID).Scan(&exist)
	if err != nil {
		return err
	} else if exist {
		return fmt.Errorf("group name " + groupName + " already exists")
	}
	return nil
}

func CheckAllGroupID(groups []int) error {
	for _, id := range groups {
		_, err := CheckGroupID(id)
		if err != nil {
			return err
		}
	}
	return nil
}

func InGroup(groupID int, deviceID string) (bool, error) {
	query := `
		SELECT EXISTS(SELECT client_id FROM client_group WHERE group_id = ? AND client_id = ?);
	`
	var exist bool
	err := mariadb.DB.QueryRow(query, groupID, deviceID).Scan(&exist)
	if err != nil {
		return false, err
	}
	return exist, nil
}