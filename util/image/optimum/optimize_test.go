package optimum

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOptimizeALL(t *testing.T) {
	if _, err := exec.LookPath(gifOptimizer); err == nil {
		buf, err := os.ReadFile(testGIF)
		require.NoError(t, err)

		out, err := Optimize(buf)
		require.NoError(t, err)
		require.Greater(t, len(out), 500)

		_ = os.WriteFile("test-optimize-compressed.gif", out, 0600)
	} else {
		t.Logf("OptimizeGIF Skip: file=%#+v: %+v", testGIF, err)
	}

	if _, err := exec.LookPath(jpgOptimizer); err == nil {
		buf, err := os.ReadFile(testJPG)
		require.NoError(t, err)

		out, err := Optimize(buf)
		require.NoError(t, err)
		require.Greater(t, len(out), 500)

		_ = os.WriteFile("test-optimize-compressed.jpg", out, 0600)
	} else {
		t.Logf("OptimizeJPG Skip: file=%#+v: %+v", testJPG, err)
	}

	if _, err := exec.LookPath(pngOptimizer); err == nil {
		buf, err := os.ReadFile(testPNG)
		require.NoError(t, err)

		out, err := Optimize(buf)
		require.NoError(t, err)
		require.Greater(t, len(out), 500)

		_ = os.WriteFile("test-optimize-compressed.png", out, 0600)
	} else {
		t.Logf("OptimizePNG Skip: file=%#+v: %+v", testPNG, err)
	}
}

func TestOptimizeGIF(t *testing.T) {
	if _, err := exec.LookPath(gifOptimizer); err != nil {
		t.Skipf("OptimizeGIF Skip: file=%#+v: %+v", testGIF, err)
	}

	f, err := os.Open(testGIF)
	require.NoError(t, err)

	out, err := OptimizeGIFReader(f)
	require.NoError(t, err)
	require.Greater(t, len(out), 500)

	_ = os.WriteFile("compressed.gif", out, 0600)
}

func TestOptimizeJPG(t *testing.T) {
	if _, err := exec.LookPath(jpgOptimizer); err != nil {
		t.Skipf("OptimizeJPG Skip: file=%#+v: %+v", testJPG, err)
	}

	f, err := os.Open(testJPG)
	require.NoError(t, err)

	out, err := OptimizeJPGReader(f)
	require.NoError(t, err)
	require.Greater(t, len(out), 500)

	_ = os.WriteFile("compressed.jpg", out, 0600)
}

func TestOptimizePNG(t *testing.T) {
	if _, err := exec.LookPath(pngOptimizer); err != nil {
		t.Skipf("OptimizePNG Skip: file=%#+v: %+v", testPNG, err)
	}

	f, err := os.Open(testPNG)
	require.NoError(t, err)

	out, err := OptimizePNGReader(f)
	require.NoError(t, err)
	require.Greater(t, len(out), 500)

	_ = os.WriteFile("compressed.png", out, 0600)
}
