package db

import (
	"os"
	"sync"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var GormDB *gorm.DB
var Mut *sync.Mutex

func init() {
	db_name := os.Getenv("GORM_SQLITE_NAME")
	var err error
	GormDB, err = gorm.Open(sqlite.Open(db_name), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	Mut = &sync.Mutex{}
}
