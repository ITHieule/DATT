package request

type CreatePaymentRequest struct {
	OrderID       int     `json:"order_id" validate:"required"`
	PaymentMethod string  `json:"payment_method" validate:"required"`
	Amount        float64 `json:"amount" validate:"required"`
}

type RefundPaymentRequest struct {
	TransactionID string `json:"transaction_id" validate:"required"`
}
