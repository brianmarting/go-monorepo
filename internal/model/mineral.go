package model

type Mineral struct {
	Id     int    `db:"id"`
	Name   string `db:"name"`
	Amount int    `db:"amount"`
}
