package services

import (
	"errors"
	"fmt"
	"strings"
	"web-api/internal/pkg/database"
	"web-api/internal/pkg/models/request"
	"web-api/internal/pkg/models/types"
)

type CartService struct {
	*BaseService
}

var Cart = &CartService{}

func (s *CartService) GetCartByUserID(userID int) ([]types.Carttypes, error) {
	var cartMap = make(map[uint]*types.Carttypes) // Dùng map để gom nhóm dữ liệu

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
		SELECT 
		    c.id, c.user_id, c.quantity,
		    p.id AS product_id, p.name, p.base_price, p.description,
		    pv.id AS product_variant_id, pv.size, pv.color, pv.price,
		    GROUP_CONCAT(pi.image_url) AS images
		FROM cart c
		LEFT JOIN product_variants pv ON c.product_variant_id = pv.id
		LEFT JOIN products p ON pv.product_id = p.id
		LEFT JOIN product_images pi ON p.id = pi.product_id
		WHERE c.user_id = ?
		GROUP BY c.id, c.user_id, p.id, p.name, p.base_price, p.description, pv.id, pv.size, pv.color, pv.price
	`

	// Thực hiện truy vấn
	rows, err := db.Raw(query, userID).Rows()
	if err != nil {
		fmt.Println("Query execution error:", err)
		return nil, err
	}
	defer rows.Close()

	// Xử lý dữ liệu từ database
	for rows.Next() {
		var cartID uint
		var userID int
		var quantity int
		var product types.ProductDetailResponse
		var variant types.ProductVariant
		var images string

		err := rows.Scan(
			&cartID, &userID, &quantity,
			&product.ID, &product.Name, &product.BasePrice, &product.Description,
			&variant.ID, &variant.Size, &variant.Color, &variant.Price,
			&images,
		)
		if err != nil {
			fmt.Println("Row scan error:", err)
			return nil, err
		}

		// Chuyển chuỗi ảnh thành mảng
		product.Images = strings.Split(images, ",")
		product.Variants = append(product.Variants, variant)

		// Nếu cartID chưa có trong map, khởi tạo mới
		if _, exists := cartMap[cartID]; !exists {
			cartMap[cartID] = &types.Carttypes{
				ID:                    cartID,
				UserID:                userID,
				Quantity:              quantity,
				ProductDetailResponse: []types.ProductDetailResponse{},
			}
		}

		// Thêm sản phẩm vào giỏ hàng tương ứng
		cartMap[cartID].ProductDetailResponse = append(cartMap[cartID].ProductDetailResponse, product)
	}

	// Chuyển map thành slice để trả về kết quả
	var carts []types.Carttypes
	for _, cart := range cartMap {
		carts = append(carts, *cart)
	}

	return carts, nil
}

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
