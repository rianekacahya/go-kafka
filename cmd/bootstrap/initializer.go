package bootstrap

import (
	"database/sql"
	newrelic "github.com/newrelic/go-agent"
	"github.com/rianekacahya/go-kafka/pkg/goconf"
	"github.com/rianekacahya/go-kafka/pkg/mysql"
	"github.com/rianekacahya/go-kafka/pkg/telemetry"
	"log"
)

func initMysql() *sql.DB {
	db, err := mysql.New(
		goconf.Config().GetString("mysql.dsn"),
		goconf.Config().GetInt("mysql.max_open"),
		goconf.Config().GetInt("mysql.max_idle"),
		goconf.Config().GetInt("mysql.timeout"),
	)
	if err != nil {
		log.Fatalf("got an error while connecting database server, error: %s", err)
	}

	return db
}

func initTelemetry() newrelic.Application {
	app, err := telemetry.New(
		goconf.Config().GetString("newrelic.id"),
		goconf.Config().GetString("newrelic.key"),

	)
	if err != nil {
		log.Fatalf("got an error while initialize telemetry, error: %s", err)
	}

	return app
}
