package persistence

import (
	"api_blog/domain/entities"
	"api_blog/infrastructure/exception"
	"context"
	"database/sql"
)

type WilayahRepository struct {
	DB *sql.DB
}

func NewWilayahRepository(db *sql.DB) *WilayahRepository {
	return &WilayahRepository{
		DB: db,
	}
}

func (r *WilayahRepository) FindAllProvinces(ctx context.Context) *[]entities.Province {
	search := ctx.Value("search").(string)
	sortBy := ctx.Value("sortBy").(string)
	sortDir := ctx.Value("sortDir").(string)

	allowedSortBy := map[string]bool{
		"id":   true,
		"name": true,
	}

	allowedSortDir := map[string]bool{
		"asc":  true,
		"desc": true,
	}

	if !allowedSortBy[sortBy] {
		sortBy = "name"
	}

	if !allowedSortDir[sortDir] {
		sortDir = "asc"
	}

	query := `select id,name from provinces`

	args := []any{}

	if search != "" {
		query += ` WHERE (lower(id) LIKE lower($1) OR lower(name) LIKE lower($1))`
		args = append(args, "%"+search+"%")
	}
	query += ` order by ` + sortBy + ` ` + sortDir

	var rows *sql.Rows
	var err error

	if len(args) > 0 {
		rows, err = r.DB.QueryContext(ctx, query, args...)
	} else {
		rows, err = r.DB.QueryContext(ctx, query)
	}

	exception.PanicIfErr(err)
	defer rows.Close()

	var provinces []entities.Province
	for rows.Next() {
		var prov entities.Province
		err = rows.Scan(&prov.Id, &prov.Name)
		exception.PanicIfErr(err)
		provinces = append(provinces, prov)
	}
	return &provinces
}

func (r *WilayahRepository) FindCitiesByProvince(ctx context.Context, provinceId string) *[]entities.City {
	search := ctx.Value("search").(string)
	sortBy := ctx.Value("sortBy").(string)
	sortDir := ctx.Value("sortDir").(string)

	allowedSortBy := map[string]bool{
		"id":   true,
		"name": true,
	}

	allowedSortDir := map[string]bool{
		"asc":  true,
		"desc": true,
	}

	if !allowedSortBy[sortBy] {
		sortBy = "name"
	}

	if !allowedSortDir[sortDir] {
		sortDir = "asc"
	}

	query := `select id,name,province_id from cities`

	args := []any{}

	if provinceId != "" {
		query += ` WHERE province_id = $1`
		args = append(args, provinceId)
	}
	if search != "" {
		if provinceId != "" {
			query += ` AND (lower(id) LIKE lower($2) OR lower(name) LIKE lower($2))`
		} else {
			query += ` WHERE (lower(id) LIKE lower($1) OR lower(name) LIKE lower($1))`
		}
		args = append(args, "%"+search+"%")
	}
	query += ` order by ` + sortBy + ` ` + sortDir

	var rows *sql.Rows
	var err error

	if len(args) > 0 {
		rows, err = r.DB.QueryContext(ctx, query, args...)
	} else {
		rows, err = r.DB.QueryContext(ctx, query)
	}

	exception.PanicIfErr(err)
	defer rows.Close()

	var cities []entities.City
	for rows.Next() {
		var city entities.City
		err = rows.Scan(&city.Id, &city.Name, &city.ProvinceId)
		exception.PanicIfErr(err)
		cities = append(cities, city)
	}
	return &cities
}

func (r *WilayahRepository) FindDistrictsByCity(ctx context.Context, cityId string) *[]entities.District {
	search := ctx.Value("search").(string)
	sortBy := ctx.Value("sortBy").(string)
	sortDir := ctx.Value("sortDir").(string)

	allowedSortBy := map[string]bool{
		"id":   true,
		"name": true,
	}

	allowedSortDir := map[string]bool{
		"asc":  true,
		"desc": true,
	}

	if !allowedSortBy[sortBy] {
		sortBy = "name"
	}

	if !allowedSortDir[sortDir] {
		sortDir = "asc"
	}

	query := `select id,name,city_id from districts`

	args := []any{}

	if cityId != "" {
		query += ` WHERE city_id = $1`
		args = append(args, cityId)
	}
	if search != "" {
		if cityId != "" {
			query += ` AND (lower(id) LIKE lower($2) OR lower(name) LIKE lower($2))`
		} else {
			query += ` WHERE (lower(id) LIKE lower($1) OR lower(name) LIKE lower($1))`
		}
		args = append(args, "%"+search+"%")
	}
	query += ` order by ` + sortBy + ` ` + sortDir

	var rows *sql.Rows
	var err error

	if len(args) > 0 {
		rows, err = r.DB.QueryContext(ctx, query, args...)
	} else {
		rows, err = r.DB.QueryContext(ctx, query)
	}

	exception.PanicIfErr(err)
	defer rows.Close()

	var districts []entities.District
	for rows.Next() {
		var district entities.District
		err = rows.Scan(&district.Id, &district.Name, &district.CityId)
		exception.PanicIfErr(err)
		districts = append(districts, district)
	}
	return &districts
}

func (r *WilayahRepository) FindSubdistrictsByDistrict(ctx context.Context, districtId string) *[]entities.Subdistrict {
	search := ctx.Value("search").(string)
	sortBy := ctx.Value("sortBy").(string)
	sortDir := ctx.Value("sortDir").(string)

	allowedSortBy := map[string]bool{
		"id":   true,
		"name": true,
	}

	allowedSortDir := map[string]bool{
		"asc":  true,
		"desc": true,
	}

	if !allowedSortBy[sortBy] {
		sortBy = "name"
	}

	if !allowedSortDir[sortDir] {
		sortDir = "asc"
	}

	query := `select id,name,district_id from subdistricts`

	args := []any{}

	if districtId != "" {
		query += ` WHERE district_id = $1`
		args = append(args, districtId)
	}
	if search != "" {
		if districtId != "" {
			query += ` AND (lower(id) LIKE lower($2) OR lower(name) LIKE lower($2))`
		} else {
			query += ` WHERE (lower(id) LIKE lower($1) OR lower(name) LIKE lower($1))`
		}
		args = append(args, "%"+search+"%")
	}

	query += ` order by ` + sortBy + ` ` + sortDir

	var rows *sql.Rows
	var err error

	if len(args) > 0 {
		rows, err = r.DB.QueryContext(ctx, query, args...)
	} else {
		rows, err = r.DB.QueryContext(ctx, query)
	}

	exception.PanicIfErr(err)
	defer rows.Close()

	var subdistricts []entities.Subdistrict
	for rows.Next() {
		var subdistrict entities.Subdistrict
		err = rows.Scan(&subdistrict.Id, &subdistrict.Name, &subdistrict.DistrictId)
		exception.PanicIfErr(err)
		subdistricts = append(subdistricts, subdistrict)
	}
	return &subdistricts
}
