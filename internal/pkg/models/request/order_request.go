package request

type CreateOrderRequest struct {
	UserID          int     `json:"user_id" validate:"required"`
	RecipientName   string  `json:"recipient_name" validate:"required"`
	RecipientPhone  string  `json:"recipient_phone" validate:"required"`
	ShippingAddressID int    `json:"shipping_address_id" validate:"required"`
	TotalPrice      float64 `json:"total_price" validate:"required"`
	Status          string  `json:"status,omitempty"`
}

type AddOrderDetailRequest struct {
	OrderID          int     `json:"order_id" validate:"required"`
	ProductVariantID int     `json:"product_variant_id" validate:"required"`
	Quantity         int     `json:"quantity" validate:"required"`
	UnitPrice        float64 `json:"unit_price" validate:"required"`
}
