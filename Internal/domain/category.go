package domain

type Category struct {
	ID           int64
	Name         string
	Slug         string
	Description  string
	IconURL      string
	ParentID     *int64
	Position     int
	CoursesCount int
}
