package optimum

import (
	"io/ioutil"
	"os"
	"os/exec"
	"testing"
)

func TestOptimizeALL(t *testing.T) {
	if _, err := exec.LookPath(gifOptimizer); err == nil {
		buf, err := ioutil.ReadFile(testGIF)
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

		_ = ioutil.WriteFile("test-optimize-compressed.gif", out, 0600)
	} else {
		t.Logf("OptimizeGIF Skip: file=%#+v: %+v", testGIF, err)
	}

	if _, err := exec.LookPath(jpgOptimizer); err == nil {
		buf, err := ioutil.ReadFile(testJPG)
		if err != nil {
			t.Fatalf("OptimizeJPG Error: file=%#+v: %+v", testJPG, err)
		}

		out, err := Optimize(buf)
		if err != nil {
			t.Fatalf("OptimizeJPG Error: file=%#+v: %+v", testJPG, err)
		}

		if len(out) <= 500 {
			t.Fatalf("OptimizeJPG Error: file=%#+v something went wrong", testGIF)
		}

		_ = ioutil.WriteFile("test-optimize-compressed.jpg", out, 0600)
	} else {
		t.Logf("OptimizeJPG Skip: file=%#+v: %+v", testJPG, err)
	}

	if _, err := exec.LookPath(pngOptimizer); err == nil {
		buf, err := ioutil.ReadFile(testPNG)
		if err != nil {
			t.Fatalf("OptimizePNG Error: file=%#+v: %+v", testPNG, err)
		}

		out, err := Optimize(buf)
		if err != nil {
			t.Fatalf("OptimizePNG Error: file=%#+v: %+v", testPNG, err)
		}

		if len(out) <= 500 {
			t.Fatalf("OptimizePNG Error: file=%#+v something went wrong", testGIF)
		}

		_ = ioutil.WriteFile("test-optimize-compressed.png", out, 0600)
	} else {
		t.Logf("OptimizePNG Skip: file=%#+v: %+v", testPNG, err)
	}
}

func TestOptimizeGIF(t *testing.T) {
	if _, err := exec.LookPath(gifOptimizer); err != nil {
		t.Skipf("OptimizeGIF Skip: file=%#+v: %+v", testGIF, err)
	}

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

	_ = ioutil.WriteFile("compressed.gif", out, 0600)
}

func TestOptimizeJPG(t *testing.T) {
	if _, err := exec.LookPath(jpgOptimizer); err != nil {
		t.Skipf("OptimizeJPG Skip: file=%#+v: %+v", testJPG, err)
	}

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

	_ = ioutil.WriteFile("compressed.jpg", out, 0600)
}

func TestOptimizePNG(t *testing.T) {
	if _, err := exec.LookPath(pngOptimizer); err != nil {
		t.Skipf("OptimizePNG Skip: file=%#+v: %+v", testPNG, err)
	}

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

	_ = ioutil.WriteFile("compressed.png", out, 0600)
}
