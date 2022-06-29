package models

type CatagoryModel struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
	Link string `db:"link"`
}
