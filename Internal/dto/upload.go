package dto

type FileUploadResponse struct {
	FileURL   string `json:"file_url"`
	FileName  string `json:"file_name"`
	SizeBytes int64  `json:"size_bytes"`
	MimeType  string `json:"mime_type"`
}

type BatchImageItemResponse struct {
	OriginalName string `json:"original_name"`
	FileURL      string `json:"file_url"`
	SizeBytes    int64  `json:"size_bytes"`
}

type BatchImageUploadResponse struct {
	Uploaded []BatchImageItemResponse `json:"uploaded"`
}
