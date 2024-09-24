package models

import "time"

type EmailLogs struct {
	ID           int       `json:"id" gorm:"primary_key;auto_increment"`
	UserID       string    `json:"user_id" gorm:"type:varchar(50);not null"`
	Email        string    `json:"email" gorm:"type:varchar(255);not null"`
	EnrollmentID int       `json:"enrollment_id" gorm:"type:int;"`
	InvoiceID    string    `json:"invoice_id" gorm:"type:varchar(50);"`
	Status       string    `json:"status" gorm:"type:varchar(20);not null;default:'pending'"`
	SentAt       time.Time `json:"sent_at" gorm:"type:timestamp;default:current_timestamp"`
	ErrorMessage string    `json:"error_message" gorm:"type:text"`
}
