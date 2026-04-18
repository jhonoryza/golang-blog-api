package controller

import (
	"api_blog/delivery/http/response"
	"api_blog/usecase/wilayah"
	"database/sql"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type WilayahController struct {
	Usecase *wilayah.WilayahUsecase
	DB      *sql.DB
}

func NewWilayahController(uc *wilayah.WilayahUsecase, db *sql.DB) *WilayahController {
	return &WilayahController{
		Usecase: uc,
		DB:      db,
	}
}

func (c *WilayahController) Provinces(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	input := wilayah.FindAllProvincesInput{
		Search:  r.URL.Query().Get("search"),
		SortBy:  r.URL.Query().Get("sortBy"),
		SortDir: r.URL.Query().Get("sortDir"),
	}

	data := c.Usecase.FindAllProvinces(r.Context(), input)

	resp := response.ApiResponse{
		Code:    http.StatusOK,
		Message: "OK",
		Data:    data,
	}
	resp.ToJson(w)
}

func (c *WilayahController) Cities(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	input := wilayah.FindCitiesByProvinceInput{
		ProvinceId: r.URL.Query().Get("provinceId"),
		Search:     r.URL.Query().Get("search"),
		SortBy:     r.URL.Query().Get("sortBy"),
		SortDir:    r.URL.Query().Get("sortDir"),
	}

	data := c.Usecase.FindCitiesByProvince(r.Context(), input)

	resp := response.ApiResponse{
		Code:    http.StatusOK,
		Message: "OK",
		Data:    data,
	}
	resp.ToJson(w)
}

func (c *WilayahController) Districts(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	input := wilayah.FindDistrictsByCityInput{
		CityId:  r.URL.Query().Get("cityId"),
		Search:  r.URL.Query().Get("search"),
		SortBy:  r.URL.Query().Get("sortBy"),
		SortDir: r.URL.Query().Get("sortDir"),
	}

	data := c.Usecase.FindDistrictsByCity(r.Context(), input)

	resp := response.ApiResponse{
		Code:    http.StatusOK,
		Message: "OK",
		Data:    data,
	}
	resp.ToJson(w)
}

func (c *WilayahController) SubDistricts(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	input := wilayah.FindSubdistrictsByDistrictInput{
		DistrictId: r.URL.Query().Get("districtId"),
		Search:     r.URL.Query().Get("search"),
		SortBy:     r.URL.Query().Get("sortBy"),
		SortDir:    r.URL.Query().Get("sortDir"),
	}

	data := c.Usecase.FindSubdistrictsByDistrict(r.Context(), input)

	resp := response.ApiResponse{
		Code:    http.StatusOK,
		Message: "OK",
		Data:    data,
	}
	resp.ToJson(w)
}