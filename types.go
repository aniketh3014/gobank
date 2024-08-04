package main

import "math/rand"

type Account struct {
	Id        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	AccNumber int64  `json:"accNumber"`
	Balance   int64  `json:"balance"`
}

// constructer for account

func NewAccount(FirstName, LastName string) *Account {
	return &Account{
		Id:        rand.Intn(100000),
		FirstName: FirstName,
		LastName:  LastName,
		AccNumber: rand.Int63n(100000000),
	}
}
