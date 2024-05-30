package models

var ItemCategory = [5]string{
	"Beverage",
	"Food",
	"Snack",
	"Condiments",
	"Additions",
}

type ItemPrice struct {
	Id    string `json:"id"`
	Price int    `json:"price"`
}
type GetItem struct {
	Id              string `json:"id"`
	Name            string `json:"name"`
	ProductCategory string `json:"productCategory"`
	Price           int    `json:"price"`
	ImageUrl        string `json:"imageUrl"`
	CreatedAt       string `json:"createdAt"`
}