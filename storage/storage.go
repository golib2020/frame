package storage

import (
	"io"
)

type Storage interface {
	Exists(path string) (bool, error)
	Get(path string) (io.Reader, error)
	Put(path string, reader io.Reader) error
	Delete(path string) error
	Copy(path1, path2 string) error
	Move(path1, path2 string) error
	Path(path string) string
	Url(path string) string
}
