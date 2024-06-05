package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName        string  `json:"first_name"`
	LastName         string  `json:"last_name"`
	Password         string  `json:"password"`
	DateOfBirth      string  `json:"date_of_birth"`
	Email            string  `json:"email"`
	AccountNo        int     `json:"account_no"`
	AvailableBalance float64 `json:"available_balance"`
	Phone            string  `json:"phone"`
	Address          string  `json:"address"`
}

type Admin struct {
	gorm.Model
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Password    string `json:"password"`
	DateOfBirth string `json:"date_of_birth"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Address     string `json:"address"`
}

//type UserProfile struct {
//	gorm.Model
//	ValidIdentity string `json:"valid_identity"`
//	PassPort string `json:"passport"`
//
//}

type Transaction struct {
	gorm.Model
	PayerAccountNumber     int     `json:"payer_account_number"`
	RecipientAccountNumber int     `json:"recipient_account_number"`
	TransactionType        string  `json:"transaction_type"`
	TransactionAmount      float64 `json:"transaction_amount"`
	// BalanceBefore          float64   `json:"balance_before"`
	// BalanceAfter           float64   `json:"balance_after"`
	TransactionDate time.Time `json:"transaction_date"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TransferRequest struct {
	AccountNumber int     `json:"account_number"`
	Amount        float64 `json:"amount"`
}
