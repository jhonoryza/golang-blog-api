package entities

import "database/sql"

type Quran struct {
	ID              int
	ExternalId      int
	Arabic          string
	Latin           string
	Transliteration string
	Translation     string
	NumAyah         int
	Page            int
	Location        string
	CreatedAt       sql.NullTime
	UpdatedAt       sql.NullTime
}

type QuranVerse struct {
	ID          int
	QuranId     int
	Ayah        int
	Page        int
	Juz         int
	Arabic      string
	Kitabah     string
	Latin       string
	Translation string
	AudioUrl    string
	CreatedAt   sql.NullTime
	UpdatedAt   sql.NullTime
}