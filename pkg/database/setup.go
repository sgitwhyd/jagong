package database

import (
	"fmt"
	"log"

	"github.com/sgitwhyd/jagong/app/models"
	"github.com/sgitwhyd/jagong/pkg/env"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func SetupDatabase() {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		env.GetEnv("DB_USER", ""),
		env.GetEnv("DB_PASSWORD", ""),
		env.GetEnv("DB_HOST", "127.0.0.1"),
		env.GetEnv("DB_PORT", "3306"),
		env.GetEnv("DB_NAME", ""),
	)
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	err = DB.AutoMigrate(&models.User{}, &models.UserSession{})
	if err != nil {
		log.Printf("AutoMigrate err:%v\n", err)
	}

	log.Print("Database setup success")
}
