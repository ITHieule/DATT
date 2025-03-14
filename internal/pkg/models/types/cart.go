package types

type Carttypes struct {
	ID                    uint                    `json:"id" gorm:"primaryKey"`
	UserID                int                     `json:"user_id" gorm:"column:user_id"`
	ProductDetailResponse []ProductDetailResponse `json:"product_variants" gorm:"foreignKey:CartID;references:ID"`
	Quantity              int                     `json:"quantity"`
}

func (Carttypes) TableName() string {
	return "cart" 
}
