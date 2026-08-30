package domain

const (
	ResourceCourse   = "course"
	ResourceLesson   = "lesson"
	ResourceResource = "resource"
	ResourceUser     = "user"
)

const (
	ActionView        = "view"
	ActionEdit        = "edit"
	ActionDelete      = "delete"
	ActionManageusers = "manage_users"
	ActionEnroll      = "enroll"
	ActionPublish     = "publish"
)

type CoursePermissions struct {
	CanView        bool `json:"can_view"`
	CanEdit        bool `json:"can_edit"`
	CanDelete      bool `json:"can_delete"`
	CanPublish     bool `json:"can_publish"`
	CanEnroll      bool `json:"can_enroll"`
	CanManageUsers bool `json:"can_manage_users"`
}
