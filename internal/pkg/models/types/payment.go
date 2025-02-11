package types

import "time"

type Payment struct {
    ID            int       `json:"id"`
    OrderID       int       `json:"order_id"`
    PaymentMethod string    `json:"payment_method"`
    PaymentStatus string    `json:"payment_status"`
    TransactionID string    `json:"transaction_id"`
    Amount        float64   `json:"amount"`
    PaymentDate   time.Time `json:"payment_date"`
}
