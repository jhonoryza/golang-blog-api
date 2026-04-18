package entities

type Province struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type City struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	ProvinceId string `json:"provinceId"`
}

type District struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	CityId string `json:"cityId"`
}

type Subdistrict struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	DistrictId string `json:"districtId"`
}