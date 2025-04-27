package repository

import (
	"api_blog/entity"
	"api_blog/exception"
	"context"
	"database/sql"
)

type QuranRepository struct {
	DB *sql.DB
}

func NewQuranRepository(db *sql.DB) *QuranRepository {
	return &QuranRepository{DB: db}
}

func (r *QuranRepository) FindAllSurah(ctx context.Context) *[]entity.Quran {
	search := ctx.Value("search").(string)
	sortBy := ctx.Value("sortBy").(string)
	sortDir := ctx.Value("sortDir").(string)

	allowedSortBy := map[string]bool{
		"name": true,
		"id":   true,
	}

	allowedSortDir := map[string]bool{
		"asc":  true,
		"desc": true,
	}

	// Default values kalau input nggak valid
	if !allowedSortBy[sortBy] {
		sortBy = "id"
	}

	if !allowedSortDir[sortDir] {
		sortDir = "desc"
	}

	args := []any{}
	query := `SELECT qurans.* FROM qurans`

	if search != "" {
		query += ` WHERE (lower(transliteration) LIKE lower($1) OR lower(translation) LIKE lower($1))`
		args = append(args, "%"+search+"%")
	}
	query += " ORDER BY " + sortBy + " " + sortDir

	rows, err := r.DB.QueryContext(ctx, query, args...)

	exception.PanicIfErr(err)
	defer rows.Close()

	var surah []entity.Quran
	for rows.Next() {
		var quran entity.Quran
		err := rows.Scan(&quran.ID, &quran.ExternalId, &quran.Arabic, &quran.Latin, &quran.Transliteration, &quran.Translation, &quran.NumAyah, &quran.Page, &quran.Location, &quran.CreatedAt, &quran.UpdatedAt)
		exception.PanicIfErr(err)
		surah = append(surah, quran)
	}

	return &surah
}

func (r *QuranRepository) FindAllAyahFromSurah(ctx context.Context) *[]entity.QuranVerse {
	search := ctx.Value("search").(string)
	sortBy := ctx.Value("sortBy").(string)
	sortDir := ctx.Value("sortDir").(string)
	surahId := ctx.Value("surahId").(string)

	allowedSortBy := map[string]bool{
		"name": true,
		"id":   true,
	}

	allowedSortDir := map[string]bool{
		"asc":  true,
		"desc": true,
	}

	// Default values kalau input nggak valid
	if !allowedSortBy[sortBy] {
		sortBy = "id"
	}

	if !allowedSortDir[sortDir] {
		sortDir = "desc"
	}

	args := []any{}
	query := `SELECT quran_verses.* FROM quran_verses WHERE quran_id = $1`

	args = append(args, surahId)

	if search != "" {
		query += ` AND (lower(transliteration) LIKE lower($2) OR lower(translation) LIKE lower($2))`
		args = append(args, "%"+search+"%")
	}
	query += " ORDER BY " + sortBy + " " + sortDir

	rows, err := r.DB.QueryContext(ctx, query, args...)

	exception.PanicIfErr(err)
	defer rows.Close()

	var ayah []entity.QuranVerse
	for rows.Next() {
		var verse entity.QuranVerse
		err := rows.Scan(&verse.ID, &verse.QuranId, &verse.Ayah, &verse.Page, &verse.Juz, &verse.Arabic, &verse.Kitabah, &verse.Latin, &verse.Translation, &verse.AudioUrl, &verse.CreatedAt, &verse.UpdatedAt)
		exception.PanicIfErr(err)
		ayah = append(ayah, verse)
	}

	return &ayah
}
