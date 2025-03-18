package services

import (
	"database/sql"
	"errors"
	"fmt"

	"web-api/internal/pkg/database"

	"web-api/internal/pkg/models/types"

	"gorm.io/gorm"
)

type OderService struct {
	*BaseService
}

var Order = &OderService{}

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
		return nil, fmt.Errorf("database connection error: %w", err)
	}
	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	// Preload dữ liệu
	err = db.
		Preload("OrderDetails.ProductVariant").
		Preload("OrderDetails.ProductVariant.Product").
		Preload("OrderDetails.ProductVariant.Product.ProductImages").
		Preload("ShippingAddress").
		Where("id = ?", orderID).
		First(&order).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("order with ID %d not found", orderID)
		}
		return nil, fmt.Errorf("query execution error: %w", err)
	}
	for i := range order.OrderDetails {
		order.OrderDetails[i].ProductVariant.ProductImages = order.OrderDetails[i].ProductVariant.Product.ProductImages
	}

	return &order, nil
}

func (s *OderService) GetCartsByUserID(userID int) ([]types.Carttypes, error) {
	var cartMap = make(map[uint]*types.Carttypes)
	var productMap = make(map[int]*types.ProductDetailResponse)

	db, err := database.FashionBusiness()
	if err != nil {
		fmt.Println("Database connection error:", err)
		return nil, err
	}
	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	query := `
		SELECT 
		    c.id AS cart_id, c.user_id, c.quantity,
		    p.id AS product_id, p.name, p.base_price, p.description,
		    pv.id AS product_variant_id, pv.size, pv.color, pv.price,
		    pi.image_url
		FROM cart c
		LEFT JOIN product_variants pv ON c.product_variant_id = pv.id
		LEFT JOIN products p ON pv.product_id = p.id
		LEFT JOIN product_images pi ON p.id = pi.product_id
		WHERE c.user_id = ?
	`

	rows, err := db.Raw(query, userID).Rows()
	if err != nil {
		fmt.Println("Query execution error:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var cartID uint
		var userID int
		var quantity int
		var productID int
		var productName string
		var productBasePrice float64
		var productDescription string
		var variantID int
		var size, color string
		var price float64
		var imageUrl sql.NullString

		err := rows.Scan(
			&cartID, &userID, &quantity,
			&productID, &productName, &productBasePrice, &productDescription,
			&variantID, &size, &color, &price,
			&imageUrl,
		)
		if err != nil {
			fmt.Println("Row scan error:", err)
			return nil, err
		}

		variant := types.ProductVariant{
			ID:    variantID,
			Size:  size,
			Color: color,
			Price: price,
		}

		if _, exists := productMap[productID]; !exists {
			productMap[productID] = &types.ProductDetailResponse{
				ID:          productID,
				Name:        productName,
				BasePrice:   productBasePrice,
				Description: productDescription,
				Variants:    []types.ProductVariant{},
				Images:      []string{},
			}
		}

		productMap[productID].Variants = append(productMap[productID].Variants, variant)

		if imageUrl.Valid {
			productMap[productID].Images = append(productMap[productID].Images, imageUrl.String)
		}

		if _, exists := cartMap[cartID]; !exists {
			cartMap[cartID] = &types.Carttypes{
				ID:                    cartID,
				UserID:                userID,
				Quantity:              quantity,
				ProductDetailResponse: []types.ProductDetailResponse{},
			}
		}
	}

	for _, product := range productMap {
		for _, cart := range cartMap {
			cart.ProductDetailResponse = append(cart.ProductDetailResponse, *product)
			break
		}
	}

	var carts []types.Carttypes
	for _, cart := range cartMap {
		carts = append(carts, *cart)
	}

	return carts, nil
}

func (s *OderService) CreateOrderFromCart(userID int, order *types.Order) error {
	db, err := database.FashionBusiness()
	if err != nil {
		fmt.Println("Database connection error:", err)
		return err
	}

	// ✅ Lấy giỏ hàng thông qua CartService
	carts, err := s.GetCartsByUserID(userID)
	if err != nil {
		return err
	}
	if len(carts) == 0 {
		return errors.New("cart is empty, cannot create order")
	}

	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	var shippingAddress types.ShippingAddress
	if err := db.Model(&types.ShippingAddress{}).
		Where("id = ? AND user_id = ?", order.ShippingAddressID, userID).
		First(&shippingAddress).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			if err := db.Model(&types.ShippingAddress{}).
				Where("user_id = ? AND is_default = ?", userID, true).
				First(&shippingAddress).Error; err != nil {
				return errors.New("Vui Lòng Chọn Địa Chỉ Mặc Định")
			}
			order.ShippingAddressID = shippingAddress.ID
		} else {
			return err
		}
	}

	// 🚀 Bắt đầu transaction
	err = db.Transaction(func(tx *gorm.DB) error {
		totalPrice := 0.0

		// Map để theo dõi tổng số lượng và tổng tiền của từng ProductVariantID
		orderDetailsMap := make(map[int]*types.OrderDetail)

		// Tính tổng tiền từ các giỏ hàng và cập nhật orderDetailsMap
		for _, cart := range carts {
			for _, product := range cart.ProductDetailResponse {
				// Nếu mảng variant rỗng thì bỏ qua sản phẩm này
				if len(product.Variants) == 0 {
					continue
				}
				// Giả sử variant được chọn là variant đầu tiên
				variant := product.Variants[0]

				// Cập nhật orderDetailsMap: nếu đã tồn tại thì cộng dồn, nếu chưa tồn tại thì khởi tạo
				if detail, exists := orderDetailsMap[variant.ID]; exists {
					detail.Quantity += cart.Quantity
					detail.TotalPrice += variant.Price * float64(cart.Quantity)
				} else {
					orderDetailsMap[variant.ID] = &types.OrderDetail{
						ProductVariantID: variant.ID,
						Quantity:         cart.Quantity,
						UnitPrice:        variant.Price,
						TotalPrice:       variant.Price * float64(cart.Quantity),
					}
				}
				totalPrice += variant.Price * float64(cart.Quantity)
			}
		}

		// Gán thông tin cho đơn hàng
		order.UserID = userID
		order.TotalPrice = totalPrice
		if order.Status == "" {
			order.Status = "Chờ xác nhận"
		}

		// Kiểm tra thông tin người nhận hàng
		if order.RecipientName == "" || order.RecipientPhone == "" {
			return errors.New("recipient name and phone are required")
		}

		// Tạo đơn hàng
		if err := tx.Model(&types.Order{}).Create(order).Error; err != nil {
			return err
		}

		// Lưu orderDetails từ map vào database
		for _, orderDetail := range orderDetailsMap {
			orderDetail.OrderID = order.ID // Gán OrderID vào từng OrderDetail
			if err := tx.Model(&types.OrderDetail{}).Create(orderDetail).Error; err != nil {
				return err
			}
		}

		// Xóa giỏ hàng sau khi đặt hàng thành công
		if err := tx.Table("cart").Where("user_id = ?", userID).Delete(&types.Carttypes{}).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		fmt.Println("Transaction failed:", err)
		return errors.New("failed to create order from cart: " + err.Error())
	}

	return nil
}
