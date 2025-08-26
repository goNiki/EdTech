package dto

// прочитать детельнее статью https://habr.com/ru/companies/first/articles/927460/

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=5"`
	Role     string `json:"role" validate:"required,oneof=user teacher admin"`
}

type RegisterResponse struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
}
