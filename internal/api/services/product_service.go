package services

import (
	"fmt"

	"web-api/internal/pkg/database"
	"web-api/internal/pkg/models/types"
)

type ProductsService struct {
	*BaseService
}

var ProductService = &ProductsService{}

func (s *ProductsService) ProductSevice() ([]types.Product, error) {
	var pro []types.Product

	// Kết nối database
	db, err := database.FashionBusiness()
	if err != nil {
		fmt.Println("Database connection error:", err)

		return nil, err
	}
	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	// Truy vấn SQL lấy ngày đặt hàng và tổng số lượng sách đã bán
	query := `
		SELECT * FROM products

	`

	// Thực hiện truy vấn và ánh xạ kết quả vào struct
	err = db.Raw(query).Scan(&pro).Error
	if err != nil {
		fmt.Println("Query execution error:", err)
		return nil, err
	}
	return pro, nil
}
