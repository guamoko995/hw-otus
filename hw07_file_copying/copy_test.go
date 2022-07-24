package main

import (
	"errors"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	outPath := "copy.txt"
	defer os.Remove(outPath)

	t.Run("offset=0 limit=0", func(t *testing.T) {
		err := Copy("testdata/input.txt", outPath, 0, 0)
		require.Nil(t, err)
		equiv(t, "testdata/out_offset0_limit0.txt", outPath)
	})

	t.Run("offset=0 limit=10", func(t *testing.T) {
		err := Copy("testdata/input.txt", outPath, 0, 10)
		require.Nil(t, err)
		equiv(t, "testdata/out_offset0_limit10.txt", outPath)
	})

	t.Run("offset=0 limit=1000", func(t *testing.T) {
		err := Copy("testdata/input.txt", outPath, 0, 1000)
		require.Nil(t, err)
		equiv(t, "testdata/out_offset0_limit1000.txt", outPath)
	})

	t.Run("offset=0 limit=10000", func(t *testing.T) {
		err := Copy("testdata/input.txt", outPath, 0, 10000)
		require.Nil(t, err)
		equiv(t, "testdata/out_offset0_limit10000.txt", outPath)
	})

	t.Run("offset=100 limit=1000", func(t *testing.T) {
		err := Copy("testdata/input.txt", outPath, 100, 1000)
		require.Nil(t, err)
		equiv(t, "testdata/out_offset100_limit1000.txt", outPath)
	})

	t.Run("offset=6000 limit=1000", func(t *testing.T) {
		err := Copy("testdata/input.txt", outPath, 6000, 1000)
		require.Nil(t, err)
		equiv(t, "testdata/out_offset6000_limit1000.txt", outPath)
	})

	t.Run("offset more syze", func(t *testing.T) {
		err := Copy("testdata/input.txt", outPath, 10000, 0)
		require.Equal(t, ErrOffsetExceedsFileSize, err)
	})

	t.Run("limit more syze", func(t *testing.T) {
		err := Copy("testdata/input.txt", outPath, 0, 10000)
		require.Nil(t, err)
		equiv(t, "testdata/out_offset0_limit0.txt", outPath)
	})

	t.Run("unsupported file", func(t *testing.T) {
		err := Copy("/dev/urandom", outPath, 0, 0)
		require.Equal(t, ErrUnsupportedFile, err)
	})

	t.Run("not a valid limit", func(t *testing.T) {
		err := Copy("testdata/input.txt", outPath, 0, -1)
		require.Equal(t, ErrNotValidLimit, err)
	})
}

func equiv(t *testing.T, path1, path2 string) {
	t.Helper()
	equiv, err := equivalentFiles(path1, path2)
	require.Nil(t, err)
	require.True(t, equiv)
}

func equivalentFiles(path1, path2 string) (bool, error) {
	f1, err := os.Open(path1)
	if err != nil {
		return false, err
	}
	defer f1.Close()
	r1 := io.Reader(f1)

	f2, err := os.Open(path2)
	if err != nil {
		return false, err
	}

	defer f2.Close()
	r2 := io.Reader(f2)

	b1 := []byte{0}
	b2 := []byte{0}

	for {
		n1, err1 := r1.Read(b1)
		n2, err2 := r2.Read(b2)
		if !errors.Is(err1, err2) || n1 != n2 || b1[0] != b2[0] {
			return false, nil
		}
		if err1 == io.EOF {
			return true, nil
		}
	}
}
