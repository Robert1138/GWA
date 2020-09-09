package db

import (
	"fmt"
	"os"
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

	env := godotenv.Load("..\\src\\goapp1\\.env")
	if env != nil {
		fmt.Println(env)
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

// GetDbConn will return the global db connection variable
func GetDbConn() *gorm.DB {
	return db
}
