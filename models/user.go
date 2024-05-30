package models

type User struct {
	Id        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	CreatedAt string `db:"createdAt" json:"createdAt"`
	UpdatedAt string `db:"updatedAt" json:"updatedAt"`
}

type GetOrderDetail struct {
	MerchantId       string  `db:"merchantId" json:"merchantId"`
	Name             string  `db:"name" json:"name"`
	MerchantCategory string  `db:"merchantCategory" json:"merchantCategory"`
	ImageUrl         string  `db:"imageUrl" json:"imageUrl"`
	Lat              float64 `db:"lat" json:"lat"`
	Lon              float64 `db:"lon" json:"lon"`
	CreatedAt        string  `db:"createdAt" json:"createdAt"`
	Products         string  `db:"products" json:"products"`
}
type OrderDetail struct {
	Merchant struct {
		MerchantId       string `db:"merchantId" json:"merchantId"`
		Name             string `db:"name" json:"name"`
		MerchantCategory string `db:"merchantCategory" json:"merchantCategory"`
		ImageUrl         string `db:"imageUrl" json:"imageUrl"`
		Location         struct {
			Lat float64 `json:"lat"`
			Lon float64 `json:"lon"`
		} `json:"location"`
		CreatedAt string `db:"createdAt" json:"createdAt"`
	} `json:"merchant"`
	Items []GetItemDetail `json:"items"`
}

type GetItemDetail struct {
	GetItem
	Quantity int `json:"quantity"`
}

type Order struct {
	OrderId string        `json:"orderId"`
	Orders  []OrderDetail `json:"orders"`
}