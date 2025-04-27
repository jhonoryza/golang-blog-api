package controller

import (
	"api_blog/exception"
	"api_blog/requests"
	"api_blog/response"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/hablullah/go-hijri"
	"github.com/hablullah/go-prayer"
	"github.com/julienschmidt/httprouter"
)

type PrayController struct {
	DB *sql.DB
}

func NewPrayController(db *sql.DB) *PrayController {
	return &PrayController{
		DB: db,
	}
}

type Schedule struct {
	Date    string
	Fajr    string
	Sunrise string
	Zuhr    string
	Asr     string
	Maghrib string
	Isha    string
}

func (c *PrayController) Index(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var req requests.PrayerRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		exception.BadRequestError(w, r, err)
		return
	}

	validate := requests.NewPrayerValidator()
	err = validate.Struct(req)
	if err != nil {
		exception.ValidationError(w, r, err)
		return
	}

	timeZone, _ := time.LoadLocation(req.Timezone)
	schedules, _ := prayer.Calculate(prayer.Config{
		Latitude:           req.Latitude,
		Longitude:          req.Longitude,
		Timezone:           timeZone,
		TwilightConvention: prayer.Kemenag(),
		AsrConvention:      prayer.Shafii,
		PreciseToSeconds:   true,
	}, req.Year)

	var data []Schedule

	for _, v := range schedules {
		temp := Schedule{
			Date:    v.Date,
			Fajr:    v.Fajr.Format("2006-06-02 15:04:05"),
			Sunrise: v.Sunrise.Format("2006-06-02 15:04:05"),
			Zuhr:    v.Zuhr.Format("2006-06-02 15:04:05"),
			Asr:     v.Asr.Format("2006-06-02 15:04:05"),
			Maghrib: v.Maghrib.Format("2006-06-02 15:04:05"),
			Isha:    v.Isha.Format("2006-06-02 15:04:05"),
		}
		data = append(data, temp)
	}

	var resp *response.ApiResponse
	if req.Raw == true {
		resp = &response.ApiResponse{
			Code:    http.StatusOK,
			Message: "OK",
			Data:    schedules,
		}

	} else {
		resp = &response.ApiResponse{
			Code:    http.StatusOK,
			Message: "OK",
			Data:    data,
		}
	}

	resp.ToJson(w)
}

func (c *PrayController) HijriCalendar(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	hijriDate, _ := hijri.CreateHijriDate(time.Now(), hijri.Default)
	data := fmt.Sprintf("%v %v %v Hijriyah", hijriDate.Day, getHijriMonthName(hijriDate.Month), hijriDate.Year)
	resp := &response.ApiResponse{
		Code:    http.StatusOK,
		Message: "OK",
		Data:    data,
	}
	resp.ToJson(w)
}

func getHijriMonthName(month int64) string {
	months := []string{
		"Muharram", "Safar", "Rabi'ul Awwal", "Rabi'ul Akhir",
		"Jumada'ul Awwal", "Jumada'ul Akhir", "Rajab", "Sha'ban",
		"Ramadhan", "Shawwal", "Dhul Qa'dah", "Dhul Hijjah",
	}
	return months[month-1]
}
