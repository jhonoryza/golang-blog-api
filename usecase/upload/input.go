package upload

type PresignInput struct {
	Filename string
	Size     int64
}

type PresignOutput struct {
	UploadUrl string            `json:"upload_url"`
	Method    string            `json:"method"`
	Headers   map[string]string `json:"headers"`
	Key       string            `json:"key"`
	ImagePath string            `json:"image_path"`
	PublicUrl string            `json:"public_url"`
	ExpiresIn int               `json:"expires_in"`
}
