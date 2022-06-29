package models

type BookModel struct {
	ID    int    `db:"id"`
	Title string `db:"title"`

	Catagories []CatagoryModel
}
