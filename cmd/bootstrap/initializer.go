package bootstrap

import (
	"database/sql"
	"github.com/rianekacahya/go-kafka/pkg/goconf"
	"github.com/rianekacahya/go-kafka/pkg/mysql"
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
