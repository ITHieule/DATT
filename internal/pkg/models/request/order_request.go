package request

type CreateOrderRequest struct {
	UserID          int     `json:"user_id" validate:"required"`
	RecipientName   string  `json:"recipient_name" validate:"required"`
	RecipientPhone  string  `json:"recipient_phone" validate:"required"`
	ShippingAddress string  `json:"shipping_address" validate:"required"`
	ShippingCity    string  `json:"shipping_city" validate:"required"`
	ShippingPostal  string  `json:"shipping_postal_code" validate:"required"`
	ShippingCountry string  `json:"shipping_country" validate:"required"`
	TotalPrice      float64 `json:"total_price" validate:"required"`
}

type AddOrderDetailRequest struct {
	OrderID          int     `json:"order_id" validate:"required"`
	ProductVariantID int     `json:"product_variant_id" validate:"required"`
	Quantity         int     `json:"quantity" validate:"required"`
	UnitPrice        float64 `json:"unit_price" validate:"required"`
}
