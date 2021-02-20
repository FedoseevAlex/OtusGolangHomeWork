package main

import (
	"errors"
	"io"
	"os"
)

var (
	ErrNoSrcFileSpecified    = errors.New("no src filename specified")
	ErrNoDstFileSpecified    = errors.New("no dst filename specified")
	ErrLimitIsNegative       = errors.New("negative limit is specified")
	ErrOffsetIsNegative      = errors.New("negative offset is specified")
	ErrReadWithoutLimit      = errors.New("no limit specified for file with unknown length")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

const defaultChunkSize int64 = 4096

func min(a, b int64) int64 {
	if a <= b {
		return a
	}
	return b
}

func Copy(fromPath, toPath string, offset, limit int64) error {
	fromFile, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer fromFile.Close()

	// Rewind file to an offset
	_, err = fromFile.Seek(offset, 0)
	if err != nil {
		return err
	}

	totalBytesToCopy := limit
	if totalBytesToCopy == 0 {
		info, err := fromFile.Stat()
		if err != nil {
			return err
		}
		totalBytesToCopy = info.Size() - offset
	}

	totalBytesWritten := int64(0)

	// Now it's okay to open dst file
	toFile, err := os.OpenFile(toPath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0664)
	if err != nil {
		return err
	}
	defer toFile.Close()

	for totalBytesToCopy != 0 {
		copyChunk := min(defaultChunkSize, totalBytesToCopy)

		written, err := io.CopyN(toFile, fromFile, copyChunk)
		totalBytesToCopy -= written
		if err != nil {
			if !errors.Is(err, io.EOF) {
				return err
			}
			totalBytesToCopy = 0
		}
		totalBytesWritten += written
	}
	return nil
}

// Validate checks input parameters and returns an error
// if some check failed.
func Validate(from, to string, offset, limit int64) error {
	if from == "" {
		return ErrNoSrcFileSpecified
	}

	if to == "" {
		return ErrNoDstFileSpecified
	}

	if offset < 0 {
		return ErrOffsetIsNegative
	}

	if limit < 0 {
		return ErrLimitIsNegative
	}

	fromInfo, err := os.Stat(from)
	if err != nil {
		return err
	}

	if offset > fromInfo.Size() {
		return ErrOffsetExceedsFileSize
	}

	if fromInfo.Size() == 0 && limit == 0 {
		return ErrReadWithoutLimit
	}

	return nil
}
