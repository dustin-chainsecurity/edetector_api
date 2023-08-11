package main

import (
	"edetector_API/internal/template"
	"fmt"
)

func main() {
	// test process template to raw template
	t := template.Template{
		Name:        "test",
		Work:        "test",
		KeywordType: "test",
		Keyword:     "test",
		HistoryAndBookmark: template.HistoryAndBookmarkStruct{
			ChromeBrowsingHistory:  true,
			ChromeLoginInfo:        true,
			ChromeBookmarks:        true,
			ChromeDownloadHistory:  true,
			ChromeKeywordSearch:    true,
			EdgeBrowsingHistory:    true,
			EdgeLoginInfo:          true,
			EdgeBookmarks:          true,
			IEBrowsingHistory:      true,
			IELoginInfo:            true,
			FirefoxBrowsingHistory: true,
			FirefoxLoginInfo:       true,
			FirefoxBookmarks:       true,
		},
		CookieAndCache: template.CookieAndCacheStruct{
			ChromeCache:    true,
			ChromeCookies:  true,
			EdgeCache:      true,
			EdgeCookies:    true,
			IECache:        true,
			FirefoxCache:   true,
			FirefoxCookies: true,
		},
		ConnectionHistory: template.ConnectionHistoryStruct{
			NetworkInfo:     true,
			NetworkResource: true,
			WirelessInfo:    true,
		},
		ProcessHistory: template.ProcessHistoryStruct{
			InstalledSoftware:       true,
			DetailedSystemService:   true,
			RemoteDesktopInfo:       true,
			SystemInfo:              true,
			Prefetch:                true,
			ScheduledTask:           true,
			NetworkTraffic:          true,
			DNSInfo:                 true,
			GeneralSystemService:    true,
			BootupProgram:           true,
			Jumplist:                true,
			MUICache:                true,
			MachineCodeHistory:      true,
			ProgramReadWriteHistory: true,
		},
		VanishingHistory: template.VanishingHistoryStruct{
			Process:        true,
			OpenedFile:     true,
			ConnectionInfo: true,
			ARPCache:       true,
		},
		RecentOpening: template.RecentOpeningStruct{
			RelatedShortcut: true,
			UserInfo:        true,
			WindowsActivity: true,
			OpenedPath:      true,
			OpenedFile:      true,
		},
		USBHistory: template.USBHistoryStruct{
			USBInfo:         true,
			SystemLogFile:   false,
			SecurityLogFile: true,
		},
		EmailHistory: template.EmailHistoryStruct{
			EmailPath: false,
			EmailList: true,
		},
	}
	raw, err := template.ToRaw(t)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(raw)
}
