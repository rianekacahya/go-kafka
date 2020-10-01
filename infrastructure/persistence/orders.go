package persistence

import "database/sql"

type orders struct {
	database *sql.DB
}

func NewOrdersRepository(database *sql.DB) *orders {
	return &orders{database}
}
