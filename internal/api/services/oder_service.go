package services

import (
	"errors"
	"fmt"

	"web-api/internal/pkg/database"

	"web-api/internal/pkg/models/types"

	"gorm.io/gorm"
)

type OderService struct {
	*BaseService
}

var OderServi = &OderService{}

func (s *OderService) GetOderByUserID(userID int) ([]types.Order, error) {
	var order []types.Order

	// Kết nối database
	db, err := database.FashionBusiness()
	if err != nil {
		fmt.Println("Database connection error:", err)
		return nil, err
	}
	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	// Truy vấn SQL lấy giỏ hàng theo user_id
	query := `
		SELECT * FROM orders WHERE user_id = ?
	`

	// Thực hiện truy vấn và ánh xạ kết quả vào struct
	err = db.Raw(query, userID).Scan(&order).Error
	if err != nil {
		fmt.Println("Query execution error:", err)
		return nil, err
	}
	return order, nil
}

func (s *OderService) GetOrderByID(orderID int) (*types.Order, error) {
	var order types.Order

	// Kết nối database
	db, err := database.FashionBusiness()
	if err != nil {
		fmt.Println("Database connection error:", err)
		return nil, fmt.Errorf("database connection error: %w", err)
	}
	dbInstance, _ := db.DB()
	defer func() {
		if err := dbInstance.Close(); err != nil {
			fmt.Println("Error closing database connection:", err)
		}
	}()

	// Preload OrderDetails và ShippingAddress vào query
	err = db.
    Preload("OrderDetails.ProductVariant").
    Preload("OrderDetails.ProductVariant.Product").
    Preload("OrderDetails.ProductVariant.ProductImages"). // Preload bảng ProductImages
    Preload("ShippingAddress").
    Where("id = ?", orderID).
    First(&order).Error



	// Kiểm tra lỗi khi không tìm thấy bản ghi
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("order with ID %d not found", orderID) // Thêm thông báo lỗi chi tiết
		}
		fmt.Println("Query execution error:", err)
		return nil, fmt.Errorf("query execution error: %w", err)
	}

	return &order, nil
}