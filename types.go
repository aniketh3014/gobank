package main

import (
	"math/rand"
	"time"
)

type TransferReq struct {
	ToAcc  int `json:"toAcc"`
	Amount int `json:"amount"`
}

type CreateAccountRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"LastName"`
}

type Account struct {
	Id        int       `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	AccNumber int64     `json:"accNumber"`
	Balance   int64     `json:"balance"`
	CreatedAt time.Time `json:"createdAt"`
}

// constructer for account

func NewAccount(FirstName, LastName string) *Account {
	return &Account{
		FirstName: FirstName,
		LastName:  LastName,
		AccNumber: rand.Int63n(100000000),
		CreatedAt: time.Now().UTC(),
	}
}
