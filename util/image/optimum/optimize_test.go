package optimum

import (
	"os"
	"testing"
)

func TestOptimizeALL(t *testing.T) {
	buf, err := os.ReadFile(testGIF)
	if err != nil {
		t.Fatalf("OptimizeGIF Error: file=%#+v: %+v", testGIF, err)
	}

	out, err := Optimize(buf)
	if err != nil {
		t.Fatalf("OptimizeGIF Error: file=%#+v: %+v", testGIF, err)
	}

	if len(out) <= 500 {
		t.Fatalf("OptimizeGIF Error: file=%#+v something went wrong", testGIF)
	}

	_ = os.WriteFile("test-optimize-compressed.gif", out, 0600)

	buf, err = os.ReadFile(testJPG)
	if err != nil {
		t.Fatalf("OptimizeJPG Error: file=%#+v: %+v", testJPG, err)
	}

	out, err = Optimize(buf)
	if err != nil {
		t.Fatalf("OptimizeJPG Error: file=%#+v: %+v", testJPG, err)
	}

	if len(out) <= 500 {
		t.Fatalf("OptimizeJPG Error: file=%#+v something went wrong", testGIF)
	}

	_ = os.WriteFile("test-optimize-compressed.jpg", out, 0600)

	buf, err = os.ReadFile(testPNG)
	if err != nil {
		t.Fatalf("OptimizePNG Error: file=%#+v: %+v", testPNG, err)
	}

	out, err = Optimize(buf)
	if err != nil {
		t.Fatalf("OptimizePNG Error: file=%#+v: %+v", testPNG, err)
	}

	if len(out) <= 500 {
		t.Fatalf("OptimizePNG Error: file=%#+v something went wrong", testGIF)
	}

	_ = os.WriteFile("test-optimize-compressed.png", out, 0600)
}

func TestOptimizeGIF(t *testing.T) {
	f, err := os.Open(testGIF)
	if err != nil {
		t.Fatalf("OptimizeGIF Error: file=%#+v: %+v", testGIF, err)
	}

	out, err := OptimizeGIFReader(f)
	if err != nil {
		t.Fatalf("OptimizeGIF Error: file=%#+v: %+v", testGIF, err)
	}

	if len(out) <= 500 {
		t.Fatalf("OptimizeGIF Error: file=%#+v something went wrong", testGIF)
	}

	_ = os.WriteFile("compressed.gif", out, 0600)
}

func TestOptimizeJPG(t *testing.T) {
	f, err := os.Open(testJPG)
	if err != nil {
		t.Fatalf("OptimizeJPG Error: file=%#+v: %+v", testJPG, err)
	}

	out, err := OptimizeJPGReader(f)
	if err != nil {
		t.Fatalf("OptimizeJPG Error: file=%#+v: %+v", testJPG, err)
	}

	if len(out) <= 500 {
		t.Fatalf("OptimizeJPG Error: file=%#+v something went wrong", testGIF)
	}

	_ = os.WriteFile("compressed.jpg", out, 0600)
}

func TestOptimizePNG(t *testing.T) {
	f, err := os.Open(testPNG)
	if err != nil {
		t.Fatalf("OptimizePNG Error: file=%#+v: %+v", testPNG, err)
	}

	out, err := OptimizePNGReader(f)
	if err != nil {
		t.Fatalf("OptimizePNG Error: file=%#+v: %+v", testPNG, err)
	}

	if len(out) <= 500 {
		t.Fatalf("OptimizePNG Error: file=%#+v something went wrong", testGIF)
	}

	_ = os.WriteFile("compressed.png", out, 0600)
}
