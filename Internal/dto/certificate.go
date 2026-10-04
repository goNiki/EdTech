package dto

import "time"

type CertificateResponse struct {
	ID              int64     `json:"id"`
	CertificateCode string    `json:"certificate_code"`
	UserID          int64     `json:"user_id"`
	CourseID        int64     `json:"course_id"`
	StudentName     string    `json:"student_name"`
	CourseTitle     string    `json:"course_title"`
	FinalScore      float64   `json:"final_score"`
	IssuedAt        time.Time `json:"issued_at"`
}

type VerifyCertificateResponse struct {
	Valid       bool                 `json:"valid"`
	Certificate *CertificateResponse `json:"certificate,omitempty"`
}
