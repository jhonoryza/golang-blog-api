package request

type QuranRequest struct {
	SurahId string `validate:"required"`
}