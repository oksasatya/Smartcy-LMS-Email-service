package seeders

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func SeedAll(db *gorm.DB) {
	// Seed all data
	EmailSeeder(db)
	EmailLogsSeeder(db)
	logrus.Println("Seed all success")
}
