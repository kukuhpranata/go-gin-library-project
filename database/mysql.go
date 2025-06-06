package database

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewMysqlClient() (*gorm.DB, error) {
	godotenv.Load()
	var (
		DB_User         = os.Getenv("DB_USERNAME")
		DB_Pass         = os.Getenv("DB_PASSWORD")
		DB_Host         = os.Getenv("DB_HOST")
		DB_Port         = os.Getenv("DB_PORT")
		DB_DatabaseName = os.Getenv("DB_DATABASE")
	)
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local",
		DB_User, DB_Pass, DB_Host, DB_Port, DB_DatabaseName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
