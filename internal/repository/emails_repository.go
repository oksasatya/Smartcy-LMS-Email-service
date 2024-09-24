package repository

import (
	"github.com/ghssni/Smartcy-LMS/Email-Service/internal/models"
	"gorm.io/gorm"
)

type EmailsRepository interface {
	InsertEmail(email *models.Email) (*models.Email, error)
	GetEmails() ([]models.Email, error)
	GetEmailById(id uint32) (*models.Email, error)
}

type emailsRepository struct {
	db *gorm.DB
}

func (r *emailsRepository) InsertEmail(email *models.Email) (*models.Email, error) {
	result := r.db.Create(&email)
	if result.Error != nil {
		return nil, result.Error
	}

	return email, nil
}

func (r *emailsRepository) GetEmails() ([]models.Email, error) {
	var emails []models.Email
	result := r.db.Find(&emails)
	if result.Error != nil {
		return nil, result.Error
	}

	return emails, nil
}

func (r *emailsRepository) GetEmailById(id uint32) (*models.Email, error) {
	var email models.Email
	result := r.db.First(&email, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &email, nil
}

func NewEmailsRepository(db *gorm.DB) EmailsRepository {
	return &emailsRepository{db}
}
