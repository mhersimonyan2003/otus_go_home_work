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
)

func Copy(fromPath, toPath string, offset, limit int64) error {
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

	if fileInfo.Mode().Perm()&0400 == 0 {
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
