package types

import "time"

type Product struct {
	ID          int       `json:"id"`
	CategoryID  int       `json:"category_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	BasePrice   float64   `json:"base_price"`
	CreatedAt   time.Time `json:"created_at"`
}


type ProductVariant struct {
    ID        int     `json:"id"`
    ProductID int     `json:"product_id"`
    Size      string  `json:"size"`
    Color     string  `json:"color"`
    Stock     int     `json:"stock"`
    Price     float64 `json:"price"`
}