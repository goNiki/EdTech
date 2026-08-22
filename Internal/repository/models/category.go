package models

type Category struct {
	ID          int64
	Name        string
	Slug        string
	Description string
	ParentID    *int64
}
