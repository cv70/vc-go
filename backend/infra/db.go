package infra

import (
	"context"
	"fmt"
	"vc-go/config"
	"vc-go/datasource/dbdao"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDB(ctx context.Context, c *config.DatabaseConfig) (*dbdao.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True", c.User, c.Password, c.Host, c.DBName)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	return dbdao.NewDB(db), err
}
