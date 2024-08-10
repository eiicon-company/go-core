package filedetect

import (
	"github.com/gabriel-vasile/mimetype"
)

const (
	MimeDoc  = "application/msword"
	MimeDocx = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	MimeXls  = "application/vnd.ms-excel"
	MimeXlsx = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	MimePpt  = "application/vnd.ms-powerpoint"
	MimePptx = "application/vnd.openxmlformats-officedocument.presentationml.presentation"
)

// IsImage checks if the given buffer is an image type
// func IsImage(buf []byte) bool {
//
// }

// IsDocument checks if the given buffer is an document type
func IsDocument(buf []byte) bool {
	mime := mimetype.Detect(buf)
	if mime == nil {
		return false
	}

	mimes := []string{
		MimePptx,
		MimePpt,
		MimeDocx,
		MimeDoc,
		MimeXlsx,
		MimeXls,
	}
	for _, m := range mimes {
		if mime.Is(m) {
			return true
		}
	}

	return false
}
