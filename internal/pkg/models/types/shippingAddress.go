package types

type ShippingAddress struct {
	ID         int     `json:"id"`
	UserID     int     `json:"user_id"`
	Province   string  `json:"province"`
	District   string  `json:"district"`
	Ward       string  `json:"ward"`
	PostalCode string  `json:"postal_code"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
	IsDefault  bool    `json:"is_default"`
}
