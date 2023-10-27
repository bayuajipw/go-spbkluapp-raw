package config

import (
	"database/sql"
	"fmt"
	"log"
	"spbkluapp/util"

	_ "github.com/go-sql-driver/mysql"
	// "gorm.io/driver/mysql"
	// "gorm.io/gorm"
)

var DB *sql.DB

// ConnectDB establishes a connection to the PostgreSQL database
func Connect() error {
	var err error

	DBHOST := util.GetDotenv("DB_HOST")
	DBUSER := util.GetDotenv("DB_USER")
	DBPASSWORD := util.GetDotenv("DB_PASSWORD")
	DBNAME := util.GetDotenv("DB_NAME")
	DBPORT := util.GetDotenv("DB_PORT")

	DB_URI := DBUSER + ":" + DBPASSWORD + "@tcp(" + DBHOST + ":" + DBPORT + ")/" + DBNAME + "?charset=utf8mb4&parseTime=True&loc=Local"

	fmt.Println()

	DB, err = sql.Open("mysql", DB_URI)

	if err != nil {
		panic(err)
	}

	log.Println("Database connected")

	return nil
}
