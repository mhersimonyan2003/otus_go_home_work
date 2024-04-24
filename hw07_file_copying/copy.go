package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrSameFile              = errors.New("same file")
	ErrInvalidOffset         = errors.New("invalid offset")
	ErrInvalidLimit          = errors.New("invalid limit")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	if offset < 0 {
		return ErrInvalidOffset
	}
	if limit < 0 {
		return ErrInvalidLimit
	}

	fromFile, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer fromFile.Close()
	toFile, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer toFile.Close()
	fromFileInfo, err := fromFile.Stat()
	if err != nil {
		return err
	}

	toFileInfo, err := toFile.Stat()
	if err != nil {
		return err
	}

	if os.SameFile(fromFileInfo, toFileInfo) {
		return ErrSameFile
	}

	fileSize := fromFileInfo.Size()
	if err != nil {
		return err
	}

	if fromFileInfo.Mode().Perm()&0o400 == 0 {
		return ErrUnsupportedFile
	}

	if offset > fileSize {
		return ErrOffsetExceedsFileSize
	}

	_, err = fromFile.Seek(offset, 0)
	if err != nil {
		return err
	}

	if limit == 0 || limit > fileSize-offset {
		limit = fileSize - offset
	}

	bar := pb.Start64(limit)
	bar.Set(pb.Bytes, true)
	bar.SetWidth(80)

	reader := bar.NewProxyReader(fromFile)

	_, err = io.CopyN(toFile, reader, limit)
	if err != nil && errors.Is(err, io.EOF) {
		return err
	}

	bar.Finish()

	return nil
}
