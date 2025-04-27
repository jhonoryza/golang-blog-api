package requests

type QuranRequest struct {
	SurahId string `validate:"required"`
}
