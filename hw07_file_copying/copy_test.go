package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

type CopyParams struct {
	Limit  int64
	Offset int64
}

func TestCopy(t *testing.T) {
	// Place your code here.
	paramsSlice := []CopyParams{
		{Offset: 0, Limit: 0},
		{Offset: 0, Limit: 10},
		{Offset: 0, Limit: 1000},
		{Offset: 0, Limit: 10000},
		{Offset: 100, Limit: 1000},
		{Offset: 6000, Limit: 1000},
		{Offset: 1000, Limit: 6000},
	}

	for _, params := range paramsSlice {
		t.Run(fmt.Sprintf("offset %d, limit %d", params.Offset, params.Limit), func(t *testing.T) {
			require.NoError(t, Copy("./testdata/input.txt", "./testdata/output.txt", params.Offset, params.Limit))

			file1, err := os.Open("./testdata/output.txt")
			t.Cleanup(func() {
				os.Remove("./testdata/output.txt")
			})

			if err != nil {
				return
			}
			defer file1.Close()

			file2, err := os.Open(fmt.Sprintf("./testdata/out_offset%d_limit%d.txt", params.Offset, params.Limit))
			if err != nil {
				return
			}
			defer file2.Close()

			scanner1 := bufio.NewScanner(file1)
			scanner2 := bufio.NewScanner(file2)

			for scanner1.Scan() && scanner2.Scan() {
				line1 := scanner1.Text()
				line2 := scanner2.Text()

				require.Equal(t, line1, line2)
			}

			require.Equal(t, scanner1.Scan(), scanner2.Scan())

			require.NoError(t, scanner1.Err())
			require.NoError(t, scanner2.Err())
		})
	}
	t.Run("invalid offset", func(t *testing.T) {
		err := Copy("./testdata/input.txt", "./testdata/output.txt", -1, 10)
		require.Error(t, err)
		require.True(t, errors.Is(err, ErrInvalidOffset))
	})
	t.Run("invalid limit", func(t *testing.T) {
		err := Copy("./testdata/input.txt", "./testdata/output.txt", 10, -1)
		require.Error(t, err)
		require.True(t, errors.Is(err, ErrInvalidLimit))
	})
	t.Run("same content, differen files", func(t *testing.T) {
		err := Copy("./testdata/same_file.txt", "./testdata/same_file_1.txt", 0, 10)
		require.NoError(t, err)
	})
	t.Run("same files", func(t *testing.T) {
		err := Copy("./testdata/same_file.txt", "./testdata/same_file.txt", 0, 10)
		require.Error(t, err)
		require.True(t, errors.Is(err, ErrSameFile))
	})
}
