package filedetect

import (
	"github.com/gabriel-vasile/mimetype"
)

const (
	MimeJpeg     = "image/jpeg"
	MimeJpeg2000 = "image/jp2"
	MimePng      = "image/png"
	MimeGif      = "image/gif"
	MimeWebp     = "image/webp"
	MimeCR2      = "image/x-canon-cr2"
	MimeTiff     = "image/tiff"
	MimeBmp      = "image/bmp"
	MimeJxr      = "image/vnd.ms-photo"
	MimePsd      = "image/vnd.adobe.photoshop"
	MimeIco      = "image/vnd.microsoft.icon"
	MimeHeif     = "image/heif"
	MimeDwg      = "image/vnd.dwg"
)

// IsImage checks if the given buffer is an image type
func IsImage(buf []byte) bool {
	mime := mimetype.Detect(buf)
	if mime == nil {
		return false
	}

	mimes := []string{
		MimeJpeg,
		// MimeJpeg2000,
		MimePng,
		MimeGif,
		MimeWebp,
		// MimeCR2,
		// MimeTiff,
		MimeBmp,
		// MimeJxr,
		// MimePsd,
		// MimeIco,
		// MimeHeif,
		// MimeDwg,
	}
	for _, m := range mimes {
		if mime.Is(m) {
			return true
		}
	}

	return false
}
