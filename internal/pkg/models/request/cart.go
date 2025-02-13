package request

type CartItem struct {
	ID              uint `gorm:"primaryKey"`
	UserID          int  `json:"user_id" gorm:"not null"`
	ProductVariantID int `json:"product_variant_id" gorm:"not null"`
	Quantity        int  `json:"quantity" gorm:"not null;default:1"`
}

type AddToCartRequest struct {
	UserID          int `json:"user_id" validate:"required"`
	ProductVariantID int `json:"product_variant_id" validate:"required"`
	Quantity        int `json:"quantity" validate:"required,min=1"`
}