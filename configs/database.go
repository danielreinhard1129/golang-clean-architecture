package configs

import (
	"fmt"
	"log"

	"github.com/danielreinhard1129/fiber-clean-arch/pkg/exception"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDatabase(config Config) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		config.Get("POSTGRES_HOST"),
		config.Get("POSTGRES_USER"),
		config.Get("POSTGRES_PASSWORD"),
		config.Get("POSTGRES_DB"),
		config.Get("POSTGRES_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	exception.PanicLogging(err)

	log.Println("✅ Database connected")
	return db
}
