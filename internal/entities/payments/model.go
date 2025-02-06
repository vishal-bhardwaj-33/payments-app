package entities

import (
	"net/http"
	"time"

	er "payments-app/internal/utils"
	"payments-app/rpc/proto"
)

// Payment model (no DB dependencies)
type Payment struct {
	ID            string    `json:"id" gorm:"primaryKey"`
	UserID        string    `json:"user_id"`
	Amount        float32   `json:"amount"`
	PaymentMethod string    `json:"payment_method"`
	Status        string    `json:"status"`
	DateCreated   time.Time `json:"date_created"`
}

// NewPayment creates a new Payment
func NewPayment(id, userID, paymentMethod, status string, amount float32) (Payment, error) {
	if id == "" || userID == "" || paymentMethod == "" || status == "" {
		return Payment{}, er.NewError(http.StatusBadRequest, "All fields are required")
	}

	return Payment{
		ID:            id,
		UserID:        userID,
		Amount:        amount,
		PaymentMethod: paymentMethod,
		Status:        status,
		DateCreated:   time.Now(),
	}, nil
}

// NewCreatePaymentResponse converts the Payment entity to a CreatePaymentResponse proto
func NewCreatePaymentResponse(payment *Payment) *proto.CreatePaymentResponse {
	return &proto.CreatePaymentResponse{
		Payment: &proto.Payment{
			Id:            payment.ID,
			UserId:        payment.UserID,
			Amount:        payment.Amount,
			PaymentMethod: payment.PaymentMethod,
			Status:        payment.Status,
			DateCreated:   payment.DateCreated.Format(time.RFC3339),
		},
	}
}

// NewUpdatePaymentResponse converts the Payment entity to an UpdatePaymentResponse proto
func NewUpdatePaymentResponse(payment *Payment) *proto.UpdatePaymentResponse {
	return &proto.UpdatePaymentResponse{
		Payment: &proto.Payment{
			Id:            payment.ID,
			UserId:        payment.UserID,
			Amount:        payment.Amount,
			PaymentMethod: payment.PaymentMethod,
			Status:        payment.Status,
			DateCreated:   payment.DateCreated.Format(time.RFC3339),
		},
	}
}

// NewGetPaymentByIDResponse converts the Payment entity to a GetPaymentByIDResponse proto
func NewGetPaymentByIDResponse(payment *Payment) *proto.GetPaymentByIDResponse {
	return &proto.GetPaymentByIDResponse{
		Payment: &proto.Payment{
			Id:            payment.ID,
			UserId:        payment.UserID,
			Amount:        payment.Amount,
			PaymentMethod: payment.PaymentMethod,
			Status:        payment.Status,
			DateCreated:   payment.DateCreated.Format(time.RFC3339),
		},
	}
}

