package mysql

import(
	"database/sql"
	"time"
	_ "github.com/go-sql-driver/mysql"
)

func New(dsn string, maxopen, maxidle, timeout int) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(time.Duration(timeout) * time.Second)
	db.SetMaxOpenConns(maxopen)
	db.SetMaxIdleConns(maxidle)

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}