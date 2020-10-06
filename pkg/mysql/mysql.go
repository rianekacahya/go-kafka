package mysql

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strings"
	"time"
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

func TransformQuery(query string, arg []interface{}) string {
	var param []interface{}
	a, _ := json.Marshal(arg)
	_ = json.Unmarshal(a, &param)

	for i:=0; i < len(param); i++ {
		if param[i] == nil {
			param[i] = "NULL"
		}else{
			param[i] = fmt.Sprintf("'%v'", param[i])
		}
	}

	return fmt.Sprintf(strings.ReplaceAll(query, "?", "%v"), param...)
}