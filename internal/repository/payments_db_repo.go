package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"payments-app/internal/database"
	entities "payments-app/internal/entities/payments"
	"payments-app/internal/utils"

	"gorm.io/gorm"
)

// DeletePaymentWithID deletes payment data by ID
func DeletePaymentWithID(paymentID string) (int, error) {
	err := database.DB.Where("id = ?", paymentID).Delete(&entities.Payment{}).Error
	if err != nil {
		log.Printf("DeletePaymentWithID: Failed to delete payment data in DB: %v", err)
		return 0, er.NewError(http.StatusServiceUnavailable, fmt.Sprintf("failed to delete payment data in DB: %s", err.Error()))
	}
	return 1, nil
}

// InsertPayment inserts payment data
func InsertPayment(payment *entities.Payment) error {
	err := database.DB.Create(&payment).Error
	if err != nil {
		log.Printf("InsertPayment: Failed to insert payment data in DB: %v", err)
		return er.NewError(http.StatusInternalServerError, fmt.Sprintf("failed to insert payment data in DB: %s", err.Error()))
	}
	return nil
}

// QueryPayment queries payment data based on payment ID and updates instance with the latest entry
func QueryPayment(payment *entities.Payment) error {
	err := database.DB.Where("id = ?", payment.ID).First(&payment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("QueryPayment: entry not found for paymentID: %s, failed with error: %s", payment.ID, err.Error())
			return err
		}
		log.Printf("QueryPayment: failed to retrieve payment data from DB: %v", err)
		return er.NewError(http.StatusInternalServerError, fmt.Sprintf("failed to retrieve payment data from DB: %s", payment.ID))
	}
	return nil
}

// UpdatePayment updates payment data
func UpdatePayment(payment *entities.Payment) error {
	if err := database.DB.Model(&entities.Payment{}).Where("id = ?", payment.ID).Updates(payment).Error; err != nil {
		log.Printf("UpdatePayment: failed to update payment data for ID: %s in DB: %v", payment.ID, err)
		return er.NewError(http.StatusInternalServerError, fmt.Sprintf("failed to update payment data for ID %s", payment.ID))
	}
	return nil
}

func UpdatePaymentWithCOALESCE(id string, amount float32, status, paymentMethod string) error {
    query := `
    UPDATE payments
    SET 
        amount = COALESCE(NULLIF($1, 0), amount), 
        status = COALESCE(NULLIF($2, ''), status),
        payment_method = COALESCE(NULLIF($3, ''), payment_method)
    WHERE id = $4`

    err := database.DB.Exec(query, amount, status, paymentMethod, id)
    if err.RowsAffected < 1 {
        log.Printf("UpdatePaymentWithCOALESCE: failed to update payment data for ID: %s in DB: %v", id, err.RowsAffected)
        return fmt.Errorf("failed to update payment: %v", err.Error)
    }

    return nil
}

// UpdateOnConflict updates payment data if already present else creates a new row
func UpdateOnConflict(payment *entities.Payment) error {
	fmt.Printf("Updating for payment %v", payment)
	err := database.DB.Where("id = ?", payment.ID).Assign(payment).FirstOrCreate(&payment).Error
	if err != nil {
		log.Printf("Failed to update payment data in DB: %v", err)
		return fmt.Errorf("failed to update payment data in DB: %w", err)
	}
	return nil
}

// FindPaymentByID finds payment data by ID and returns updated instance
func FindPaymentByID(paymentID string) (*entities.Payment, error) {
	var payment entities.Payment
	err := database.DB.Where("id = ?", paymentID).First(&payment).Error
	if err != nil {
		return nil, er.NewError(http.StatusInternalServerError, fmt.Sprintf("failed to find payment data for ID %s", paymentID))
	}
	return &payment, nil
}

// FetchRandomPayments fetch data from external API
func FetchRandomPayments(count int) ([]*entities.Payment, error) {
	url := fmt.Sprintf("http://localhost:3000/payments?count=%d", count)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to call external API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("external API returned status: %d", resp.StatusCode)
	}

	// Parse response body
	var payments []entities.Payment
	if err := json.NewDecoder(resp.Body).Decode(&payments); err != nil {
		return nil, fmt.Errorf("failed to parse external API response: %v", err)
	}

	// Convert to []*proto.Payment
	var paymentPointers []*entities.Payment
	for i := range payments {
		paymentPointers = append(paymentPointers, &payments[i])
	}

	return paymentPointers, nil
}