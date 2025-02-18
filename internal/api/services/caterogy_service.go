package services

import (
	"fmt"

	"web-api/internal/pkg/database"

	"web-api/internal/pkg/models/types"
)

type CaterogyService struct {
	*BaseService
}

var Caterogy = &CaterogyService{}

// 🔹 Lấy toàn bộ danh sách WaterRecords
func (s *CaterogyService) GetCaterogySevice() ([]types.Category, error) {
	var records []types.Category

	// Kết nối database
	db, err := database.FashionBusiness()
	if err != nil {
		return nil, fmt.Errorf("lỗi kết nối database: %w", err)
	}
	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	// Thực hiện truy vấn
	query := `SELECT * FROM categories`
	err = db.Raw(query).Scan(&records).Error
	if err != nil {
		return nil, fmt.Errorf("lỗi truy vấn dữ liệu: %w", err)
	}

	return records, nil
}
