package repository

import "payment-system-one/internal/models"

func (p *Postgres) FindAdminByEmail(email string) (*models.Admin, error) {
	admin := &models.Admin{}

	if err := p.DB.Where("email = ?", email).First(&admin).Error; err != nil {
		return nil, err
	}
	return admin, nil
}

func (p *Postgres) CreateAdmin(admin *models.Admin) error {
	if err := p.DB.Create(admin).Error; err != nil {
		return err
	}
	return nil
}
