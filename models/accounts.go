package models

import (
	"errors"
	"fmt"
	"goapp1/util"

	"github.com/jinzhu/gorm"
)

type User struct {
	UserID   uint   `gorm:"primary_key;column:UserID; AUTO_INCREMENT"`
	Username string `gorm:"column:Username"`
	Password string `gorm:"column:Password"`
}

type Userdetail struct {
	UserDetailID int
	DOB          string
	Email        string
	PhoneNumber  string
	UserID       string
}

// AddUser will add a user to the database
func AddUser(newUser User) (bool, error) {
	if exists, existingUser := ValidUser(newUser.Username); exists != false {
		fmt.Println(existingUser)
		return false, util.ErrUserExists
	}
	if len(newUser.Password) < 8 {
		return false, util.ErrPasswordInvalidFormat
	}

	err := GetDbConn().Create(&newUser).Error

	if err != nil && err != gorm.ErrInvalidSQL {
		fmt.Println("invalid syntax")
		return false, errors.New("Invalid syntax")
	}

	return true, nil

}

// ValidUser checks if user exists TODO returns bool if it exists and the user struct
func ValidUser(currUsername string) (bool, User) {
	user := User{}
	db.LogMode(true)
	db := GetDbConn()
	err := db.Table("user").Where("Username=?", currUsername).First(&user).Error

	// err if successful will always return nil
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		fmt.Println("Error occured connecting to DB")
		return false, user
	} else if !gorm.IsRecordNotFoundError(err) { // record is found
		return true, user
	} else {
		return false, user // returns false if record is not found ..... and anything else bad that happens
	}
}

/*
func UserID(username string) int {
	user := User{}
	db := GetDbConn()
	err := db.Where("Username = ?", username).First(&user)

	fmt.Printf("UserID: %d", user.UserID)
	return 0
}
*/

// TODO checkDbErr checks db errors ..... maybe some logging service on the way
// possible gorm errors that can be returned ErrUnaddressable, ErrInvalidSQL, ErrRecordNotFound
func checkDbErr() {
}

// notes
// with gorm to put more tags on struct ses a semi colon ie gorm:blah:soemthing;key:bleh
// useful querries
//     db.Table("users").Where("Username=?", currUsername).Find(&user)
//     db.Table("users").Where("Username=?", currUsername).Select("*")
//     db.Raw("SELECT * FROM users WHERE Username = ?", currUsername)
