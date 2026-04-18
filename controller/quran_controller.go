package controller

import (
	"api_blog/cache"
	"api_blog/entity"
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
	DB    *sql.DB
	Cache *cache.CacheManager
}

func NewQuranController(db *sql.DB, cacheMgr *cache.CacheManager) *QuranController {
	return &QuranController{
		DB:    db,
		Cache: cacheMgr,
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

	ctx := r.Context()

	if !cache.HasSearchOrSortParams(r) {
		cacheKey := cache.BuildSurahCacheKey(r)
		var cachedData []entity.Quran
		if err := c.Cache.Get(ctx, cacheKey, &cachedData); err == nil && cachedData != nil {
			resp := response.ApiResponse{
				Code:    http.StatusOK,
				Message: "OK",
				Data:    cachedData,
			}
			resp.ToJson(w)
			return
		}
	}

	ctx = context.WithValue(ctx, "search", search)
	ctx = context.WithValue(ctx, "sortBy", sortBy)
	ctx = context.WithValue(ctx, "sortDir", sortDir)
	r = r.WithContext(ctx)

	quranRepo := repository.NewQuranRepository(c.DB)
	surah := quranRepo.FindAllSurah(r.Context())

	if !cache.HasSearchOrSortParams(r) {
		cacheKey := cache.BuildSurahCacheKey(r)
		c.Cache.Set(ctx, cacheKey, surah)
	}

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

	ctx := r.Context()

	if !cache.HasSearchOrSortParams(r) {
		cacheKey := cache.BuildAyahCacheKey(r, surahId)
		var cachedData []entity.QuranVerse
		if err := c.Cache.Get(ctx, cacheKey, &cachedData); err == nil && cachedData != nil {
			resp := response.ApiResponse{
				Code:    http.StatusOK,
				Message: "OK",
				Data:    cachedData,
			}
			resp.ToJson(w)
			return
		}
	}

	ctx = context.WithValue(ctx, "search", search)
	ctx = context.WithValue(ctx, "sortBy", sortBy)
	ctx = context.WithValue(ctx, "sortDir", sortDir)
	ctx = context.WithValue(ctx, "surahId", surahId)
	r = r.WithContext(ctx)

	quranRepo := repository.NewQuranRepository(c.DB)
	ayah := quranRepo.FindAllAyahFromSurah(r.Context())

	if !cache.HasSearchOrSortParams(r) {
		cacheKey := cache.BuildAyahCacheKey(r, surahId)
		c.Cache.Set(ctx, cacheKey, ayah)
	}

	resp := response.ApiResponse{
		Code:    http.StatusOK,
		Message: "OK",
		Data:    ayah,
	}

	resp.ToJson(w)
}