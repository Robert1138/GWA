package account

import (
	"github.com/Robert1138/GWA/util"
	"github.com/Robert1138/GWA/util/db"

	"github.com/jinzhu/gorm"
)

func UpdatePassword(userID int, newPassword string) error {
	var err error
	newPassword, err = util.HashPassword(newPassword)
	if err != nil {
		return err
	}
	err = db.GetDbConn().Table("user").Where("UserID=?", userID).Update("UserPassword", newPassword).Error
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		log.Error("no record found for password change UpdatePassword accountmodel.go")
		return err
	} else if !gorm.IsRecordNotFoundError(err) {
		return nil
	} else {
		return err
	}
}
