package services

import (

	"fmt"

	"web-api/internal/pkg/database"

	"web-api/internal/pkg/models/types"

)

type FashionService struct {
	*BaseService
}

var FashionBusiness = &FashionService{}

// 🔹 Lấy toàn bộ danh sách WaterRecords
func (s *FashionService) FashionSevice() ([]types.FashionBusiness, error) {
	var records []types.FashionBusiness

	// Kết nối database
	db, err := database.FashionBusiness()
	if err != nil {
		return nil, fmt.Errorf("lỗi kết nối database: %w", err)
	}
	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	// Thực hiện truy vấn
	query := `SELECT * FROM WaterRecords`
	err = db.Raw(query).Scan(&records).Error
	if err != nil {
		return nil, fmt.Errorf("lỗi truy vấn dữ liệu: %w", err)
	}

	return records, nil
}