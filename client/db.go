package client

import (
	"context"
	"meDemo/constant"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func DB() *gorm.DB {
	return db
}

func DBWithContext(ctx context.Context) *gorm.DB {
	return db.WithContext(ctx)
}

func ConnectDB() error {
	return ConnectDBWithConfig(&gorm.Config{})
}

func ConnectDBWithConfig(gormOption gorm.Option) error {
	dbInt, err := gorm.Open(postgres.Open(constant.DatabaseURL()), gormOption)
	if err != nil {
		return err
	}

	db = dbInt
	return nil
}
