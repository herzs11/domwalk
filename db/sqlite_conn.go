package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var GormDB *gorm.DB

func init() {
	var err error
	GormDB, err = gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
}
