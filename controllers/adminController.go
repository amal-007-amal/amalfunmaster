package controllers

import (
	"database/sql"
	"encoding/csv"
	"net/http"
	"os"
	"sample/server/database"
	"sample/server/models"
	"strings"

	"github.com/gin-gonic/gin"
)

type AdminController struct{}

func (cntrl AdminController) ResetData(ctx *gin.Context) {
	SEATS_FILE_PATH := "/home/user/Downloads/AmalFlurn-master/database/MockData/Seats.csv"
	SEAT_CLASSES_FILE_PATH := "/home/user/Downloads/AmalFlurn-master/database/MockData/SeatPricing.csv"
	database.SqlDB.Exec("SET FOREIGN_KEY_CHECKS = 0;")
	database.SqlDB.Exec("TRUNCATE TABLE seats;")
	database.SqlDB.Exec("TRUNCATE TABLE bookings;")
	database.SqlDB.Exec("TRUNCATE TABLE seat_classes;")
	database.SqlDB.Exec("SET FOREIGN_KEY_CHECKS = 1;")

	classesFile, _ := os.Open(SEAT_CLASSES_FILE_PATH)
	seatsFile, _ := os.Open(SEATS_FILE_PATH)
	defer classesFile.Close()
	defer seatsFile.Close()

	csvReader := csv.NewReader(classesFile)
	rows, _ := csvReader.ReadAll()
	var seatClasses []models.SeatClass
	for _, row := range rows[1:] {
		var seatClass models.SeatClass
		seatClass.Title = row[1]
		fetchPrice(&seatClass.MinPrice, row[2])
		fetchPrice(&seatClass.NormalPrice, row[3])
		fetchPrice(&seatClass.MaxPrice, row[4])
		seatClasses = append(seatClasses, seatClass)
	}
	database.SqlDB.Create(&seatClasses)

	csvReader = csv.NewReader(seatsFile)
	rows, _ = csvReader.ReadAll()
	var seats []models.Seat
	for _, row := range rows[1:] {
		var seatClass models.SeatClass
		var seat models.Seat
		seat.ID = row[1]
		database.SqlDB.Find(&seatClass, "title = ?", row[2])
		seat.SeatClassID = int(seatClass.ID)
		seats = append(seats, seat)
	}
	database.SqlDB.Create(&seats)
	database.SqlDB.Exec("UPDATE seat_classes SET seat_count = (SELECT COUNT(*) FROM seats WHERE seats.seat_class_id = seat_classes.id);")
	ctx.Status(http.StatusNoContent)
}

func fetchPrice(val *sql.NullFloat64, str string) {
	ts := strings.Trim(str, "$")
	if ts == "" {
		val.Scan(nil)
	} else {
		val.Scan(ts)
	}
}
