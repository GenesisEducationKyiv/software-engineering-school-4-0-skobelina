package repo

import (
	"log"
	"os"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
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

func Mock() (sqlmock.Sqlmock, *gorm.DB, error) {
	testDB, mockDB, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 testDB,
		PreferSimpleProtocol: true,
	}))
	return mockDB, db, err
}
