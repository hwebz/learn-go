package initializers

import (
	"fmt"
	"log"

	"github.com/hwebz/go-postgres-crud-rest-api/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB(env *Env) {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Saigon", env.DBHost, env.DBUser, env.DBPassword, env.DBName, env.DBPort)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database")
	}

	DB.Logger = logger.Default.LogMode(logger.Info)

	DB.AutoMigrate(&models.Feedback{})

	fmt.Println("🚀 Connected to the database")
}

type CreateFeedbackSchema struct {
	Name     string   `json:"name" validate:"required"`
	Email    string   `json:"email" validate:"required,email"`
	Feedback string   `json:"feedback" validate:"required"`
	Rating   *float32 `json:"rating" validate:"required"`
	Status   string   `json:"status,omitempty"`
}

type UpdateFeedbackSchema struct {
	Name     string   `json:"name,omitempty"`
	Email    string   `json:"email,omitempty"`
	Feedback string   `json:"feedback,omitempty"`
	Rating   *float32 `json:"rating,omitempty"`
	Status   string   `json:"status,omitempty"`
}
