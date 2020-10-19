package main

import (
	"github.com/Robert1138/GWA/pkg/api"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	api.StartAPI()
}
