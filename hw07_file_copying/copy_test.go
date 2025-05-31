package main

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestCopyHappyPaths(t *testing.T) {
	type tc struct {
		offset, limit int64
		etalon        string
	}
	cases := []tc{
		{0, 0, "out_offset0_limit0.txt"},
		{0, 10, "out_offset0_limit10.txt"},
		{0, 1000, "out_offset0_limit1000.txt"},
		{0, 10000, "out_offset0_limit10000.txt"},
		{100, 1000, "out_offset100_limit1000.txt"},
		{6000, 1000, "out_offset6000_limit1000.txt"},
	}

	for _, c := range cases {
		tmpDir := t.TempDir()
		dst := filepath.Join(tmpDir, "out.txt")

		if err := Copy("testdata/input.txt", dst, c.offset, c.limit); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		want, _ := os.ReadFile(filepath.Join("testdata", c.etalon))
		got, _ := os.ReadFile(dst)
		if string(got) != string(want) {
			t.Errorf("mismatch for offset=%d limit=%d", c.offset, c.limit)
		}
	}
}

func TestCopyOffsetTooLarge(t *testing.T) {
	tmpDir := t.TempDir()
	err := Copy("testdata/input.txt", filepath.Join(tmpDir, "out.txt"), 1<<20, 0)
	if !errors.Is(err, ErrOffsetExceedsFileSize) {
		t.Fatalf("expected ErrOffsetExceedsFileSize, got %v", err)
	}
}

func TestCopyUnsupportedFile(t *testing.T) {
	tmpDir := t.TempDir()
	err := Copy("/dev/urandom", filepath.Join(tmpDir, "out.txt"), 0, 0)
	if !errors.Is(err, ErrUnsupportedFile) {
		t.Fatalf("expected ErrUnsupportedFile, got %v", err)
	}
}
