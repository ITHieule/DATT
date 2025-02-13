package services

import (
	"fmt"
	"strings"

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

func (s *ProductsService) Product_imageSevice() ([]types.Product_image, error) {
	var pro []types.Product_image

	// Kết nối database
	db, err := database.FashionBusiness()
	if err != nil {
		fmt.Println("Database connection error:", err)
		return nil, err
	}
	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	// Truy vấn SQL lấy tất cả hình ảnh của sản phẩm và nhóm chúng thành một chuỗi
	query := `
		SELECT products.id, 
		       products.name,
		       products.base_price, 
		       GROUP_CONCAT(product_images.image_url ORDER BY product_images.id) AS image_urls
		FROM fashion_shop.products products
		LEFT JOIN fashion_shop.product_images product_images 
		ON products.id = product_images.product_id
		GROUP BY products.id;
	`

	// Thực hiện truy vấn và ánh xạ kết quả vào struct
	err = db.Raw(query).Scan(&pro).Error
	if err != nil {
		fmt.Println("Query execution error:", err)
		return nil, err
	}

	// Chuyển đổi `ImageURLsRaw` thành `ImageURLs`
	for i := range pro {
		if pro[i].ImageURLsRaw != "" { // Kiểm tra nếu có dữ liệu ảnh
			pro[i].ImageURLs = strings.Split(pro[i].ImageURLsRaw, ",")
		} else {
			pro[i].ImageURLs = []string{} // Nếu không có ảnh, trả về slice rỗng
		}
	}

	return pro, nil
}
