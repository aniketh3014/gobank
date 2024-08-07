package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	DeleteAccount(int) error
	UpdateAccount(int) error
	GetAccoountById(int) (*Account, error)
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

	pDb := &PostgresDb{
		db: db,
	}

	if err := pDb.init(); err != nil {
		return nil, err
	}

	return pDb, nil
}

func (d *PostgresDb) init() error {
	return d.CreateAccountTable()
}

func (d *PostgresDb) CreateAccountTable() error {
	query := `CREATE TABLE IF NOT EXISTS account(
		id SERIAL PRIMARY KEY,
		first_name VARCHAR(50),
		last_name VARCHAR(50),
		accNumber SERIAL,
		balance SERIAL,
		created_at TIMESTAMP
	)`
	_, err := d.db.Exec(query)
	return err
}

func (d *PostgresDb) CreateAccount(acc *Account) error {
	query := `INSERT INTO account (first_name, last_name, accNumber, balance, created_at) VALUES ($1, $2, $3, $4, $5)`
	res, err := d.db.Query(query, acc.FirstName, acc.LastName, acc.AccNumber, acc.Balance, acc.CreatedAt)

	if err != nil {
		return err
	}
	fmt.Println(res)
	return nil
}
func (d *PostgresDb) DeleteAccount(id int) error {

	_, err := d.db.Query(`DELETE FROM account WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("unable to delete account")
	}
	return nil
}
func (d *PostgresDb) UpdateAccount(id int) error {
	return nil
}
func (d *PostgresDb) GetAccoountById(id int) (*Account, error) {

	data, err := d.db.Query(`SELECT * FROM account WHERE id = $1`, id)
	if err != nil {
		return nil, err
	}

	if data.Next() {
		account := &Account{}
		err := data.Scan(&account.Id, &account.FirstName, &account.LastName, &account.AccNumber, &account.Balance, &account.CreatedAt)
		if err != nil {
			return nil, err
		}
		return account, nil
	}
	return nil, fmt.Errorf("Account not found %d", id)
}
