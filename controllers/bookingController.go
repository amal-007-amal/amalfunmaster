package controllers

import (
	"net/http"
	"sample/server/database"
	"sample/server/models"

	"github.com/gin-gonic/gin"
)

type BookingController struct{}

type BookSeatsRequest struct {
	SeatIds []string `binding:"required"`
	Name    string   `binding:"required"`
	Phone   string   `binding:"required"`
	Email   string   `binding:"required"`
}

type BookSeatsResponse struct {
	BookingId int
	Amount    float64
}

func (cntrl BookingController) BookSeats(ctx *gin.Context) {
	var request BookSeatsRequest
	err := ctx.BindJSON(&request)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		ctx.Abort()
	} else {
		var seats []models.Seat
		result := database.SqlDB.Preload("SeatClass").Find(&seats, "id in ?", request.SeatIds)
		if result.Error != nil {
			ctx.Status(http.StatusInternalServerError)
			ctx.Abort()
			return
		}
		if len(seats) != len(request.SeatIds) {
			ctx.Status(http.StatusBadRequest)
			ctx.Abort()
			return
		}
		for _, seat := range seats {
			if seat.BookingID.Valid {
				ctx.Status(http.StatusConflict)
				ctx.Abort()
				return
			}
		}
		booking := models.Booking{Name: request.Name, Email: request.Email, Phone: request.Phone, Seats: seats}
		result = database.SqlDB.Create(&booking)
		if result.Error != nil {
			ctx.Status(http.StatusInternalServerError)
			ctx.Abort()
			return
		}
		totalAmount := 0.0
		for _, seat := range seats {
			totalAmount = totalAmount + (calcClassPrice(seat.SeatClass))
		}

		for _, seat := range seats {
			database.SqlDB.Exec("UPDATE seat_classes SET booked_count = booked_count + 1 WHERE seat_classes.id = ?", seat.SeatClassID)
		}

		response := BookSeatsResponse{BookingId: int(booking.ID), Amount: totalAmount}
		ctx.JSON(http.StatusOK, response)

		//Call to texttosms api
	}
}

type GetUserBookingsResponse struct {
	BookingId uint
}

func (cntrl BookingController) GetUserBookings(ctx *gin.Context) {
	if ctx.Query("id") != "" {
		var bookings []models.Booking
		result := database.SqlDB.Where("email = ? OR phone = ?", ctx.Query("id"), ctx.Query("id")).Find(&bookings)
		if result.Error != nil {
			ctx.Status(http.StatusInternalServerError)
			ctx.Abort()
			return
		}
		var response []GetUserBookingsResponse
		for _, booking := range bookings {
			response = append(response, GetUserBookingsResponse{booking.ID})
		}
		ctx.JSON(http.StatusOK, response)
	} else {
		ctx.Status(http.StatusBadRequest)
		ctx.Abort()
	}
}
