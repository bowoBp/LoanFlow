package db

import (
	"fmt"
	"github.com/bowoBp/LoanFlow/pkg/reader"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

func Default() (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		reader.GetEnv("PGHOST"),
		reader.GetEnv("PGUSER"),
		reader.GetEnv("PGPASSWORD"),
		reader.GetEnv("PGDB"),
		reader.GetEnv("PGPORT"),
		reader.GetEnv("PGSSL"),
	)

	newLogger := gormLogger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		gormLogger.Config{
			SlowThreshold:             200 * time.Millisecond, // Slow SQL threshold
			LogLevel:                  gormLogger.Info,        // Log level
			IgnoreRecordNotFoundError: true,                   // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,                   // Disable color
		},
	)

	pgqldb, err := gorm.Open(postgres.New(postgres.Config{
		DSN: dsn,
	}), &gorm.Config{
		Logger:                 newLogger,
		SkipDefaultTransaction: true,
	})
	if err != nil {
		return pgqldb, err
	} else {
		return pgqldb, nil
	}
}
