package bootstrap

import (
	"database/sql"
	"github.com/go-redis/redis/v7"
	newrelic "github.com/newrelic/go-agent"
	"github.com/rianekacahya/go-kafka/pkg/goconf"
	"github.com/rianekacahya/go-kafka/pkg/golog"
	"github.com/rianekacahya/go-kafka/pkg/goredis"
	"github.com/rianekacahya/go-kafka/pkg/mysql"
	"github.com/rianekacahya/go-kafka/pkg/telemetry"
	"go.uber.org/zap"
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

func initRedis() *redis.Client {
	client, err := goredis.New(
		goconf.Config().GetString("redis.host"),
		goconf.Config().GetString("redis.password"),
	)

	if err != nil {
		log.Fatalf("got an error while connecting to redis server, error: %s", err)
	}

	return client
}

func initLogger() *zap.Logger {
	logger, err := golog.New()
	if err != nil {
		log.Fatalf("got an error while initialize zap logger, error: %s", err)
	}

	return logger
}
