package optimum

import (
	"bytes"
	"io"
	"os/exec"

	"github.com/codeskyblue/go-sh"
	"github.com/h2non/filetype"
	"golang.org/x/xerrors"
)

// Optimize reduce image size
func Optimize(buf []byte) ([]byte, error) {
	if !filetype.IsImage(buf) {
		return nil, xerrors.New("file is not an image")
	}

	kind, err := filetype.Match(buf)
	if err != nil {
		return nil, xerrors.Errorf("ext %v is not supported: %+v", kind, err)
	}

	switch kind.Extension {
	default:
		return nil, xerrors.Errorf("ext %s is not supported", kind.Extension)
	case "jpeg", "jpg":
		return OptimizeJPG(buf)
	case "gif":
		return OptimizeGIF(buf)
	case "png":
		return OptimizePNG(buf)
	}
}

const (
	jpgOptimizer = "cjpeg"
	pngOptimizer = "pngquant"
	gifOptimizer = "gifsicle"
)

// OptimizeGIFReader reduce GIF size
func OptimizeGIFReader(reader io.Reader) ([]byte, error) {
	path, err := exec.LookPath(gifOptimizer)
	if err != nil {
		return nil, err
	}

	return sh.Command(path, "--optimize=3").
		SetStdin(reader).
		Command("cat", "-").
		Output()
}

// OptimizeGIF reduce GIF size
func OptimizeGIF(buf []byte) ([]byte, error) {
	return OptimizeGIFReader(bytes.NewReader(buf))
}

// OptimizeJPGReader reduce JPG size
func OptimizeJPGReader(reader io.Reader) ([]byte, error) {
	path, err := exec.LookPath(jpgOptimizer)
	if err != nil {
		return nil, err
	}

	return sh.Command(path, "-quality", "50,80").
		SetStdin(reader).
		Output()
}

// OptimizeJPG reduce JPG size
func OptimizeJPG(buf []byte) ([]byte, error) {
	return OptimizeJPGReader(bytes.NewReader(buf))
}

// OptimizePNGReader reduce PNG size
func OptimizePNGReader(reader io.Reader) ([]byte, error) {
	path, err := exec.LookPath(pngOptimizer)
	if err != nil {
		return nil, err
	}

	return sh.Command(path, "--quality", "50-80", "--speed", "3", "-").
		SetStdin(reader).
		Output()
}

// OptimizePNG reduce PNG size
func OptimizePNG(buf []byte) ([]byte, error) {
	return OptimizePNGReader(bytes.NewReader(buf))
}
