package controller

import (
"api_blog/domain/entities"
	"api_blog/infrastructure/cache"
	"api_blog/infrastructure/exception"
	request "api_blog/delivery/http/request"
	"api_blog/delivery/http/response"
	"api_blog/usecase/quran"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
)

type QuranController struct {
	Usecase *quran.QuranUsecase
	Cache   *cache.CacheManager
}

func NewQuranController(uc *quran.QuranUsecase, cacheMgr *cache.CacheManager) *QuranController {
	return &QuranController{
		Usecase: uc,
		Cache:   cacheMgr,
	}
}

func (c *QuranController) Index(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	input := quran.FindAllSurahInput{
		Search:  r.URL.Query().Get("search"),
		SortBy:  r.URL.Query().Get("sortBy"),
		SortDir: r.URL.Query().Get("sortDir"),
	}

	if input.SortDir == "" {
		input.SortDir = "asc"
	}
	if input.SortBy == "" {
		input.SortBy = "external_id"
	}

	ctx := r.Context()

	if !cache.HasSearchOrSortParams(r) {
		cacheKey := cache.BuildSurahCacheKey(r)
		var cachedData []entities.Quran
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

	surah := c.Usecase.FindAllSurah(r.Context(), input)

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

	req := request.QuranRequest{
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

	if !cache.HasSearchOrSortParams(r) && surahId != "" {
		cacheKey := cache.BuildAyahCacheKey(r, surahId)
		var cachedData []entities.QuranVerse
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

	input := quran.FindAllAyahInput{
		SurahId: surahId,
		Search:  search,
		SortBy:  sortBy,
		SortDir: sortDir,
	}

	ayah := c.Usecase.FindAllAyahFromSurah(r.Context(), input)

	if !cache.HasSearchOrSortParams(r) && surahId != "" {
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