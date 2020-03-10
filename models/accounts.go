package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

type Account struct {
	gorm.Model
	Id       int
	Name     string
	Age      int
	FavColor string
}

func (account *Account) GetInfo() {
	record := GetDbConn()
	if record == nil {
		fmt.Println("Issue getting connection from account")
	} else {
		fmt.Println("Db connection established using account receiver")
	}
}

/*
Id
Name
Age
FavColor
*/
