package requests

type CreatePostRequest struct {
	Title         string  `json:"title" validate:"required,min=3"`
	Content       string  `json:"content" validate:"required,min=3"`
	Summary       string  `json:"summary" validate:"required"`
	PublishedAt   *string `json:"published_at" validate:"omitempty"`
	AuthorId      int     `json:"author_id" validate:"required"`
	IsMarkdown    bool    `json:"is_markdown" validate:"required,boolean"`
	ImageUrl      string  `json:"image_url" validate:"required"`
	ImageTwUrl    string  `json:"image_tw_url" validate:"required"`
	ImageThumbUrl string  `json:"image_thumb_url" validate:"required"`
}

type UpdatePostRequest struct {
	Title         string  `json:"title" validate:"required,min=3"`
	Slug          *string `json:"slug" validate:"omitempty,min=3,alphanum"`
	Content       string  `json:"content" validate:"required,min=3"`
	Summary       string  `json:"summary" validate:"required"`
	PublishedAt   *string `json:"published_at" validate:"omitempty"`
	AuthorId      int     `json:"author_id" validate:"required"`
	IsMarkdown    bool    `json:"is_markdown" validate:"required,boolean"`
	ImageUrl      string  `json:"image_url" validate:"required"`
	ImageTwUrl    string  `json:"image_tw_url" validate:"required"`
	ImageThumbUrl string  `json:"image_thumb_url" validate:"required"`
}
