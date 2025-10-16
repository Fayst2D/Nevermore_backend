package dto

import (
	"bytes"
	"mime/multipart"
)

type FileInfo struct {
	File   multipart.File
	Header *multipart.FileHeader
}

type BytesFile struct {
	*bytes.Reader
}

func (b *BytesFile) Close() error {
	return nil
}
