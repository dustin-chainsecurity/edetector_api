package schedule

import (
	"database/sql"
	"edetector_API/pkg/mariadb/query"
	"strings"

	"golang.org/x/exp/slices"
)

type schedule_columns struct {
	clientID         sql.NullString
	scanSchedule     sql.NullString
	collectSchedule  sql.NullString
	fileSchedule     sql.NullString
}

func loadSchedule() []schedule_columns {
	schedules := []schedule_columns{}
	results := query.LoadStoredSchedule()
	for _, r := range results {
		tmp := schedule_columns {
			clientID:         r[0],
			scanSchedule:     r[1],
			collectSchedule:  r[2],
			fileSchedule:     r[3],
		}
		schedules = append(schedules, tmp)
	}
	return schedules
}

func processSchedule() (map[string][]string, map[string][]string, map[string][]string){
	scan := make(map[string][]string)
	collect := make(map[string][]string)
	file := make(map[string][]string)
	schedules := loadSchedule()
	for _, s := range schedules {
		// scan schedule
		if s.scanSchedule.Valid {
			parts := strings.Split(s.scanSchedule.String, ",")
			for _, t := range parts {
				if !slices.Contains(scan[t], s.clientID.String) {
					scan[t] = append(scan[t], s.clientID.String)
				}
			}
		}
		// collect schedule
		if s.collectSchedule.Valid {
			if !slices.Contains(collect[s.collectSchedule.String], s.clientID.String) {
				collect[s.collectSchedule.String] = append(collect[s.collectSchedule.String], s.clientID.String)
			}
		}
		// collect schedule
		if s.collectSchedule.Valid {
			if !slices.Contains(file[s.fileSchedule.String], s.clientID.String) {
				file[s.fileSchedule.String] = append(file[s.collectSchedule.String], s.clientID.String)
			}
		}
	}
	return scan, collect, file
}