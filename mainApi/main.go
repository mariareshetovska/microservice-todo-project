package main

import (
	"log"
	"mainApi/config"
	"mainApi/pkg/database"
	"mainApi/pkg/router"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("error initializing configs: %s", err.Error())
	}

	urlPostgres := database.GetUrl(config.Db)

	db, err := database.New(urlPostgres)
	if err != nil {
		logrus.WithError(err).Fatal("Error verifing database")
	}
	logrus.Info("Database is ready to use")
	router, err := router.NewRouter(db)
	if err != nil {
		logrus.WithError(err).Fatal("Error building router")
	}

	if err := http.ListenAndServe(":"+config.Server.Port, router); err != nil && err != http.ErrServerClosed {
		logrus.WithError(err).Error("Server failed")
	}

}
