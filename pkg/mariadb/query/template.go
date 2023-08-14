package query

import (
	"edetector_API/pkg/mariadb"
	"strings"

	"github.com/google/uuid"
)

func LoadRawTemplate(id string) (RawTemplate, error) {
	var raw RawTemplate
	err := mariadb.DB.QueryRow("SELECT * FROM analysis_template WHERE template_id = ?", id).Scan(
		&raw.ID, 
		&raw.Name, 
		&raw.Work, 
		&raw.KeywordType, 
		&raw.Keyword, 
		&raw.HistoryAndBookmark, 
		&raw.CookieAndCache, 
		&raw.ConnectionHistory, 
		&raw.ProcessHistory, 
		&raw.VanishingHistory, 
		&raw.RecentOpening, 
		&raw.USBHistory, 
		&raw.EmailHistory,
	)
	if err != nil {
		return RawTemplate{}, err
	}
	return raw, nil
}

func LoadAllRawTemplate() ([]RawTemplate, error) {
	var raws []RawTemplate
	rows, err := mariadb.DB.Query("SELECT * FROM analysis_template")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var raw RawTemplate
		err := rows.Scan(
			&raw.ID, 
			&raw.Name, 
			&raw.Work, 
			&raw.KeywordType, 
			&raw.Keyword, 
			&raw.HistoryAndBookmark, 
			&raw.CookieAndCache, 
			&raw.ConnectionHistory, 
			&raw.ProcessHistory, 
			&raw.VanishingHistory, 
			&raw.RecentOpening, 
			&raw.USBHistory, 
			&raw.EmailHistory,
		)
		if err != nil {
			return nil, err
		}
		raws = append(raws, raw)
	}
	return raws, nil
}

type RawTemplate struct {
	ID                 string
	Name               string 
	Work               string
	KeywordType        string
	Keyword            string
	HistoryAndBookmark string  // "0110100110101"
	CookieAndCache     string
	ConnectionHistory  string
	ProcessHistory     string
	VanishingHistory   string
	RecentOpening      string
	USBHistory         string
	EmailHistory       string
}

func AddTemplate(raw RawTemplate) (string, error) {
	template_id := uuid.NewString()
	template_id = strings.ReplaceAll(template_id, "-", "")

	query := `
		INSERT INTO analysis_template (template_id, template_name, work,keyword_type,keyword,history_and_bookmark,cookie_and_cache,connection_history,process_history,vanishing_history,recent_opening,usb_history,email_history) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := mariadb.DB.Exec(query,
		template_id,
		raw.Name,
		raw.Work,
		raw.KeywordType,
		raw.Keyword,
		raw.HistoryAndBookmark,
		raw.CookieAndCache,
		raw.ConnectionHistory,
		raw.ProcessHistory,
		raw.VanishingHistory,
		raw.RecentOpening,
		raw.USBHistory,
		raw.EmailHistory,
	)
	if err != nil {
		return "", err
	}
	return template_id, nil
}

func DeleteTemplate(id string) error {
	query := "DELETE FROM analysis_template WHERE template_id = ?"
	_, err := mariadb.DB.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

func UpdateTemplate(id string, raw RawTemplate) error {
	query := `
		UPDATE analysis_template 
		SET template_name = ?, work = ?, keyword_type = ?, keyword = ?, history_and_bookmark = ?, cookie_and_cache = ?, connection_history = ?, process_history = ?, vanishing_history = ?, recent_opening = ?, usb_history = ?, email_history = ?
		WHERE template_id = ?
	`
	_, err := mariadb.DB.Exec(query,
		raw.Name,
		raw.Work,
		raw.KeywordType,
		raw.Keyword,
		raw.HistoryAndBookmark,
		raw.CookieAndCache,
		raw.ConnectionHistory,
		raw.ProcessHistory,
		raw.VanishingHistory,
		raw.RecentOpening,
		raw.USBHistory,
		raw.EmailHistory,
		id,
	)
	if err != nil {
		return err
	}
	return nil
}

func CheckTemplateID (id string) (bool, error) {
	var exist bool
	query := "SELECT EXISTS (SELECT * FROM analysis_template WHERE template_id = ?)"
	err := mariadb.DB.QueryRow(query, id).Scan(&exist)
	if err != nil {
		return false, err
	}
	return exist, nil
}

func CheckTemplateName (name string) (bool, error) {
	var exist bool
	query := "SELECT EXISTS (SELECT * FROM analysis_template WHERE template_name = ?)"
	err := mariadb.DB.QueryRow(query, name).Scan(&exist)
	if err != nil {
		return false, err
	}
	return exist, nil
}

func CheckTemplateNameForUpdate (id string, name string) (bool, error) {
	var exist bool
	query := "SELECT EXISTS (SELECT * FROM analysis_template WHERE template_name = ? AND template_id != ?)"
	err := mariadb.DB.QueryRow(query, name, id).Scan(&exist)
	if err != nil {
		return false, err
	}
	return exist, nil
}