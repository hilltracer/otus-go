package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"time"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	if offset < 0 || limit < 0 {
		return fmt.Errorf("offset and limit must be non-negative")
	}

	src, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer src.Close()

	info, err := src.Stat()
	if err != nil {
		return err
	}
	if !info.Mode().IsRegular() {
		return ErrUnsupportedFile
	}

	size := info.Size()
	if offset > size {
		return ErrOffsetExceedsFileSize
	}

	// Destination
	dst, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Offset
	if _, err = src.Seek(offset, io.SeekStart); err != nil {
		return err
	}

	// Size for copying
	toCopy := limit
	if limit == 0 || limit > size-offset {
		toCopy = size - offset
	}

	const bufSize = 16
	buf := make([]byte, bufSize)
	var copied int64

	for copied < toCopy {
		remain := toCopy - copied
		if remain < bufSize {
			buf = buf[:remain]
		}

		n, rErr := src.Read(buf)
		if n > 0 {
			wN, wErr := dst.Write(buf[:n])
			if wErr != nil {
				return wErr
			}
			if wN != n {
				return io.ErrShortWrite
			}
			copied += int64(n)
			printProgress(copied, toCopy)
		}
		if rErr != nil {
			if rErr == io.EOF {
				break
			}
			return rErr
		}
		time.Sleep(1 * time.Millisecond) // test progress =)
	}
	fmt.Print("\r100%\n")
	return nil
}

func printProgress(done, total int64) {
	pct := int(float64(done) / float64(total) * 100)
	fmt.Printf("\r%3d%%", pct)
}
