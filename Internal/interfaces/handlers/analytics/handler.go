package analytics

import (
	"edtech/internal/interfaces/middleware/auth"
	"edtech/internal/service"
)

type AnalyticsHandler struct {
	analyticsService service.AnalyticsServices
	authMiddleware    auth.AuthMiddleware
}

func NewAnalyticsHandler(
	analyticsService service.AnalyticsServices,
	authMiddleware auth.AuthMiddleware,
) *AnalyticsHandler {
	return &AnalyticsHandler{
		analyticsService: analyticsService,
		authMiddleware:    authMiddleware,
	}
}
