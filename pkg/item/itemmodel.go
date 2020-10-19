package item

import (
	"fmt"

	"github.com/Robert1138/GWA/util/db"

	"github.com/jinzhu/gorm"
	//"goapp1/util/db"
)

type Item struct {
	ItemID         int    `gorm:"primary_key;column:ItemID; AUTO_INCREMENT"`
	ItemTitle      string `gorm:"column:ItemTitle"`
	ItemSubtitle   string `gorm:"column:ItemSubtitle"`
	ItemContent    string `gorm:"column:ItemContent"`
	ItemCurrentBid int    `gorm:"column:ItemCurrentBid"`
	ItemTotalBids  int    `gorm:"column:ItemTotalBids"`
	ItemStartDate  uint   `gorm:"column:ItemStartDate"`
	ItemEndDate    uint   `gorm:"column:ItemEndDate"`
	ItemHasReserve int    `gorm:"column:ItemHasReserve"`
	ItemReserve    int    `gorm:"column:ItemReserve"`
}

// GetItem returns an Item struct for the corresponding item ID.  If the item does not exist it will return empty or an error
func GetItem(itemID int) Item {
	item := Item{}
	err := db.GetDbConn().Table("item").Where("ItemID=?", itemID).First(&item).Error
	fmt.Println("getitem")

	if err != nil && !gorm.IsRecordNotFoundError(err) {
		fmt.Println("couldn't connect to db")
	} else if !gorm.IsRecordNotFoundError(err) {
		fmt.Println("found item")
		return item
	} else {
		fmt.Println("Something really went wrong")
	}
	return item
}

// at the moment this will just return a small list of items
func ItemList() {

}
