package dto

// прочитать детельнее статью https://habr.com/ru/companies/first/articles/927460/

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=5"`
	Username string `json:"username" validate:"required"`
	Role     string `json:"role" validate:"required,oneof=user teacher admin"`
}

type RegisterResponse struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=5"`
}

type LoginResponce struct {
	ID           int64  `json:"id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresID    int64  `json:"expires_in"`
}

type EnrolleRequest struct {
	UserEmail string `json:"user_email"`
	CourseID  int64  `json:"courseid"`
	Role      string `json:"role"`
}
