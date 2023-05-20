package database

import (
	"fmt"
	"log"
	"sample/server/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// global variable used to db operations
var (
	SqlDB *gorm.DB
)

// this is function connect the db and also close the connection
func Init(sqlUrl string) {
	fmt.Println("mysql test table")
	var err error
	SqlDB, err = gorm.Open(mysql.Open(sqlUrl), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}
	SqlDB.AutoMigrate(&models.SeatClass{})
	SqlDB.AutoMigrate(&models.Booking{})
	SqlDB.AutoMigrate(&models.Seat{})
}

func Close() {
	//
}
