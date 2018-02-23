package store

import "io"

type Storage interface {
	Put(src io.Reader) (string, error)
	Get(key string) (io.Reader, error)
}
