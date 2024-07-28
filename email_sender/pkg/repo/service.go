package repo

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect(databaseURL string) (*gorm.DB, error) {
	logLevel := logger.Info
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logLevel,
			IgnoreRecordNotFoundError: false,
			Colorful:                  true,
		},
	)
	return gorm.Open(postgres.New(postgres.Config{
		DSN: databaseURL,
	}), &gorm.Config{
		Logger: newLogger,
	})
}
