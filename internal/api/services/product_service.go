package services

import (
	"fmt"
	"strings"
	"time"

	"web-api/internal/pkg/database"
	"web-api/internal/pkg/models/request"
	"web-api/internal/pkg/models/types"

	"gorm.io/gorm"
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

// GetProductByID lấy chi tiết sản phẩm theo ID
func (s *ProductsService) GetProductByID(productID int) (*types.ProductDetailResponse, error) {
	// Kết nối database
	db, err := database.FashionBusiness()
	if err != nil {
		fmt.Println("Database connection error:", err)
		return nil, err
	}
	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	// Biến chứa thông tin sản phẩm
	var product struct {
		ID          int     `json:"id"`
		Name        string  `json:"name"`
		BasePrice   float64 `json:"base_price"`
		Description string  `json:"description"`
		ImageURLs   string  `json:"image_urls"`
	}

	// Truy vấn sản phẩm và hình ảnh
	query := `
		SELECT p.id, p.name, p.base_price, p.description, 
		       COALESCE(GROUP_CONCAT(pi.image_url ORDER BY pi.id), '') AS image_urls
		FROM products p
		LEFT JOIN product_images pi ON p.id = pi.product_id
		WHERE p.id = ?	
		GROUP BY p.id;
	`

	err = db.Raw(query, productID).Scan(&product).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("product not found")
		}
		return nil, fmt.Errorf("query error: %w", err)
	}

	// Lấy danh sách biến thể sản phẩm
	var variants []types.ProductVariant
	err = db.Table("product_variants").
		Where("product_id = ?", productID).
		Find(&variants).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get product variants: %w", err)
	}

	// Chuyển đổi danh sách ảnh thành slice
	imageURLs := []string{}
	if product.ImageURLs != "" {
		imageURLs = strings.Split(product.ImageURLs, ",")
	}

	// Trả về kết quả
	return &types.ProductDetailResponse{
		ID:          product.ID,
		Name:        product.Name,
		BasePrice:   product.BasePrice,
		Description: product.Description,
		Variants:    variants,
		Images:      imageURLs,
	}, nil
}

// AddProductService thêm sản phẩm mới, bao gồm biến thể và ảnh sản phẩm
func (s *ProductsService) AddProductService(requestParams *request.CreateProductRequest) (types.Product, error) {
	// Kết nối database
	db, err := database.FashionBusiness()
	if err != nil {
		fmt.Println("Database connection error:", err)
		return types.Product{}, err
	}
	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	// Bắt đầu transaction để đảm bảo tính toàn vẹn dữ liệu
	tx := db.Begin()

	// Thêm sản phẩm mới
	product := types.Product{
		CategoryID:  requestParams.CategoryID,
		Name:        requestParams.Name,
		Description: requestParams.Description,
		BasePrice:   requestParams.BasePrice,
	}

	if err := tx.Create(&product).Error; err != nil {
		tx.Rollback() // Hoàn tác nếu có lỗi
		fmt.Println("Error inserting product:", err)
		return types.Product{}, err
	}

	// Thêm danh sách biến thể nếu có
	var variants []types.ProductVariant
	for _, variant := range requestParams.Variants {
		variant.ProductID = product.ID // Gán ID sản phẩm cho biến thể
		variants = append(variants, variant)
	}

	if len(variants) > 0 {
		if err := tx.Create(&variants).Error; err != nil {
			tx.Rollback() // Hoàn tác nếu có lỗi
			fmt.Println("Error inserting product variants:", err)
			return types.Product{}, err
		}
	}

	// Thêm danh sách ảnh nếu có
	var images []types.ProductImage
	for _, imageURL := range requestParams.Images {
		image := types.ProductImage{
			ProductID: product.ID,
			ImageURL:  imageURL,
		}
		images = append(images, image)
	}

	if len(images) > 0 {
		if err := tx.Create(&images).Error; err != nil {
			tx.Rollback() // Hoàn tác nếu có lỗi
			fmt.Println("Error inserting product images:", err)
			return types.Product{}, err
		}
	}

	// Commit transaction nếu không có lỗi
	tx.Commit()

	return product, nil
}

func (s *ProductsService) UpdateProductService(requestParams *request.CreateProductRequest) (types.Product, error) {
	// Kết nối database
	db, err := database.FashionBusiness()
	if err != nil {
		fmt.Println("Database connection error:", err)
		return types.Product{}, err
	}
	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	// Bắt đầu transaction để đảm bảo tính toàn vẹn dữ liệu
	tx := db.Begin()

	// Kiểm tra xem sản phẩm có tồn tại không
	var product types.Product
	if err := tx.First(&product, requestParams.Id).Error; err != nil {
		tx.Rollback()
		fmt.Println("Product not found:", err)
		return types.Product{}, err
	}

	// Cập nhật thông tin sản phẩm
	product.Name = requestParams.Name
	product.Description = requestParams.Description
	product.BasePrice = requestParams.BasePrice
	product.CategoryID = requestParams.CategoryID

	if err := tx.Save(&product).Error; err != nil {
		tx.Rollback()
		fmt.Println("Error updating product:", err)
		return types.Product{}, err
	}

	// Xóa các biến thể cũ và thêm mới
	if err := tx.Where("product_id = ?", product.ID).Delete(&types.ProductVariant{}).Error; err != nil {
		tx.Rollback()
		fmt.Println("Error deleting old product variants:", err)
		return types.Product{}, err
	}

	var variants []types.ProductVariant
	for _, variant := range requestParams.Variants {
		variant.ProductID = product.ID
		variants = append(variants, variant)
	}

	if len(variants) > 0 {
		if err := tx.Create(&variants).Error; err != nil {
			tx.Rollback()
			fmt.Println("Error inserting new product variants:", err)
			return types.Product{}, err
		}
	}

	// Xóa các ảnh cũ và thêm mới
	if err := tx.Where("product_id = ?", product.ID).Delete(&types.ProductImage{}).Error; err != nil {
		tx.Rollback()
		fmt.Println("Error deleting old product images:", err)
		return types.Product{}, err
	}

	var images []types.ProductImage
	for _, imageURL := range requestParams.Images {
		image := types.ProductImage{
			ProductID: product.ID,
			ImageURL:  imageURL,
		}
		images = append(images, image)
	}

	if len(images) > 0 {
		if err := tx.Create(&images).Error; err != nil {
			tx.Rollback()
			fmt.Println("Error inserting new product images:", err)
			return types.Product{}, err
		}
	}

	// Commit transaction nếu không có lỗi
	tx.Commit()

	return product, nil
}

func (s *ProductsService) DeleteproductSevice(Id int) error {

	// Kết nối database
	db, err := database.FashionBusiness()
	if err != nil {
		fmt.Println("Database connection error:", err)

		return err
	}
	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	// Thực thi lệnh DELETE
	result := db.Exec("DELETE FROM products WHERE id = ?", Id)

	// Kiểm tra lỗi truy vấn
	if result.Error != nil {
		fmt.Println("Query execution error:", result.Error)
		return result.Error
	}

	// Kiểm tra số dòng bị ảnh hưởng (nếu ID không tồn tại, sẽ không xóa được)
	if result.RowsAffected == 0 {
		fmt.Println("No products found with ID:", Id)
		return fmt.Errorf("không tìm thấy products với ID %d", Id)
	}

	fmt.Println("Deleted products successfully!")
	return nil
}

func (s *ProductsService) SearchProductService(requestParams *request.CreateProductRequest) ([]types.Product, error) {
	var products []types.Product

	// Kết nối database
	db, err := database.FashionBusiness()
	if err != nil {
		fmt.Println("Database connection error:", err)
		return nil, err
	}
	dbInstance, _ := db.DB()
	defer dbInstance.Close()

	// Truy vấn SQL tìm kiếm sản phẩm
	query := `
		SELECT * FROM products 
		WHERE id = ? OR name LIKE ?
	`
	searchName := "%" + requestParams.Name + "%" // Tìm kiếm gần đúng

	err = db.Raw(query, requestParams.Id, searchName).Scan(&products).Error
	if err != nil {
		fmt.Println("Query execution error:", err)
		return nil, err
	}
	return products, nil
}

func (s *ProductsService) GetLatestProductHots(limit int) ([]types.ProductHot, error) {
    var results []struct {
        ProductID   int       `json:"product_id"`
        CategoryID  int       `json:"category_id"`
        Name        string    `json:"name"`
        Description string    `json:"description"`
        BasePrice   float64   `json:"base_price"`
        CreatedAt   time.Time `json:"created_at"`
        ImageURL    string    `json:"image_url"`
    }

    db, err := database.FashionBusiness()
    if err != nil {
        fmt.Println("Database connection error:", err)
        return nil, err
    }
    dbInstance, _ := db.DB()
    defer dbInstance.Close()

    // Truy vấn lấy sản phẩm mới nhất với hình ảnh
    query := `
        SELECT p.id AS product_id, p.category_id, p.name, p.description, p.base_price, p.created_at, 
               COALESCE(pi.image_url, '') AS image_url
        FROM products p
        LEFT JOIN product_images pi ON pi.product_id = p.id
        ORDER BY p.created_at DESC
        LIMIT ?
    `

    // Thực hiện truy vấn
    err = db.Raw(query, limit).Scan(&results).Error
    if err != nil {
        fmt.Println("Query execution error:", err)
        return nil, err
    }

    // Chuyển đổi kết quả thành dạng `ProductHot`
    var productHots []types.ProductHot
    for _, item := range results {
        productHots = append(productHots, types.ProductHot{
            ImageURL:  item.ImageURL,
            Product: types.Product{
                ID:          item.ProductID,
                CategoryID:  item.CategoryID,
                Name:        item.Name,
                Description: item.Description,
                BasePrice:   item.BasePrice,
                CreatedAt:   item.CreatedAt,
            },
        })
    }

    return productHots, nil
}

