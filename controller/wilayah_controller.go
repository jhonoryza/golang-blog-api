package controller

import (
	"api_blog/repository"
	"api_blog/response"
	"context"
	"database/sql"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type WilayahController struct {
	DB *sql.DB
}

func NewWilayahController(db *sql.DB) *WilayahController {
	return &WilayahController{
		DB: db,
	}
}

func (c *WilayahController) Provinces(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	search := r.URL.Query().Get("search")
	sortBy := r.URL.Query().Get("sortBy")
	sortDir := r.URL.Query().Get("sortDir")

	// add parameter to context
	ctx := r.Context()
	ctx = context.WithValue(ctx, "search", search)
	ctx = context.WithValue(ctx, "sortBy", sortBy)
	ctx = context.WithValue(ctx, "sortDir", sortDir)
	r = r.WithContext(ctx)

	repo := repository.NewWilayahRepository(c.DB)
	data := repo.Provinces(r.Context())

	resp := response.ApiResponse{
		Code:    http.StatusOK,
		Message: "OK",
		Data:    data,
	}
	resp.ToJson(w)
}

func (c *WilayahController) Cities(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	search := r.URL.Query().Get("search")
	sortBy := r.URL.Query().Get("sortBy")
	sortDir := r.URL.Query().Get("sortDir")
	provinceId := r.URL.Query().Get("provinceId")

	// add parameter to context
	ctx := r.Context()
	ctx = context.WithValue(ctx, "search", search)
	ctx = context.WithValue(ctx, "sortBy", sortBy)
	ctx = context.WithValue(ctx, "sortDir", sortDir)
	ctx = context.WithValue(ctx, "provinceId", provinceId)
	r = r.WithContext(ctx)

	repo := repository.NewWilayahRepository(c.DB)
	data := repo.Cities(r.Context())

	resp := response.ApiResponse{
		Code:    http.StatusOK,
		Message: "OK",
		Data:    data,
	}
	resp.ToJson(w)
}

func (c *WilayahController) Districts(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	search := r.URL.Query().Get("search")
	sortBy := r.URL.Query().Get("sortBy")
	sortDir := r.URL.Query().Get("sortDir")
	cityId := r.URL.Query().Get("cityId")

	// add parameter to context
	ctx := r.Context()
	ctx = context.WithValue(ctx, "search", search)
	ctx = context.WithValue(ctx, "sortBy", sortBy)
	ctx = context.WithValue(ctx, "sortDir", sortDir)
	ctx = context.WithValue(ctx, "cityId", cityId)
	r = r.WithContext(ctx)

	repo := repository.NewWilayahRepository(c.DB)
	data := repo.Districts(r.Context())

	resp := response.ApiResponse{
		Code:    http.StatusOK,
		Message: "OK",
		Data:    data,
	}
	resp.ToJson(w)
}

func (c *WilayahController) SubDistricts(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	search := r.URL.Query().Get("search")
	sortBy := r.URL.Query().Get("sortBy")
	sortDir := r.URL.Query().Get("sortDir")
	districtId := r.URL.Query().Get("districtId")

	// add parameter to context
	ctx := r.Context()
	ctx = context.WithValue(ctx, "search", search)
	ctx = context.WithValue(ctx, "sortBy", sortBy)
	ctx = context.WithValue(ctx, "sortDir", sortDir)
	ctx = context.WithValue(ctx, "districtId", districtId)
	r = r.WithContext(ctx)

	repo := repository.NewWilayahRepository(c.DB)
	data := repo.SubDistricts(r.Context())

	resp := response.ApiResponse{
		Code:    http.StatusOK,
		Message: "OK",
		Data:    data,
	}
	resp.ToJson(w)
}
