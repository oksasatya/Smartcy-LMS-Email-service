package seeders

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/ghssni/Smartcy-LMS/Email-Service/internal/models"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

func EmailSeeder(db *gorm.DB) {
	var emails []models.Email

	gofakeit.Seed(0)

	for i := 0; i < 10; i++ {
		email := models.Email{
			UserID:    gofakeit.UUID(),
			EmailType: gofakeit.RandomString([]string{"forgot_password", "payment_due", "payment_success"}),
			Email:     gofakeit.Email(),
			CreatedAt: time.Now(),
		}
		emails = append(emails, email)
	}
	if err := db.Create(&emails).Error; err != nil {
		logrus.Fatalf("Failed to seed Emails: %v", err)
	} else {
		logrus.Println("Emails seeded successfully")
	}
}
