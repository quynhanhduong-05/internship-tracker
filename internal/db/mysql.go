package db

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectMySQL() *gorm.DB {
	dsn := "root:Qa180705@@tcp(127.0.0.1:3306)/internship_tracker?parseTime=true"

	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}

	log.Println("Connected to MySQL (GORM)")
	return database
}
