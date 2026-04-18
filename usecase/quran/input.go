package quran

type FindAllSurahInput struct {
	Search  string
	SortBy  string
	SortDir string
}

type FindAllAyahInput struct {
	SurahId string
	Search  string
	SortBy  string
	SortDir string
}