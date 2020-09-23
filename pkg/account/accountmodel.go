package account

import (
	"goapp1/util/db"

	"github.com/jinzhu/gorm"
)

func UpdatePassword(userID int, newPassword string) bool {
	err := db.GetDbConn().Table("user").Where("UserID=?", userID).Update("UserPassword", newPassword).Error
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		log.Error("no record found for password change UpdatePassword accountmodel.go")
		return false
	} else if !gorm.IsRecordNotFoundError(err) {
		return true
	} else {
		return false
	}
}
