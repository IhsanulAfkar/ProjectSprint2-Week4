package forms

type EstimatePrice struct {
	UserLocation struct {
		Lat  float64 `json:"lat"`
		Long float64 `json:"long"`
	} `json:"userLocation"`
	Orders []EstimateOrder `json:"orders"`
}

type EstimateOrder struct {
	MerchantId      string              `json:"merchantId"`
	IsStartingPoint bool                `json:"isStartingPoint"`
	Items           []EstimateOrderItem `json:"items"`
}

type EstimateOrderItem struct {
	ItemId   string `json:"itemId"`
	Quantity int    `json:"quantity"`
}
