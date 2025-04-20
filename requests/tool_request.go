package requests

type CreateToolRequest struct {
	Name        string  `json:"name" validate:"required,min=2,max=255"`
	Description *string `json:"description" validate:"omitempty,min=2"`
	Link        string  `json:"link" validate:"required,min=2,max=255"`
	IsPublished bool    `json:"is_published" validate:"required,boolean"`
	Type        string  `json:"type" validate:"required"` // dev tools,php packages,go packages,tutorial
}

type UpdateToolRequest struct {
	Name        string  `json:"name" validate:"required,min=2,max=255"`
	Description *string `json:"description" validate:"omitempty,min=2"`
	Link        string  `json:"link" validate:"required,min=2,max=255"`
	IsPublished bool    `json:"is_published" validate:"required,boolean"`
	Type        string  `json:"type" validate:"required"` // dev tools,php packages,go packages,tutorial
}
