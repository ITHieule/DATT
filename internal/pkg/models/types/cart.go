package types

type Carttypes struct {
	Id                 uint `json:"id"`
	User_Id            int  `json:"user_id"`
	Product_Variant_Id int  `json:"product_variant_id"`
	Quantity           int  `json:"quantity"`
}
