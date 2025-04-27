package controller

import (
	"api_blog/exception"
	"api_blog/repository"
	"api_blog/requests"
	"api_blog/response"
	"context"
	"database/sql"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
)

type QuranController struct {
	DB *sql.DB
}

func NewQuranController(db *sql.DB) *QuranController {
	return &QuranController{
		DB: db,
	}
}

func (c *QuranController) Index(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	search := r.URL.Query().Get("search")
	sortBy := r.URL.Query().Get("sortBy")
	sortDir := r.URL.Query().Get("sortDir")

	if sortDir == "" {
		sortDir = "asc"
	}
	if sortBy == "" {
		sortBy = "external_id"
	}
	// add parameter to context
	ctx := r.Context()
	ctx = context.WithValue(ctx, "search", search)
	ctx = context.WithValue(ctx, "sortBy", sortBy)
	ctx = context.WithValue(ctx, "sortDir", sortDir)
	r = r.WithContext(ctx)

	// get all tools
	quranRepo := repository.NewQuranRepository(c.DB)
	surah := quranRepo.FindAllSurah(r.Context())

	resp := response.ApiResponse{
		Code:    http.StatusOK,
		Message: "OK",
		Data:    surah,
	}

	resp.ToJson(w)
}

func (c *QuranController) Show(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	search := r.URL.Query().Get("search")
	sortBy := r.URL.Query().Get("sortBy")
	sortDir := r.URL.Query().Get("sortDir")
	surahId := r.URL.Query().Get("surahId")

	req := requests.QuranRequest{
		SurahId: surahId,
	}

	validate := validator.New()

	if err := validate.Struct(req); err != nil {
		exception.ErrorHandler(w, r, err)
		return
	}

	if sortDir == "" {
		sortDir = "asc"
	}
	if sortBy == "" {
		sortBy = "id"
	}
	// add parameter to context
	ctx := r.Context()
	ctx = context.WithValue(ctx, "search", search)
	ctx = context.WithValue(ctx, "sortBy", sortBy)
	ctx = context.WithValue(ctx, "sortDir", sortDir)
	ctx = context.WithValue(ctx, "surahId", surahId)
	r = r.WithContext(ctx)

	// get all tools
	quranRepo := repository.NewQuranRepository(c.DB)
	ayah := quranRepo.FindAllAyahFromSurah(r.Context())

	resp := response.ApiResponse{
		Code:    http.StatusOK,
		Message: "OK",
		Data:    ayah,
	}

	resp.ToJson(w)
}
