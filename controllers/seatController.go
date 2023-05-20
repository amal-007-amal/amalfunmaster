package controllers

import (
	"errors"
	"net/http"
	"sample/server/database"
	"sample/server/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SeatController struct{}

type (
	GetSeatPricingResponse struct {
		SeatId     string
		Class      string
		IsBooked   bool
		ClassPrice float64
	}
	GetAllSeatsResponse struct {
		SeatId   string
		Class    string
		IsBooked bool
	}
)

// handler load all the seats informations
func (cntrl SeatController) GetAllSeats(ctx *gin.Context) {
	var seats []models.Seat
	result := database.SqlDB.Preload("SeatClass").Find(&seats)
	if result.Error != nil {
		ctx.Status(http.StatusInternalServerError)
		ctx.Abort()
		return
	}
	var response []GetAllSeatsResponse
	for _, seat := range seats {
		response = append(response, GetAllSeatsResponse{seat.ID, seat.SeatClass.Title, seat.BookingID.Valid})
	}
	ctx.JSON(http.StatusOK, response)
}

// load the seatpricing
func (cntrl SeatController) GetSeatPricing(ctx *gin.Context) {
	if ctx.Param("id") != "" {
		var seat models.Seat
		result := database.SqlDB.Preload("SeatClass").First(&seat, "id = ?", ctx.Param("id"))
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				ctx.Status(http.StatusNotFound)
				ctx.Abort()
				return
			}
			ctx.Status(http.StatusInternalServerError)
			ctx.Abort()
			return
		}
		response := new(GetSeatPricingResponse)
		response.SeatId = seat.ID
		response.Class = seat.SeatClass.Title
		response.ClassPrice = calcClassPrice(seat.SeatClass)
		response.IsBooked = seat.BookingID.Valid
		ctx.JSON(http.StatusOK, response)
	} else {
		ctx.Status(http.StatusBadRequest)
		ctx.Abort()
	}
}

func calcClassPrice(seatClass models.SeatClass) float64 {
	ratio := float64(seatClass.BookedCount) / float64(seatClass.SeatCount)
	if (ratio < 0.4) && seatClass.MinPrice.Valid {
		return seatClass.MinPrice.Float64
	} else if (ratio < 0.6) && seatClass.NormalPrice.Valid {
		return seatClass.NormalPrice.Float64
	} else if seatClass.MaxPrice.Valid {
		return seatClass.MaxPrice.Float64
	} else if seatClass.NormalPrice.Valid {
		return seatClass.NormalPrice.Float64
	} else if seatClass.MinPrice.Valid {
		return seatClass.MinPrice.Float64
	}
	return 0
}
