package services

import (
	"errors"
	"fmt"
	"web-api/internal/pkg/database"
	"web-api/internal/pkg/models/request"
)

type CartService struct {
	*BaseService
}

var Cart = &CartService{}

func (s *CartService) AddToCart(userID, productVariantID, quantity int) error {
	if quantity < 1 {
		return errors.New("quantity must be at least 1")
	}
	// Kết nối database
	db, err := database.FashionBusiness()
	if err != nil {
		return fmt.Errorf("lỗi kết nối database: %w", err)
	}
	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	// Thực hiện truy vấn kiểm tra sản phẩm trong giỏ hàng
	var record request.CartItem
	query := `SELECT * FROM cart WHERE user_id = ? AND product_variant_id = ?`
	err = db.Raw(query, userID, productVariantID).Scan(&record).Error
	if err != nil {
		return fmt.Errorf("lỗi truy vấn dữ liệu: %w", err)
	}

	// Nếu sản phẩm đã có trong giỏ hàng, cập nhật số lượng
	if record.ID != 0 {
		updateQuery := `UPDATE cart SET quantity = quantity + ? WHERE user_id = ? AND product_variant_id = ?`
		if err := db.Exec(updateQuery, quantity, userID, productVariantID).Error; err != nil {
			return fmt.Errorf("lỗi cập nhật giỏ hàng: %w", err)
		}
	} else {
		// Nếu chưa có, thêm mới
		insertQuery := `INSERT INTO cart (user_id, product_variant_id, quantity) VALUES (?, ?, ?)`
		if err := db.Exec(insertQuery, userID, productVariantID, quantity).Error; err != nil {
			return fmt.Errorf("lỗi thêm sản phẩm vào giỏ hàng: %w", err)
		}
	}

	return nil
}

func (s *CartService) UpdateCartQuantity(userID, productVariantID, quantity int) error {
	if quantity < 1 {
		return errors.New("quantity must be at least 1")
	}

	// Kết nối database
	db, err := database.FashionBusiness()
	if err != nil {
		return fmt.Errorf("lỗi kết nối database: %w", err)
	}
	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	// Kiểm tra xem sản phẩm có tồn tại trong giỏ hàng không
	var count int
	checkQuery := `SELECT COUNT(*) FROM cart WHERE user_id = ? AND product_variant_id = ?`
	if err := db.Raw(checkQuery, userID, productVariantID).Scan(&count).Error; err != nil {
		return fmt.Errorf("lỗi truy vấn giỏ hàng: %w", err)
	}

	if count == 0 {
		return errors.New("sản phẩm không tồn tại trong giỏ hàng")
	}

	// Cập nhật số lượng sản phẩm trong giỏ hàng
	updateQuery := `UPDATE cart SET quantity = ? WHERE user_id = ? AND product_variant_id = ?`
	if err := db.Exec(updateQuery, quantity, userID, productVariantID).Error; err != nil {
		return fmt.Errorf("lỗi cập nhật số lượng sản phẩm trong giỏ hàng: %w", err)
	}

	return nil
}

func (s *CartService) RemoveFromCart(userID, productVariantID int) error {
	// Kết nối database
	db, err := database.FashionBusiness()
	if err != nil {
		return fmt.Errorf("lỗi kết nối database: %w", err)
	}
	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	// Kiểm tra xem sản phẩm có tồn tại trong giỏ hàng không
	var count int
	checkQuery := `SELECT COUNT(*) FROM cart WHERE user_id = ? AND product_variant_id = ?`
	if err := db.Raw(checkQuery, userID, productVariantID).Scan(&count).Error; err != nil {
		return fmt.Errorf("lỗi truy vấn giỏ hàng: %w", err)
	}

	if count == 0 {
		return errors.New("sản phẩm không tồn tại trong giỏ hàng")
	}

	// Xóa sản phẩm khỏi giỏ hàng
	deleteQuery := `DELETE FROM cart WHERE user_id = ? AND product_variant_id = ?`
	if err := db.Exec(deleteQuery, userID, productVariantID).Error; err != nil {
		return fmt.Errorf("lỗi xóa sản phẩm khỏi giỏ hàng: %w", err)
	}

	return nil
}
