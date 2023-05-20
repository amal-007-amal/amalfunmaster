package main

import (
	"log"
	"sample/server/controllers"
	"sample/server/database"

	"github.com/gin-gonic/gin"
)

func handleRequests() {
	seatController := new(controllers.SeatController)
	bookingController := new(controllers.BookingController)
	adminController := new(controllers.AdminController)

	router := gin.Default()
	// handling cors error
	// router.Use(cors.New(cors.Config{
	// 	AllowOrigins:     []string{"*"},
	// 	AllowMethods:     []string{"*"},
	// 	AllowHeaders:     []string{"Origin", "Set-Cookies"},
	// 	ExposeHeaders:    []string{"Content-Length"},
	// 	AllowCredentials: true,
	// 	MaxAge:           12 * time.Hour,
	// }))

	router.GET("/seats/", seatController.GetAllSeats)
	router.GET("/seats/:id", seatController.GetSeatPricing)
	router.POST("/booking/", bookingController.BookSeats)
	router.GET("/bookings/", bookingController.GetUserBookings)
	router.GET("/admin/resetData", adminController.ResetData)
	log.Fatal(router.Run(":8000"))
}

func main() {
	database.Init("root:mynewpassword@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local")
	handleRequests()
	defer database.Close()
}
