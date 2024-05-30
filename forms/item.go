package forms

type CreateItem struct {
	Name            string `json:"name"`
	ProductCategory string `json:"productCategory"`
	Price           int    `json:"price"`
	ImageUrl        string `json:"imageUrl"`
}