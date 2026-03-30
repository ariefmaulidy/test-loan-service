package storage

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"time"
)

type Uploader interface {
	Upload(ctx context.Context, filename string, r io.Reader, size int64, contentType string) (publicURL string, err error)
}

type DummyUploader struct {
	rnd *rand.Rand
}

func NewDummyUploader() *DummyUploader {
	return &DummyUploader{rnd: rand.New(rand.NewSource(time.Now().UnixNano()))}
}

func (d *DummyUploader) Upload(ctx context.Context, filename string, r io.Reader, size int64, contentType string) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	_, _ = io.Copy(io.Discard, r)
	return randomDummyImageURL(d.rnd), nil
}

func randomDummyImageURL(rnd *rand.Rand) string {
	w := 300 + rnd.Intn(501)
	h := 200 + rnd.Intn(401)
	bg := fmt.Sprintf("%06x", rnd.Intn(0x1000000))
	fg := fmt.Sprintf("%06x", rnd.Intn(0x1000000))
	exts := []string{".png", ".jpg", ".gif"}
	ext := exts[rnd.Intn(len(exts))]
	return fmt.Sprintf("https://dummyimage.com/%dx%d/%s/%s%s", w, h, bg, fg, ext)
}
