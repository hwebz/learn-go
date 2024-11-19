package initializers

import (
	"fmt"
	"log"
	"os"

	"github.com/go-jwt-auth/models"
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
		log.Fatal("Failed to connect to the Database: \n", err.Error())
	}

	DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	DB.Logger = logger.Default.LogMode(logger.Info)

	err = DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("Migration Failed: \n", err.Error())
		os.Exit(1)
	}

	fmt.Println("🚀 Database connected successfully!")
}
