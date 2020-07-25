package thumb

import (
	"io/ioutil"
	"testing"

	"github.com/h2non/bimg"
)

func TestDetermine(t *testing.T) {
	buf, _ := ioutil.ReadFile(testGIF)
	if "gif" != bimg.DetermineImageTypeName(buf) {
		t.Fatalf("Determine Error: file=%#+v", testGIF)
	}

	buf, _ = ioutil.ReadFile(testJPG)
	if "jpeg" != bimg.DetermineImageTypeName(buf) {
		t.Fatalf("Determine Error: file=%#+v", testJPG)
	}

	buf, _ = ioutil.ReadFile(testPNG)
	if "png" != bimg.DetermineImageTypeName(buf) {
		t.Fatalf("Determine Error: file=%#+v", testPNG)
	}
}

func TestResize(t *testing.T) {
	buffer, err := bimg.Read(testPNG)
	if err != nil {
		t.Fatalf("Resize Error: file=%s: %+v", testPNG, err)
	}

	newImage, err := bimg.NewImage(buffer).Resize(800, 600)
	if err != nil {
		t.Fatalf("Resize Error: file=%#+v: %+v", testPNG, err)
	}

	size, err := bimg.NewImage(newImage).Size()
	if size.Width != 800 && size.Height != 600 {
		t.Fatalf("The image size is valid: %+v", err)
	}

	// bimg.Write("resize.png", newImage)
}

func TestForceResize(t *testing.T) {
	buffer, err := bimg.Read(testPNG)
	if err != nil {
		t.Fatalf("Resize Error: file=%s: %+v", testPNG, err)
	}

	newImage, err := bimg.NewImage(buffer).ForceResize(800, 600)
	if err != nil {
		t.Fatalf("Resize Error: file=%#+v: %+v", testPNG, err)
	}

	size, err := bimg.NewImage(newImage).Size()
	if size.Width != 800 && size.Height != 600 {
		t.Fatalf("The image size is valid: %+v", err)
	}

	// bimg.Write("force-resize.png", newImage)
}

func TestThumbnail(t *testing.T) {
	buffer, err := bimg.Read(testJPG)
	if err != nil {
		t.Fatalf("Resize Error: file=%s: %+v", testJPG, err)
	}

	newImage, err := bimg.NewImage(buffer).Thumbnail(100)
	if err != nil {
		t.Fatalf("Resize Error: file=%#+v: %+v", testJPG, err)
	}

	size, err := bimg.NewImage(newImage).Size()
	if size.Width != 100 && size.Height != 100 {
		t.Fatalf("The image size is valid: %+v", err)
	}

	// bimg.Write("thumbnail.jpg", newImage)
}
