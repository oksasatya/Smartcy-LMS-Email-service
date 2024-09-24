package models

import (
	"time"
)

type Email struct {
	ID        uint32    `json:"id" gorm:"primary_key;auto_increment"`
	UserID    string    `json:"user_id" gorm:"type:varchar(50);not null"`
	EmailType string    `json:"email_type" gorm:"type:varchar(50);not null"`
	Email     string    `json:"email" gorm:"type:varchar(255);not null"`
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamp;default:current_timestamp"`
}
