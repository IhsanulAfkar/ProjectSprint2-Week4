package forms

type CreateMerchant struct {
	Name             string `json:"name"`
	MerchantCategory string `json:"merchantCategory"`
	ImageUrl         string `json:"imageUrl"`
	Location         struct {
		Lat float32 `json:"lat"`
		Lon float32 `json:"lon"`
	} `json:"location"`
}
