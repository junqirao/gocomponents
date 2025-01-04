package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type minioCli struct {
	cfg Config
	cli *minio.Client
}

func newMinio(cfg Config) (*minioCli, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessId, cfg.Secret, ""),
		Secure: cfg.SSL,
	})
	if err != nil {
		return nil, err
	}
	return &minioCli{cfg: cfg, cli: client}, nil
}

func (c minioCli) Put(ctx context.Context, name string, reader io.Reader) (key string, err error) {
	bb := &bytes.Buffer{}
	size, err := bb.ReadFrom(reader)
	if err != nil {
		return
	}
	info, err := c.cli.PutObject(ctx, c.cfg.Bucket, name, reader, size, minio.PutObjectOptions{})
	if err != nil {
		return
	}
	key = info.Key
	return
}

func (c minioCli) Get(ctx context.Context, name string) (reader io.ReadCloser, err error) {
	return c.cli.GetObject(ctx, c.cfg.Bucket, name, minio.GetObjectOptions{})
}

func (c minioCli) Delete(ctx context.Context, name string) (err error) {
	return c.cli.RemoveObject(ctx, c.cfg.Bucket, name, minio.RemoveObjectOptions{})
}

func (c minioCli) SignGetUrl(ctx context.Context, name string, expires int64, attachment ...bool) (s string, err error) {
	reqParams := make(url.Values)
	if len(attachment) > 0 && attachment[0] {
		reqParams.Set("response-content-disposition", fmt.Sprintf(`attachment; filename="%s"`, name))
	}

	u, err := c.cli.PresignedGetObject(ctx, c.cfg.Bucket, name, time.Duration(expires)*time.Second, reqParams)
	if err != nil {
		return
	}
	s = u.String()
	return
}

func (c minioCli) SignPutUrl(ctx context.Context, name string, expires int64) (s string, err error) {
	u, err := c.cli.PresignedPutObject(ctx, c.cfg.Bucket, name, time.Duration(expires)*time.Second)
	if err != nil {
		return
	}
	s = u.String()
	return
}
