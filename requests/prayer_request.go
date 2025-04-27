package requests

import "github.com/go-playground/validator/v10"

type PrayerRequest struct {
	Timezone  string  `json:"timezone" validate:"required,min=1,max=255"`
	Year      int     `json:"year" validate:"required,yearValid"`
	Longitude float64 `json:"longitude" validate:"required,longitudeValid"`
	Latitude  float64 `json:"latitude" validate:"required,latitudeValid"`
	Raw       bool    `json:"raw" validate:"boolean"`
}

func NewPrayerValidator() *validator.Validate {
	validate := validator.New()
	validate.RegisterValidation("yearValid", yearValidator)
	validate.RegisterValidation("longitudeValid", longitudeValidator)
	validate.RegisterValidation("latitudeValid", latitudeValidator)

	return validate
}

// Validasi custom untuk tahun
func yearValidator(fl validator.FieldLevel) bool {
	year, ok := fl.Field().Interface().(int)
	if !ok {
		return false
	}
	return year >= 1900 && year <= 9999
}

// Validasi custom untuk longitude
func longitudeValidator(fl validator.FieldLevel) bool {
	longitude, ok := fl.Field().Interface().(float64)
	if !ok {
		return false
	}
	return longitude >= -180 && longitude <= 180
}

// Validasi custom untuk latitude
func latitudeValidator(fl validator.FieldLevel) bool {
	latitude, ok := fl.Field().Interface().(float64)
	if !ok {
		return false
	}
	return latitude >= -90 && latitude <= 90
}
