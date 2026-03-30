package loan

import (
	"context"
	"io"
)

func (u *Usecase) UploadFile(ctx context.Context, filename string, r io.Reader, size int64, contentType string) (publicURL string, err error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}
	return u.uploader.Upload(ctx, filename, r, size, contentType)
}
