package repository

import (
	"payment-system-one/internal/models"
	"time"
)

func (p *Postgres) FindUserByEmail(email string) (*models.User, error) {
	user := &models.User{}

	if err := p.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// create a user in thye database
func (p *Postgres) CreateUser(user *models.User) error {
	if err := p.DB.Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (p *Postgres) UpdateUser(user *models.User) error {
	if err := p.DB.Save(user).Error; err != nil {
		return err
	}
	return nil
}

// FindUserByAccountNumber
func (p *Postgres) FindUserByAccountNumber(accountNumber int) (*models.User, error) {
	user := &models.User{}

	if err := p.DB.Where("account_no = ?", accountNumber).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// TransferFunds
func (p *Postgres) TransferFunds(user *models.User, recipient *models.User, amount float64) error {
	tx := p.DB.Begin()

	// deduct the amount from the payer
	user.AvailableBalance -= amount
	// add the amount to the recipient
	recipient.AvailableBalance += amount

	// save the transaction for the payer
	if err := tx.Save(user).Error; err != nil {
		tx.Rollback()
		return err
	}

	// save the transaction for the recipient
	if err := tx.Save(recipient).Error; err != nil {
		tx.Rollback()
		return err
	}

	// save the transaction in the transaction table
	transaction := &models.Transaction{
		PayerAccountNumber:     user.AccountNo,
		RecipientAccountNumber: recipient.AccountNo,
		TransactionType:        "debit",
		TransactionAmount:      amount,
		TransactionDate:        time.Now(),
	}

	// save the transaction
	if err := tx.Create(transaction).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

// Transaction
func (p *Postgres) Transaction(account_no int) ([]models.Transaction, error) {
	transactions := []models.Transaction{}

	if err := p.DB.Where("payer_account_number = ? OR recipient_account_number = ?", account_no, account_no).Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}