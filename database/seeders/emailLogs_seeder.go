package seeders

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/ghssni/Smartcy-LMS/Email-Service/internal/models"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func EmailLogsSeeder(db *gorm.DB) {
	// Seed EmailLogs data
	var emailLogs []models.EmailLogs

	gofakeit.Seed(0)

	for i := 0; i < 10; i++ {
		emailLog := models.EmailLogs{
			UserID:       gofakeit.UUID(),
			Email:        gofakeit.Email(),
			EnrollmentID: gofakeit.Number(1, 100),
			InvoiceID:    gofakeit.UUID(),
			Status:       gofakeit.RandomString([]string{"pending", "sent", "failed"}),
			SentAt:       gofakeit.Date(),
			ErrorMessage: gofakeit.Sentence(10),
		}
		emailLogs = append(emailLogs, emailLog)
	}

	if err := db.Create(&emailLogs).Error; err != nil {
		logrus.Fatalf("Failed to seed EmailLogs: %v", err)
	} else {
		logrus.Println("EmailLogs seeded successfully")
	}

}
