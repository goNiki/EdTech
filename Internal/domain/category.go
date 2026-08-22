package domain

type Category struct {
	ID          int64  `db:"id"`
	Name        string `db:"name"`
	Slug        string `db:"slug"`
	Description string `db:"description"`
	ParentID    *int64 `db:"parent_id"`
}
