package main

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	log "github.com/sirupsen/logrus"
)

func (app *application) checkConnection() bool {
	count := 5

	for count > 0 {
		if app.db.Ping() == nil {
			log.Info("Connection to Database is OK")
			return true
		}
		count--
		log.WithFields(log.Fields{
			"retry": count,
			"err":   "Can't connect to Database",
		}).Info("Checking connection")
	}
	log.Fatal("Can't connect to Database")
	return false
}

func (app *application) initScheduler() {
	dbConnectionString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", app.ctx.DBUser, app.ctx.DBPassword, app.ctx.DBHost, app.ctx.DBPort, app.ctx.DBName)
	m, err := migrate.New(
		"file://scripts/migrations",
		dbConnectionString)
	if err != nil {
		log.Fatal("new: ", err)
	}
	if err := m.Up(); err != nil {
		log.Info("Start migrations")
		if err.Error() == "no change" {
			log.Info("INITIAL-Database --> No changes in database")
		} else {
			log.Fatal("up: ", err)
		}
	}
	log.Info("Migrations are done")
}