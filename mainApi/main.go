package main

import (
	"mainApi/pkg/database"
	"mainApi/pkg/router"
	"net/http"

	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	if err := godotenv.Load(); err != nil {
		logrus.WithError(err).Fatal("error loading env variables")
	}

	urlPostgres := database.GetUrl(database.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBname:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	db, err := database.New(urlPostgres)
	if err != nil {
		logrus.WithError(err).Fatal("Error verifing database")
	}
	logrus.Info("Database is ready to use")
	router, err := router.NewRouter(db)
	if err != nil {
		logrus.WithError(err).Fatal("Error building router")
	}
	const port = "8080"
	if err := http.ListenAndServe(":8080", router); err != nil && err != http.ErrServerClosed {
		logrus.WithError(err).Error("Server failed")
	}

}
