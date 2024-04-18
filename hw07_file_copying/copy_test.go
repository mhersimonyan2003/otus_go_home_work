package main

import (
	"bufio"
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
			defer os.Remove("./testdata/output.txt")

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
}
