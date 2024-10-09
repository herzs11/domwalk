package db

import (
	"fmt"
	"os"
	"sync"

	"github.com/fatih/color"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DomwalkDB struct {
	*gorm.DB
	DBName string
}

var GormDB DomwalkDB
var Mut *sync.Mutex

func GormDBConnect(db_name string) error {
	if _, err := os.Stat(db_name); os.IsNotExist(err) {
		color.Yellow("Database file does not exist. Creating...\n")
	}
	Mut = &sync.Mutex{}
	gdb, err := gorm.Open(
		sqlite.Open(db_name), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Error),
		},
	)
	GormDB = DomwalkDB{DB: gdb, DBName: db_name}
	if err != nil {
		return fmt.Errorf("failed to connect database: %s", err)
	}
	return nil
}
