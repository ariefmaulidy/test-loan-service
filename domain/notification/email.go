package notification

import (
	"context"
	"fmt"
	"log"
)

type EmailSender interface {
	SendEmail(ctx context.Context, to string, subject string, body string) error
}

type DummyEmailSender struct{}

func NewDummyEmailSender() *DummyEmailSender {
	return &DummyEmailSender{}
}

func (d *DummyEmailSender) SendEmail(ctx context.Context, to string, subject string, body string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	fmt.Printf("[DUMMY EMAIL] To: %s | Subject: %s | Body: %s\n", to, subject, body)
	log.Printf("[DUMMY EMAIL] sent to %s", to)
	return nil
}
