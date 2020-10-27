package db

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" // blank import is for gorm
	"github.com/joho/godotenv"
)

/*
	We do not close db everytime we get it.
	The db object lives for the life of the program.
	However the following needs to be set for safety.
	connection.DB().SetConnMaxLifetime(time.Minute * 10)
	connection.DB().SetMaxIdleConns(1)
	connection.DB().SetMaxOpenConns(0)
*/

var db *gorm.DB

func init() {

	configPath, err := filepath.Abs("../src/github.com/Robert1138/GWA/.env")
	err = godotenv.Load(configPath)
	if err != nil {
		fmt.Println(err)
	}
	dbName := os.Getenv("db_name")
	dbPassword := os.Getenv("db_pass")
	dbHost := os.Getenv("db_host")
	dbUser := os.Getenv("db_user")
	dbLogmode := false

	if os.Getenv("db_logmode") == "true" {
		dbLogmode = true
	}
	// mysql connection
	dbString := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true&loc=Local", dbUser, dbPassword, dbHost, dbName)

	connection, e := gorm.Open("mysql", dbString)
	if e != nil {
		fmt.Println(e)
	}
	connection.DB().SetConnMaxLifetime(time.Minute * 10)
	connection.DB().SetMaxIdleConns(1)
	connection.DB().SetMaxOpenConns(0)
	db = connection
	db.SingularTable(true) // this makes sure gorm is using singular table names in queries
	db.LogMode(dbLogmode)
	fmt.Println("DB initialized")
}

func InitDB() error {
	dbName := os.Getenv("db_name")
	dbPassword := os.Getenv("db_pass")
	dbHost := os.Getenv("db_host")
	dbUser := os.Getenv("db_user")
	dbLogmode := false

	if os.Getenv("db_logmode") == "true" {
		dbLogmode = true
	}
	// mysql connection
	dbString := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true&loc=Local", dbUser, dbPassword, dbHost, dbName)

	connection, e := gorm.Open("mysql", dbString)
	if e != nil {
		fmt.Println(e)
		return e
	}
	connection.DB().SetConnMaxLifetime(time.Minute * 10)
	connection.DB().SetMaxIdleConns(1)
	connection.DB().SetMaxOpenConns(0)
	db = connection
	db.SingularTable(true) // this makes sure gorm is using singular table names in queries
	db.LogMode(dbLogmode)
	fmt.Println("DB initialized")
	return nil
}

// GetDbConn will return the global db connection variable
func GetDbConn() *gorm.DB {
	return db
}
