package request

type CreateAddressRequest struct {
    UserID     int    `json:"user_id" validate:"required"`
    FullAddress string `json:"full_address" validate:"required"`
    City       string `json:"city" validate:"required"`
    PostalCode string `json:"postal_code" validate:"required"`
    Country    string `json:"country" validate:"required"`
}
