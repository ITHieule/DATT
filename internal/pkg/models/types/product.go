package types

import (
	"time"
)

type Product struct {
	ID          int       `json:"id"`
	CategoryID  int       `json:"category_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	BasePrice   float64   `json:"base_price"`
	CreatedAt   time.Time `json:"created_at"`
}

type Product_image struct {
	ID           int      `json:"id"`
	Name         string   `json:"name"`
	BasePrice    float64  `json:"base_price"`
	ImageURLs    []string `json:"image_urls" gorm:"-"`
	ImageURLsRaw string   `json:"-" gorm:"column:image_urls"` // Trường tạm để nhận dữ liệu từ SQL
}

type ProductDetailResponse struct {
	ID        int              `json:"id"`
	Name      string           `json:"name"`
	BasePrice float64          `json:"base_price"`
	Variants  []ProductVariant `json:"variants"`
	Images []string `json:"images"`
}

// ProductVariant đại diện cho một biến thể của sản phẩm
type ProductVariant struct {
	ID        int     `json:"id"`
	ProductID int     `json:"product_id"`
	Size      string  `json:"size"`
	Color     string  `json:"color"`
	Stock     int     `json:"stock"`
	Price     float64 `json:"price"`
}
