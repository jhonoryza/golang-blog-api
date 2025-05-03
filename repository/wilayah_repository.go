package repository

import (
	"api_blog/entity"
	"api_blog/exception"
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

func (r *WilayahRepository) Provinces(ctx context.Context) *[]entity.Province {
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

	// Default values kalau input nggak valid
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

	var provinces []entity.Province
	for rows.Next() {
		var prov entity.Province
		err = rows.Scan(&prov.Id, &prov.Name)
		exception.PanicIfErr(err)
		provinces = append(provinces, prov)
	}
	return &provinces
}

func (r *WilayahRepository) Cities(ctx context.Context) *[]entity.City {
	search := ctx.Value("search").(string)
	sortBy := ctx.Value("sortBy").(string)
	sortDir := ctx.Value("sortDir").(string)
	provinceId := ctx.Value("provinceId").(string)

	allowedSortBy := map[string]bool{
		"id":   true,
		"name": true,
	}

	allowedSortDir := map[string]bool{
		"asc":  true,
		"desc": true,
	}

	// Default values kalau input nggak valid
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
		query += ` AND (lower(id) LIKE lower($2) OR lower(name) LIKE lower($2))`
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

	var cities []entity.City
	for rows.Next() {
		var city entity.City
		err = rows.Scan(&city.Id, &city.Name, &city.ProvinceId)
		exception.PanicIfErr(err)
		cities = append(cities, city)
	}
	return &cities
}

func (r *WilayahRepository) Districts(ctx context.Context) *[]entity.District {
	search := ctx.Value("search").(string)
	sortBy := ctx.Value("sortBy").(string)
	sortDir := ctx.Value("sortDir").(string)
	cityId := ctx.Value("cityId").(string)

	allowedSortBy := map[string]bool{
		"id":   true,
		"name": true,
	}

	allowedSortDir := map[string]bool{
		"asc":  true,
		"desc": true,
	}

	// Default values kalau input nggak valid
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
		query += ` AND (lower(id) LIKE lower($2) OR lower(name) LIKE lower($2))`
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

	var districts []entity.District
	for rows.Next() {
		var district entity.District
		err = rows.Scan(&district.Id, &district.Name, &district.CityId)
		exception.PanicIfErr(err)
		districts = append(districts, district)
	}
	return &districts
}

func (r *WilayahRepository) SubDistricts(ctx context.Context) *[]entity.Subdistrict {
	search := ctx.Value("search").(string)
	sortBy := ctx.Value("sortBy").(string)
	sortDir := ctx.Value("sortDir").(string)
	districtId := ctx.Value("districtId").(string)

	allowedSortBy := map[string]bool{
		"id":   true,
		"name": true,
	}

	allowedSortDir := map[string]bool{
		"asc":  true,
		"desc": true,
	}

	// Default values kalau input nggak valid
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
		query += ` AND (lower(id) LIKE lower($2) OR lower(name) LIKE lower($2))`
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

	var subdistricts []entity.Subdistrict
	for rows.Next() {
		var subdistrict entity.Subdistrict
		err = rows.Scan(&subdistrict.Id, &subdistrict.Name, &subdistrict.DistrictId)
		exception.PanicIfErr(err)
		subdistricts = append(subdistricts, subdistrict)
	}
	return &subdistricts
}
