package certificate

import (
	"log/slog"
	"net/http"
	"strconv"

	"edtech/internal/dto"
	"edtech/internal/interfaces/middleware/auth"
	"edtech/internal/interfaces/response"
	"edtech/internal/service"
	errorsAPP "edtech/pkg/errors"

	"github.com/go-chi/chi/v5"
)

type CertificateHandler struct {
	certService    service.CertificateServices
	log            *slog.Logger
	authMiddleware auth.AuthMiddleware
}

func NewCertificateHandler(
	certService service.CertificateServices,
	log *slog.Logger,
	authMiddleware auth.AuthMiddleware,
) *CertificateHandler {
	return &CertificateHandler{
		certService:    certService,
		log:            log,
		authMiddleware: authMiddleware,
	}
}

// GetOrIssueCertificate handles GET /api/v1/courses/{courseid}/certificate
func (h *CertificateHandler) GetOrIssueCertificate(w http.ResponseWriter, r *http.Request) {
	const op = "http.handlers.certificate.GetOrIssueCertificate"

	userID := h.authMiddleware.GetUserID(r.Context())
	if userID == 0 {
		response.HandleError(w, r, h.log, errorsAPP.ErrUnauthorized, op)
		return
	}

	courseID, err := strconv.ParseInt(chi.URLParam(r, "courseid"), 10, 64)
	if err != nil {
		response.HandleError(w, r, h.log, errorsAPP.ErrInvalidURLParam, op)
		return
	}

	cert, err := h.certService.GetOrIssueCertificate(r.Context(), userID, courseID)
	if err != nil {
		response.HandleError(w, r, h.log, err, op)
		return
	}

	response.OK(w, r, dto.CertificateResponse{
		ID:              cert.ID,
		CertificateCode: cert.CertificateCode,
		UserID:          cert.UserID,
		CourseID:        cert.CourseID,
		StudentName:     cert.StudentName,
		CourseTitle:     cert.CourseTitle,
		FinalScore:      cert.FinalScore,
		IssuedAt:        cert.IssuedAt,
	})
}

// VerifyCertificate handles GET /api/v1/certificates/verify/{code}
func (h *CertificateHandler) VerifyCertificate(w http.ResponseWriter, r *http.Request) {
	const op = "http.handlers.certificate.VerifyCertificate"

	code := chi.URLParam(r, "code")
	if code == "" {
		response.HandleError(w, r, h.log, errorsAPP.ErrInvalidURLParam, op)
		return
	}

	cert, err := h.certService.VerifyCertificate(r.Context(), code)
	if err != nil {
		response.HandleError(w, r, h.log, err, op)
		return
	}

	response.OK(w, r, dto.VerifyCertificateResponse{
		Valid: true,
		Certificate: &dto.CertificateResponse{
			ID:              cert.ID,
			CertificateCode: cert.CertificateCode,
			UserID:          cert.UserID,
			CourseID:        cert.CourseID,
			StudentName:     cert.StudentName,
			CourseTitle:     cert.CourseTitle,
			FinalScore:      cert.FinalScore,
			IssuedAt:        cert.IssuedAt,
		},
	})
}
