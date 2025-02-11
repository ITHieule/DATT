package types

type ShippingAddress struct {
	ID          int    `json:"id"`
	UserID      int    `json:"user_id"`
	FullAddress string `json:"full_address"`
	City        string `json:"city"`
	PostalCode  string `json:"postal_code"`
	Country     string `json:"country"`
	IsDefault   bool   `json:"is_default"`
}
