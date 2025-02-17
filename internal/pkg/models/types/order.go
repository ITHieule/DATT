package types

import "time"

type Order struct {
	ID             int           `json:"id"`
	UserID         int           `json:"user_id"`
	RecipientName  string        `json:"recipient_name"`
	RecipientPhone string        `json:"recipient_phone"`
	TotalPrice     float64       `json:"total_price"`
	Status         string        `json:"status"`
	CreatedAt      time.Time     `json:"created_at"`
	OrderDetails   []OrderDetail `json:"order_details" gorm:"foreignKey:OrderID;references:ID"` 
}

type OrderDetail struct {
	ID               int            `json:"id"`
	OrderID          int            `json:"order_id"`
	ProductVariantID int            `json:"product_variant_id"`
	Quantity         int            `json:"quantity"`
	UnitPrice        float64        `json:"unit_price"`
	TotalPrice       float64        `json:"total_price"`
	ProductVariant   ProductVariant `json:"product_variant" gorm:"foreignKey:ProductVariantID;references:ID"`
}
