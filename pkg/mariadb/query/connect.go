package query

import (
	"context"
	"edetector_API/pkg/logger"
	"edetector_API/pkg/mariadb"
	"time"
)

func CheckConnection(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			logger.Info("Check connection stopped.")
			return
		default:
			query := `SELECT id FROM user WHERE id = 1;`
			_, err := mariadb.DB.Exec(query)
			if err != nil {
				logger.Error("MariaDB connection lost, reconnecting...")
				err = mariadb.Connect_init()
				if err != nil {
					logger.Error("Failed to connect to MariaDB.")
					return
				}
			}
			time.Sleep(1 * time.Hour)
		}
	}
}