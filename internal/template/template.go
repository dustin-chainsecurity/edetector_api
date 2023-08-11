package template

import (
	"edetector_API/pkg/mariadb/query"
	"fmt"
	"reflect"
)

var CategoryMap map[string][]string
var Field map[string]string

func init() {

	CategoryMap = map[string][]string {
		"main":				 	{ "history_and_bookmark", "cookie_and_cache", "connection_history", "process_history", "vanishing_history", "recent_opening", "usb_history", "email_history" },
		"history_and_bookmark": { "ChromeBrowsingHistory", "ChromeLoginInfo", "ChromeBookmarks", "ChromeDownloadHistory", "ChromeKeywordSearch", "EdgeBrowsingHistory", "EdgeLoginInfo", "EdgeBookmarks", "IEBrowsingHistory", "IELoginInfo", "FirefoxBrowsingHistory", "FirefoxLoginInfo", "FirefoxBookmarks"},
		"cookie_and_cache":     { "ChromeCache", "ChromeCookies", "EdgeCache", "EdgeCookies", "IECache", "FirefoxCache", "FirefoxCookies" },
		"connection_history":   { "NetworkInfo", "NetworkResource", "WirelessInfo" },
		"process_history":      { "InstalledSoftware", "DetailedSystemService", "RemoteDesktopInfo", "SystemInfo", "Prefetch", "ScheduledTask", "NetworkTraffic", "DNSInfo", "GeneralSystemService", "BootupProgram", "Jumplist", "MUICache", "MachineCodeHistory", "ProgramReadWriteHistory" },
		"vanishing_history":    { "Process", "OpenedFile", "ConnectionInfo", "ARPCache" },
		"recent_opening":       { "RelatedShortcut", "UserInfo", "WindowsActivity", "OpenedPath", "OpenedFile" },
		"usb_history":          { "USBInfo", "SystemLogFile", "SecurityLogFile" },
		"email_history":        { "EmailPath", "EmailList" },
	}

	Field = map[string]string {
		"history_and_bookmark": "HistoryAndBookmark",
		"cookie_and_cache":     "CookieAndCache",
		"connection_history":   "ConnectionHistory",
		"process_history":      "ProcessHistory",
		"vanishing_history":    "VanishingHistory",
		"recent_opening":       "RecentOpening",
		"usb_history":          "USBHistory",
		"email_history":        "EmailHistory",
	}
}

type Template struct {
	ID                  int                  `json:"template_id"`
	Name                string               `json:"template_name"`
	Work                string               `json:"work"`
	KeywordType         string               `json:"keyword_type"`
	Keyword             string               `json:"keyword"`
	HistoryAndBookmark  historyAndBookmark   `json:"history_and_bookmark"`
	CookieAndCache      cookieAndCache       `json:"cookie_and_cache"`
	ConnectionHistory   connectionHistory    `json:"connection_history"`
	ProcessHistory      processHistory       `json:"process_history"`
	VanishingHistory    vanishingHistory     `json:"vanishing_history"`
	RecentOpening       recentOpening        `json:"recent_opening"`
	USBHistory          usbHistory           `json:"usb_history"`
	EmailHistory        emailHistory         `json:"email_history"`
}

type historyAndBookmark struct {
	ChromeBrowsingHistory   bool `json:"chrome_browsing_history"`
	ChromeLoginInfo 	    bool `json:"chrome_login_info"`
	ChromeBookmarks         bool `json:"chrome_bookmark"`
	ChromeDownloadHistory   bool `json:"chrome_download_history"`
	ChromeKeywordSearch     bool `json:"chrome_keyword_search"`
	EdgeBrowsingHistory     bool `json:"edge_browsing_history"`
	EdgeLoginInfo           bool `json:"edge_login_info"`
	EdgeBookmarks           bool `json:"edge_bookmark"`
	IEBrowsingHistory       bool `json:"ie_browsing_history"`
	IELoginInfo             bool `json:"ie_login_info"`
	FirefoxBrowsingHistory  bool `json:"firefox_browsing_history"`
	FirefoxLoginInfo        bool `json:"firefox_login_info"`
	FirefoxBookmarks        bool `json:"firefox_bookmark"`
}

type cookieAndCache struct {
	ChromeCache			 bool `json:"chrome_cache"`
	ChromeCookies		 bool `json:"chrome_cookie"`
	EdgeCache			 bool `json:"edge_cache"`
	EdgeCookies			 bool `json:"edge_cookie"`
	IECache				 bool `json:"ie_cache"`
	FirefoxCache		 bool `json:"firefox_cache"`
	FirefoxCookies		 bool `json:"firefox_cookie"`
}

type connectionHistory struct {
	NetworkInfo			 bool `json:"network_info"`
	NetworkResource		 bool `json:"network_resource"`
	WirelessInfo		 bool `json:"wireless_info"`
}

type processHistory struct {
	InstalledSoftware			bool `json:"installed_software"`
	DetailedSystemService		bool `json:"detailed_system_service"`
	RemoteDesktopInfo			bool `json:"remote_desktop_info"`
	SystemInfo					bool `json:"system_info"`
	Prefetch					bool `json:"prefetch"`
	ScheduledTask				bool `json:"scheduled_task"`
	NetworkTraffic				bool `json:"network_traffic"`
	DNSInfo						bool `json:"dns_info"`
	GeneralSystemService		bool `json:"general_system_service"`
	BootupProgram				bool `json:"bootup_program"`
	Jumplist					bool `json:"jumplist"`
	MUICache					bool `json:"mui_cache"`
	MachineCodeHistory			bool `json:"machine_code_history"`
	ProgramReadWriteHistory		bool `json:"program_read_write_history"`
}

type vanishingHistory struct {
	Process				bool `json:"process"`
	OpenedFile			bool `json:"opened_file"`
	ConnectionInfo		bool `json:"connection_info"`
	ARPCache			bool `json:"arp_cache"`
}

type recentOpening struct {
	RelatedShortcut	bool `json:"related_shortcut"`
	UserInfo		bool `json:"user_info"`
	WindowsActivity	bool `json:"windows_activity"`
	OpenedPath		bool `json:"opened_path"`
	OpenedFile		bool `json:"opened_file"`
}

type usbHistory struct {
	USBInfo				bool `json:"usb_info"`
	SystemLogFile		bool `json:"system_log_file"`
	SecurityLogFile		bool `json:"security_log_file"`
}

type emailHistory struct {
	EmailPath			bool `json:"email_path"`
	EmailList			bool `json:"email_list"`
}

func ProcessRawTemplate(raw query.RawTemplate) (Template, error) {
	var template Template
	for _, category := range CategoryMap["main"] {
		err := processCategory(category, &template, raw)
		if err != nil {
			return template, err
		}
	}
	return template, nil
}

func processCategory(category string, template *Template, raw query.RawTemplate) (err error) {
	categoryValues, found := CategoryMap[category]
	if !found {
		return fmt.Errorf("invalid category map key %s", category)
	}

	field, found := Field[category]
	if !found {
		return fmt.Errorf("invalid field map key %s", category)
	}
	target := reflect.ValueOf(template).Elem().FieldByName(field)
	if !target.IsValid() || !target.CanSet() {
		return fmt.Errorf("invalid field name %s in category %s", field, category)
	}
	rawSlice := reflect.ValueOf(&raw).Elem().FieldByName(field).String()
	targetValue := reflect.ValueOf(target).Elem()
	for i, key := range categoryValues {
		field := targetValue.FieldByName(key)
		if field.IsValid() && field.CanSet() {
			field.SetBool(rawSlice[i] == '1')
		} else {
			return fmt.Errorf("invalid field name %s", key)
		}
	}

	return nil
}