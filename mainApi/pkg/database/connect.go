package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

func Connect(connUrl string) (*sqlx.DB, error) {
	logrus.Debug("Connecting to database")
	sqlxdb, err := sqlx.Open("postgres", connUrl)

	if err != nil {
		return nil, err
	}
	err = sqlxdb.Ping()
	if err != nil {
		return nil, err
	}

	logrus.Info("Connecting to database")
	return sqlxdb, nil
}

func New(dbURL string) (Database, error) {
	con, err := Connect(dbURL)
	if err != nil {
		return nil, err
	}
	d := &database{
		conn: con,
	}
	return d, nil
}

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBname   string
	SSLMode  string
}

func GetUrl(cfg Config) string {
	url := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBname, cfg.Password, cfg.SSLMode)
	return url
}
