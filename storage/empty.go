package storage

import (
	"context"
	"errors"
	"io"
)

var (
	insEmpty = &empty{}
	ErrEmpty = errors.New("empty storage")
)

type empty struct{}

func (e empty) Put(ctx context.Context, name string, r io.Reader) (key string, err error) {
	err = ErrEmpty
	return
}

func (e empty) Get(ctx context.Context, name string) (rc io.ReadCloser, err error) {
	err = ErrEmpty
	return
}

func (e empty) Delete(ctx context.Context, name string) (err error) {
	err = ErrEmpty
	return
}

func (e empty) SignGetUrl(ctx context.Context, name string, expires int64, attachment ...bool) (s string, err error) {
	err = ErrEmpty
	return
}

func (e empty) SignPutUrl(ctx context.Context, name string, expires int64) (s string, err error) {
	err = ErrEmpty
	return
}
