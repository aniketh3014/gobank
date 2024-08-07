package main

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	DeleteAccount(int) error
	UpdateAccount(int) error
	GetAccoountById(int) *Account
}

type PostgresDb struct {
	db *sql.DB
}

func NewPostgresDb() (*PostgresDb, error) {
	connstr := "user=postgres dbname=postgres password=secret sslmode=disable"
	db, err := sql.Open("postgres", connstr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresDb{
		db: db,
	}, nil
}
