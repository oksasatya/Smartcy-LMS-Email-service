package database

import (
	"fmt"
	"github.com/ghssni/Smartcy-LMS/Email-Service/database/migrations"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"sync"
)

var (
	db   *gorm.DB
	once sync.Once
)

// InitDB is a function to initialize the database connection
func InitDB() (*gorm.DB, error) {
	var err error

	once.Do(func() {
		host := os.Getenv("DB_HOST")
		user := os.Getenv("DB_USER")
		password := os.Getenv("DB_PASSWORD")
		dbname := os.Getenv("DB_NAME")
		port := os.Getenv("DB_PORT")

		// Build DSN
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbname, port)

		// Open connection to the database
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			PrepareStmt: true,
			Logger:      logger.Default.LogMode(logger.Info),
		})
		if err != nil {
			log.Fatalf("Error connecting to the database: %v", err)
		}

		// Run migrations
		if err := migrations.Migrate(db); err != nil {
			err = fmt.Errorf("migration failed: %v", err)
			return
		}
	})

	// Return instance or error
	return db, err
}
