package request

type CreatePaymentRequest struct {
	ID            int     `json:"id"`
	OrderID       int     `json:"order_id"`
	PaymentMethod string  `json:"payment_method"`
	PaymentStatus string  `json:"payment_status"`
	TransactionID *string `json:"transaction_id,omitempty"`
	Amount        float64 `json:"amount"`
	PaymentDate   string  `json:"payment_date"`
}

type RefundPaymentRequest struct {
	TransactionID string `json:"transaction_id" validate:"required"`
}
