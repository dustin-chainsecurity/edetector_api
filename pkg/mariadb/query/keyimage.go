package query

import (
	"edetector_API/config"
	"edetector_API/pkg/mariadb"
	"encoding/json"
	"io/ioutil"
	"os"
)

type KeyImage struct {
	Id      int    `json:"id"`
	Type    string `json:"type"`
	Apptype string `json:"apptype"`
	Path    string `json:"path"`
	Keyword string `json:"keyword"`
}

func LoadStaticImageToDB() (error) {
	var keyImages []KeyImage
	file, err := os.Open(config.Viper.GetString("KEYIMAGE_PATH"))
	if err != nil {
		return err
	}
	defer file.Close()
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(bytes), &keyImages)
	if err != nil {
		return err
	}
	for _, keyImage := range keyImages {
		query := "INSERT INTO key_image (type, apptype, path, keyword) VALUES (?, ?, ?, ?)"
		_, err := mariadb.DB.Exec(query, keyImage.Type, keyImage.Apptype, keyImage.Path, keyImage.Keyword)
		if err != nil {
			return err
		}
	}
	return nil
}

func ResetCustomizedImage() error {
	query := "DELETE FROM key_image WHERE type = 'customized'"
	_, err := mariadb.DB.Exec(query)
	if err != nil {
		return err
	}
	var keyImages []KeyImage
	file, err := os.Open(config.Viper.GetString("KEYIMAGE_PATH"))
	if err != nil {
		return err
	}
	defer file.Close()
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(bytes), &keyImages)
	if err != nil {
		return err
	}
	for _, keyImage := range keyImages {
		if (keyImage.Type == "advanced") {
			query := "INSERT INTO key_image (type, apptype, path, keyword) VALUES ('customized', ?, ?, ?)"
			_, err := mariadb.DB.Exec(query, keyImage.Apptype, keyImage.Path, keyImage.Keyword)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func GetKeyImageByType(keyType string) ([]KeyImage, error) {
	var keyImages []KeyImage
	query := "SELECT id, type, apptype, path, keyword FROM key_image WHERE type = ?"
	res, err := mariadb.DB.Query(query, keyType)
	if err != nil {
		return nil, err
	}
	defer res.Close()
	for res.Next() {
		var keyImage KeyImage
		err := res.Scan(&keyImage.Id, &keyImage.Type, &keyImage.Apptype, &keyImage.Path, &keyImage.Keyword)
		if err != nil {
			return nil, err
		}
		keyImages = append(keyImages, keyImage)
	}
	return keyImages, nil
}

func AddKeyImage(keyImage KeyImage) (int, error) {
	query := "INSERT INTO key_image (type, apptype, path, keyword) VALUES (?, ?, ?, ?)"
	result, err := mariadb.DB.Exec(query, keyImage.Type, keyImage.Apptype, keyImage.Path, keyImage.Keyword)
	if err != nil {
		return -1, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}
	return int(id), nil
}

func DeleteKeyImage(ids []int) error {
	query := "DELETE FROM key_image WHERE id = ?"
	for _, id := range ids {
		_, err := mariadb.DB.Exec(query, id)
		if err != nil {
			return err
		}
	}
	return nil
}

func UpdateKeyImage(keyImage KeyImage) error {
	query := "UPDATE key_image SET type = ?, apptype = ?, path = ?, keyword = ? WHERE id = ?"
	_, err := mariadb.DB.Exec(query, keyImage.Type, keyImage.Apptype, keyImage.Path, keyImage.Keyword, keyImage.Id)
	if err != nil {
		return err
	}
	return nil
}

func CheckKeyImage(keyType string, path string, keyword string) (bool, error) {
	query := "SELECT COUNT(*) FROM key_image WHERE type = ? AND path = ? AND keyword = ?"
	var count int
	err := mariadb.DB.QueryRow(query, keyType, path, keyword).Scan(&count)
	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

func CheckImageID(ID int) (bool, error) {
	query := "SELECT EXISTS(SELECT * FROM key_image WHERE id = ?)"
	var exist bool
	err := mariadb.DB.QueryRow(query, ID).Scan(&exist)
	if err != nil {
		return false, err
	}
	return exist, nil
}