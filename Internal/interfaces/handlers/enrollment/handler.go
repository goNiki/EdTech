package enrollment

import (
	"edtech/internal/interfaces/middleware/auth"
	"edtech/internal/service"

	"github.com/go-playground/validator/v10"
)

type EnrollmentHandler struct {
	enrollmentService service.EnrolledServices
	authMiddleware    auth.AuthMiddleware
	validator         *validator.Validate
}

func NewEnrollmentHandler(
	enrollmentService service.EnrolledServices,
	authMiddleware auth.AuthMiddleware,
) *EnrollmentHandler {
	return &EnrollmentHandler{
		enrollmentService: enrollmentService,
		authMiddleware:    authMiddleware,
		validator:         validator.New(),
	}
}
