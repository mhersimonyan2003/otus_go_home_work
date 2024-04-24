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
	ErrInvalidOffset         = errors.New("invalid offset")
	ErrInvalidLimit          = errors.New("invalid limit")
	ErrSameFile              = errors.New("files are the same")
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

	fileInfo, err := fromFile.Stat()
	fileSize := fileInfo.Size()
	if err != nil {
		return err
	}

	if fileInfo.Mode().Perm()&0o400 == 0 {
		return ErrUnsupportedFile
	}

	if offset > fileSize {
		return ErrOffsetExceedsFileSize
	}

	toFile, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer toFile.Close()

	toFileInfo, err := toFile.Stat()
	if err != nil {
		return err
	}

	if os.SameFile(fileInfo, toFileInfo) {
		return ErrUnsupportedFile
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
