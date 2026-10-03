package dto

type FileUploadResponse struct {
	FileURL   string `json:"file_url"`
	FileName  string `json:"file_name"`
	SizeBytes int64  `json:"size_bytes"`
	MimeType  string `json:"mime_type"`
}
