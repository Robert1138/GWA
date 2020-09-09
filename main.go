package main

import (
	"goapp1/pkg/api"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	api.StartAPI()
}
