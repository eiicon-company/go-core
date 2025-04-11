package optimum

import (
	"bytes"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"

	"github.com/gabriel-vasile/mimetype"
	"golang.org/x/xerrors"

	"github.com/eiicon-company/go-core/util/filedetect"
)

// Optimize reduce image size
func Optimize(buf []byte) ([]byte, error) {
	if !filedetect.IsImage(buf) {
		return nil, xerrors.New("file is not an image")
	}

	mime := mimetype.Detect(buf)
	if mime == nil {
		return nil, xerrors.Errorf("file is not supported")
	}

	switch mime.Extension() {
	default:
		return nil, xerrors.Errorf("ext %s is not supported", mime.Extension())
	case ".jpeg", ".jpg":
		return OptimizeJPG(buf)
	case ".gif":
		return OptimizeGIF(buf)
	case ".png":
		return OptimizePNG(buf)
	}
}

// OptimizeGIFReader re-encodes GIF without external tools
func OptimizeGIFReader(reader io.Reader) ([]byte, error) {
	img, _, err := image.Decode(reader)
	if err != nil {
		return nil, xerrors.Errorf("failed to decode GIF image: %w", err)
	}

	var buf bytes.Buffer
	err = gif.Encode(&buf, img, nil)
	if err != nil {
		return nil, xerrors.Errorf("failed to encode GIF image: %w", err)
	}

	return buf.Bytes(), nil
}

// OptimizeGIF re-encodes GIF without external tools
func OptimizeGIF(buf []byte) ([]byte, error) {
	return OptimizeGIFReader(bytes.NewReader(buf))
}

// OptimizeJPGReader reduce JPG size
func OptimizeJPGReader(reader io.Reader) ([]byte, error) {
	img, _, err := image.Decode(reader)
	if err != nil {
		return nil, xerrors.Errorf("failed to decode JPEG image: %w", err)
	}

	var buf bytes.Buffer
	err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: 80})
	if err != nil {
		return nil, xerrors.Errorf("failed to encode JPEG image: %w", err)
	}

	return buf.Bytes(), nil
}

// OptimizeJPG reduce JPG size
func OptimizeJPG(buf []byte) ([]byte, error) {
	return OptimizeJPGReader(bytes.NewReader(buf))
}

// OptimizePNGReader reduce PNG size
func OptimizePNGReader(reader io.Reader) ([]byte, error) {
	img, _, err := image.Decode(reader)
	if err != nil {
		return nil, xerrors.Errorf("failed to decode PNG image: %w", err)
	}

	var buf bytes.Buffer
	encoder := png.Encoder{CompressionLevel: png.BestCompression}
	err = encoder.Encode(&buf, img)
	if err != nil {
		return nil, xerrors.Errorf("failed to encode PNG image: %w", err)
	}

	return buf.Bytes(), nil
}

// OptimizePNG reduce PNG size
func OptimizePNG(buf []byte) ([]byte, error) {
	return OptimizePNGReader(bytes.NewReader(buf))
}
