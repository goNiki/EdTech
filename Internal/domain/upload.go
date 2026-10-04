package domain

import "io"

type FileUploadResult struct {
	FileURL   string
	FileName  string
	SizeBytes int64
	MimeType  string
}

type BatchFileItem struct {
	Reader   io.Reader
	Filename string
	Size     int64
}

type BatchUploadResultItem struct {
	OriginalName string
	FileURL      string
	SizeBytes    int64
	MimeType     string
}
