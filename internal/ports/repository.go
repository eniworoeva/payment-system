package ports

import "payment-system-one/internal/models"

type Repository interface {
	FindUserByEmail(email string) (*models.User, error)
	TokenInBlacklist(token *string) bool
	CreateUser(user *models.User) error
	UpdateUser(user *models.User) error
	FindAdminByEmail(email string) (*models.Admin, error)
	CreateAdmin(admin *models.Admin) error
	FindUserByAccountNumber(accountNumber int) (*models.User, error)
	TransferFunds(user *models.User, recipient *models.User, amount float64) error
	Transaction(account_no int) ([]models.Transaction, error)
}
