package wilayah

type FindAllProvincesInput struct {
	Search  string
	SortBy  string
	SortDir string
}

type FindCitiesByProvinceInput struct {
	ProvinceId string
	Search     string
	SortBy     string
	SortDir    string
}

type FindDistrictsByCityInput struct {
	CityId  string
	Search  string
	SortBy  string
	SortDir string
}

type FindSubdistrictsByDistrictInput struct {
	DistrictId string
	Search     string
	SortBy     string
	SortDir    string
}