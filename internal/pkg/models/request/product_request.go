package request

import "web-api/internal/pkg/models/types"

type CreateProductRequest struct {
	Id          int                     `json:"id" validate:"required"`
	CategoryID  int                     `json:"category_id" validate:"required"`
	Name        string                  `json:"name" validate:"required"`
	Description string                  `json:"description"`
	BasePrice   float64                  `json:"base_price" validate:"required"`
	Status      string                   `json:"status"`
	Created_at  string                   `json:"created_at"`
	Variants    []types.ProductVariant   `json:"variants"` // Danh sách biến thể
	Images      []string                 `json:"images"`   // Danh sách URL ảnh
}
type CreateProductVariantRequest struct {
	Id        int     `json:"id" validate:"required"`
	ProductID int     `json:"product_id" validate:"required"`
	Size      string  `json:"size" validate:"required"`
	Color     string  `json:"color" validate:"required"`
	Stock     int     `json:"stock"`
	Price     float64 `json:"price" validate:"required"`
}
