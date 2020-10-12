package User

// This contains structs and methods for the User model.  It is not a handler package

type User struct {
	UserID   uint   `gorm:"primary_key;column:UserID; AUTO_INCREMENT"`
	UserName string `gorm:"column:UserName"`
	Password string `gorm:"column:UserPassword"`
}
