package models

var MerchantCategory = [6]string{
	"SmallRestaurant",
	"MediumRestaurant",
	"LargeRestaurant",
	"MerchandiseRestaurant",
	"BoothKiosk",
	"ConvenienceStore",
}

type Merchant struct {
	Id               string  `json:"id"`
	Name             string  `json:"name"`
	MerchantCategory string  `db:"merchantCategory" json:"merchantCategory"`
	ImageUrl         string  `db:"imageUrl" json:"imageUrl"`
	Lat              float64 `db:"lat" json:"lat"`
	Lon              float64 `db:"lon" json:"lon"`
	CreatedAt        string  `db:"createdAt" json:"createdAt"`
}

type GetMerchant struct {
	MerchantId       string `json:"merchantId"`
	Name             string `json:"name"`
	MerchantCategory string `json:"merchantCategory"`
	ImageUrl         string `json:"imageUrl"`
	Location         struct {
		Lat float64 `json:"lat"`
		Lon float64 `json:"lon"`
	} `json:"location"`
	CreatedAt string `json:"createdAt"`
}

type NearbyMerchant struct {
	Merchant DistanceMerchant `json:"merchant"`
	Items    []GetItem        `json:"items"`
}

// dev

type DistanceMerchant struct {
	MerchantId       string `json:"merchantId"`
	Name             string `json:"name"`
	MerchantCategory string `json:"merchantCategory"`
	ImageUrl         string `json:"imageUrl"`
	Location         struct {
		Lat float64 `json:"lat"`
		Lon float64 `json:"lon"`
	} `json:"location"`
	CreatedAt string `json:"createdAt"`
}
