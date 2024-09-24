package repository

import (
	"github.com/ghssni/Smartcy-LMS/Email-Service/internal/models"
	"gorm.io/gorm"
)

type EmailsLogRepository interface {
	InsertEmailLog(log *models.EmailLogs) (*models.EmailLogs, error)
	GetEmailLogById(id string) (*models.EmailLogs, error)
	GetAllEmailLogs() ([]models.EmailLogs, error)
	UpdateStatusLogs(id uint32, status string, errorMessage string) error
}

type emailsLogRepository struct {
	db *gorm.DB
}

func NewEmailsLogRepository(db *gorm.DB) EmailsLogRepository {
	return &emailsLogRepository{db}
}

// InsertEmailLog Insert a new email log
func (r *emailsLogRepository) InsertEmailLog(log *models.EmailLogs) (*models.EmailLogs, error) {
	result := r.db.Create(&log)
	if result.Error != nil {
		return nil, result.Error
	}
	return log, nil
}

// GetEmailLogById Get email log by ID
func (r *emailsLogRepository) GetEmailLogById(id string) (*models.EmailLogs, error) {
	var log models.EmailLogs
	result := r.db.First(&log, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &log, nil
}

// GetAllEmailLogs Get all email logs
func (r *emailsLogRepository) GetAllEmailLogs() ([]models.EmailLogs, error) {
	var logs []models.EmailLogs
	result := r.db.Find(&logs)
	if result.Error != nil {
		return nil, result.Error
	}
	return logs, nil
}

// UpdateStatusLogs Update status and error message for a specific log
func (r *emailsLogRepository) UpdateStatusLogs(id uint32, status string, errorMessage string) error {
	result := r.db.Model(&models.EmailLogs{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":        status,
		"error_message": errorMessage,
	})
	if result.Error != nil {
		return result.Error
	}
	return nil
}
