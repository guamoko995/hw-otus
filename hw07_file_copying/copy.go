package main

import (
	"io"
	"os"

	pb "github.com/cheggaaa/pb/v3"
	"github.com/pkg/errors"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrNotValidLimit         = errors.New("not a valid limit")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	// limit не может быть отрицательным
	if limit < 0 {
		return ErrNotValidLimit
	}

	// получение reader
	from, err := os.Open(fromPath)
	if err != nil {
		// return ErrUnsupportedFile
		return errors.Wrapf(err, "failed open file %q", fromPath)
	}
	defer from.Close()
	reader := io.Reader(from)

	// создание writer
	to, err := os.Create(toPath)
	if err != nil {
		return errors.Wrapf(err, "failed create file %q", toPath)
	}
	defer to.Close()
	writer := io.Writer(to)

	// получение информации о копируемом файле
	fi, err := from.Stat()
	if err != nil {
		return errors.Wrapf(err, "could not get information about the file %q", fromPath)
	}

	// файлы с неизвестной длиной не поддерживаются
	if fi.Size() == 0 {
		return ErrUnsupportedFile
	}

	// получение размера с учетом отступа
	if sizeWithOffset := fi.Size() - offset; limit == 0 || limit > sizeWithOffset {
		limit = sizeWithOffset
	}

	// отступ не может превышать размер
	if limit < 0 {
		return ErrOffsetExceedsFileSize
	}

	// установка отступа
	if _, err = from.Seek(offset, 0); err != nil {
		return errors.Wrapf(err, "failed to set offset %d from file %q", offset, toPath)
	}

	// создание прогресс-бара и привязка его к writer через прокси
	bar := pb.Full.Start64(limit)
	barWriter := bar.NewProxyWriter(writer)

	// копирование через прокси
	io.CopyN(barWriter, reader, limit)
	bar.Finish()

	to.Chmod(fi.Mode())
	return nil
}
