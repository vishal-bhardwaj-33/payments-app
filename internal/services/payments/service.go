package payments

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	entities "payments-app/internal/entities/payments"
	repository "payments-app/internal/repository"
	er "payments-app/internal/utils"
	"payments-app/rpc/proto"

	"gorm.io/gorm"
)


func (s *PaymentServiceServerImpl) CreatePayment(ctx context.Context, req *proto.CreatePaymentRequest) (*proto.CreatePaymentResponse, error) {
	payment, err := entities.NewPayment(req.GetPayment().GetId(), req.GetPayment().GetUserId(), req.GetPayment().GetPaymentMethod(), req.GetPayment().GetStatus(), req.GetPayment().GetAmount())
	if err != nil {
		return nil, err
	}
	// fmt.Printf("Hello I am here in the create Payment ")

	// Check if payment already exists
	_, queryDBErr := repository.FindPaymentByID(payment.ID)
	if queryDBErr == nil {
		return nil, er.NewError(http.StatusConflict, fmt.Sprintf("Payment with ID %s already exists", payment.ID))
	}

	// Insert the payment into the database
	if dbErr := repository.InsertPayment(&payment); dbErr != nil {
		return nil, dbErr
	}

	// Return the created payment
	return entities.NewCreatePaymentResponse(&payment), nil
}

func (s *PaymentServiceServerImpl) UpdatePayment(ctx context.Context, req *proto.UpdatePaymentRequest) (*proto.UpdatePaymentResponse, error) {
    payment := req.GetPayment()
    paymentId := payment.GetId()

    // Check if the payment exists
    _, err := repository.FindPaymentByID(paymentId)
    if err != nil {
        return nil, fmt.Errorf("failed to find payment with ID %s: %v", paymentId, err)
    }

    // Directly call the repository function to update using COALESCE
    if err := repository.UpdatePaymentWithCOALESCE(paymentId, payment.GetAmount(), payment.GetStatus(), payment.GetPaymentMethod()); err != nil {
        return nil, fmt.Errorf("failed to update payment in database: %v", err)
    }

    // Fetch and return the updated payment record
    updatedPayment, err := repository.FindPaymentByID(paymentId)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch updated payment: %v", err)
    }

	
    return entities.NewUpdatePaymentResponse(updatedPayment),nil
}



func (s *PaymentServiceServerImpl) GetPaymentByID(ctx context.Context, req *proto.GetPaymentByIDRequest) (*proto.GetPaymentByIDResponse, error) {
	paymentID := req.GetId()

	// Query the payment from the database
	payment, err := repository.FindPaymentByID(paymentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, er.NewError(http.StatusNotFound, fmt.Sprintf("Payment with ID %s not found", paymentID))
		}
		return nil, err
	}

	// Return the payment details
	return entities.NewGetPaymentByIDResponse(payment), nil
}

func (s *PaymentServiceServerImpl) RandomDummyData(ctx context.Context, req *proto.RandomDummyDataRequest) (*proto.RandomDummyDataResponse, error) {
	count := int(req.GetCount()) // Get the count from the request

	// Fetch random payment data from external API
	payments, err := repository.FetchRandomPayments(count, s.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch random payments: %v", err)
	}

	// Convert []*entities.Payment to []*proto.Payment
	var protoPayments []*proto.Payment
	for _, p := range payments {
		protoPayments = append(protoPayments, &proto.Payment{
			Id:            p.ID,
			UserId:        p.UserID,
			Amount:        p.Amount,
			PaymentMethod: p.PaymentMethod,
			DateCreated:   p.DateCreated.Format(time.RFC3339),
			Status:        p.Status,
		})
	}

	// Return the payments
	return &proto.RandomDummyDataResponse{
		Payments: protoPayments,
	}, nil
}

