package request

type PresignUploadRequest struct {
	Filename string `json:"filename" validate:"required"`
	Size     int64  `json:"size" validate:"required,max=1048576"`
}
