package mariadb

import (
	"database/sql"
	"edetector_API/config"
	"fmt"

	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-mysql-org/go-mysql/replication"
	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func Connect_init() error {
	var err error
	dbUser := config.Viper.GetString("MARIADB_USER")
	dbPass := config.Viper.GetString("MARIADB_PASSWORD")
	dbHost := config.Viper.GetString("MARIADB_HOST")
	dbPort := config.Viper.GetInt("MARIADB_PORT")
	dbName := config.Viper.GetString("MARIADB_DATABASE")

	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	DB, err = sql.Open("mysql", connectionString)
	if err != nil {
		return err
	}
	fmt.Println("MARIADB Connected")
	return nil
}

func Replicate() (*replication.BinlogStreamer, error) {
	cfg := replication.BinlogSyncerConfig{
		ServerID: 100,
		Flavor:   "mariadb",
		Host:     config.Viper.GetString("MARIADB_HOST"),
		Port:     3306,
		User:     config.Viper.GetString("MARIADB_USER"),
		Password: config.Viper.GetString("MARIADB_PASSWORD"),
	}
	syncer := replication.NewBinlogSyncer(cfg)
	streamer, err := syncer.StartSync(mysql.Position{})
	if err != nil {
		return nil, err
	}
	return streamer, err
}