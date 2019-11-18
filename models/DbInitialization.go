package models

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/joho/godotenv"
)

var db *gorm.DB

func init() {
	env := godotenv.Load("..\\src\\goapp1\\.env")
	if env != nil {
		fmt.Println(env)
	}

	/*
		Apparently we do not close db everytime we get it.
		The db object lives for the life of the program.
		However the following needs to be set for safety.

		db.DB().SetConnMaxLifetime(time.Minute * whatever is determined)
		db.DB().SetMaxIdleConns()
		db.DB().SetMaxOpenConns()
	*/
	dbName := os.Getenv("db_name")
	dbPassword := os.Getenv("db_pass")
	dbHost := os.Getenv("db_host")
	dbUser := os.Getenv("db_user")

	// mysql connection
	dbString := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true&loc=Local", dbUser, dbPassword, dbHost, dbName)
	fmt.Println(dbString)

	connection, e := gorm.Open("mysql", dbString)
	if e != nil {
		fmt.Println(e)
	}

	db = connection

	// connection.Close()

}

func DbTest() {
	fmt.Println("Db init")
}

// GetDbConn() will return the global db connection variable
func GetDbConn() *gorm.DB {
	return db
}
