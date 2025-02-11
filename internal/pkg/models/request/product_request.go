package request

type CreateProductRequest struct {
	CategoryID  int     `json:"category_id" validate:"required"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	BasePrice   float64 `json:"base_price" validate:"required"`
}

type CreateProductVariantRequest struct {
	ProductID int     `json:"product_id" validate:"required"`
	Size      string  `json:"size" validate:"required"`
	Color     string  `json:"color" validate:"required"`
	Stock     int     `json:"stock"`
	Price     float64 `json:"price" validate:"required"`
}
