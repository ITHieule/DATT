package services

import (
	"fmt"

	"web-api/internal/pkg/database"

	"web-api/internal/pkg/models/types"
)

type OderService struct {
	*BaseService
}

var OderServi = &OderService{}

func (s *OderService) GetCartByUserID(userID int) ([]types.Order, error) {
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
