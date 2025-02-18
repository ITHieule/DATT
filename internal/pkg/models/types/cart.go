package types

// Carttypes đại diện cho giỏ hàng của người dùng
type Carttypes struct {
	Id              uint             `json:"id"`
	User_Id         int              `json:"user_id"`          // ID người dùng
	ProductDetailResponse []ProductDetailResponse `json:"product_variants"` // Mảng các sản phẩm trong giỏ hàng
	Quantity        int              `json:"quantity"`         // Tổng số lượng sản phẩm trong giỏ
}
