package types

import (
	"time"
)

type Product struct {
	ID            int            `json:"id"`
	CategoryID    int            `json:"category_id"`
	Name          string         `json:"name"`
	Description   string         `json:"description"`
	BasePrice     float64        `json:"base_price"`
	CreatedAt     time.Time      `json:"created_at"`
	ProductImages []ProductImage `json:"product_images" gorm:"foreignKey:ProductID;references:ID"`
}
type Product_image struct {
	ID           int      `json:"id"`
	Name         string   `json:"name"`
	BasePrice    float64  `json:"base_price"`
	ImageURLs    []string `json:"image_urls" gorm:"-"`
	ImageURLsRaw string   `json:"-" gorm:"column:image_urls"`
}

type ProductDetailResponse struct {
	ID          int              `json:"id"`
	Name        string           `json:"name"`
	BasePrice   float64          `json:"base_price"`
	Description string           `json:"description"`
	Variants    []ProductVariant `json:"variants" gorm:"foreignKey:ProductID;references:ID"`
	Images      []string         `json:"images"`
}

type ProductVariant struct {
	ID            int            `json:"id"`
	ProductID     int            `json:"product_id"`
	Product       Product        `json:"product" gorm:"foreignKey:ProductID;references:ID"`
	Size          string         `json:"size"`
	Color         string         `json:"color"`
	Stock         int            `json:"stock"`
	Price         float64        `json:"price"`
	ProductImages []ProductImage `json:"product_images" gorm:"foreignKey:ProductID;references:ProductID"`
}

type ProductImage struct {
	ID        int    `json:"id"`
	ProductID int    `json:"product_id"`
	ImageURL  string `json:"image_url"`
}

type ProductHot struct {
	Product  Product `json:"product"`
	ImageURL string  `json:"image_url"`
}
