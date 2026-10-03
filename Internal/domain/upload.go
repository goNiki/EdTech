package domain

type FileUploadResult struct {
	FileURL   string
	FileName  string
	SizeBytes int64
	MimeType  string
}
