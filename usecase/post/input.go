package post

type FindAllInput struct {
	Search  string
	SortBy  string
	SortDir string
}

type FindBySlugInput struct {
	Slug string
}

type CreateInput struct {
	Title        string
	Summary      *string
	Content      string
	PublishedAt  *string
	AuthorId     int
	IsMarkdown   bool
	ImageUrl     string
	ImageTwUrl   string
	ImageThumbUrl string
}

type UpdateInput struct {
	Title        string
	Summary      *string
	Content      string
	Slug         *string
	PublishedAt  *string
	AuthorId     int
	IsMarkdown   bool
	ImageUrl     string
	ImageTwUrl   string
	ImageThumbUrl string
}

type DeleteInput struct {
	Slug string
}