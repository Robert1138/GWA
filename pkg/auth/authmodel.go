package auth

import (
	"errors"
	"fmt"
	u "goapp1/pkg/user"
	"goapp1/util"
	"goapp1/util/db"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	UserID   uint   `gorm:"primary_key;column:UserID; AUTO_INCREMENT"`
	UserName string `gorm:"column:UserName"`
	Password string `gorm:"column:UserPassword"`
}

// Account represents a Users account details
type Account struct {
	AccountID          int
	AccountFirstName   string
	AccountLastName    string
	AccountPhoneNumber string
	AccountEmail       string
	AccountUserID      string
}

// AddUser will add a user to the database
func AddUser(newUser u.User) (bool, error) {
	if exists, existingUser := ValidUser(newUser.UserName); exists != false {
		fmt.Println(existingUser)
		return false, util.ErrUserExists
	}
	if len(newUser.Password) < 8 {
		return false, util.ErrPasswordInvalidFormat
	}
	var err error // gross go through this and clean it up
	newUser.Password, err = hashPassword(newUser.Password)
	if err != nil {
		fmt.Println("password couldnt be encrypted")
	}

	err = db.GetDbConn().Create(&newUser).Error

	if err != nil && err != gorm.ErrInvalidSQL {
		fmt.Println("invalid syntax")
		return false, errors.New("Invalid syntax")
	}

	return true, nil

}

// ValidUser checks if user exists TODO returns bool if it exists and the user struct
func ValidUser(currUserName string) (bool, u.User) {
	user := u.User{}
	db.GetDbConn().LogMode(true)

	err := db.GetDbConn().Table("user").Where("UserName=?", currUserName).First(&user).Error

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

// not used, determining if needed
func verifyUser(currUserName string, currUserPassword string) (bool, u.User) {
	user := u.User{}
	db.GetDbConn().LogMode(true)

	err := db.GetDbConn().Table("user").Where("UserName=? AND UserPassword=?", currUserName, currUserPassword).First(&user).Error

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

// ValidUserPassword checks if user exists and if the password matches.  Returns a User struct and a bool if user pass matched
func ValidUserPassword(currUserName string, password string) (u.User, bool) {
	user := u.User{}
	db.GetDbConn().LogMode(true)

	err := db.GetDbConn().Table("user").Where("UserName=?", currUserName).First(&user).Error

	// err if successful will always return nil
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		fmt.Println("Error occured connecting to DB")
		return u.User{}, false
	} else if !gorm.IsRecordNotFoundError(err) { // record is found
		if checkPassword(user.Password, password) { // only if this is true will this func return a populated User struct
			fmt.Println("user pass matches")
			return user, true
		}
		fmt.Println("user pass does not match")

		return u.User{}, false
	} else {
		return u.User{}, false // returns false if record is not found ..... and anything else bad that happens
	}
}

// the cost has profound effect on Response time
// cost of 12 adds over 500ms, 14 1600ms to 1900ms
func hashPassword(password string) (string, error) {
	passBytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(passBytes), err
}

func checkPassword(hashedPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

// notes
// with gorm to put more tags on struct ses a semi colon ie gorm:blah:soemthing;key:bleh
// useful querries
//     db.Table("users").Where("UserName=?", currUserName).Find(&user)
//     db.Table("users").Where("UserName=?", currUserName).Select("*")
//     db.Raw("SELECT * FROM users WHERE UserName = ?", currUserName)
