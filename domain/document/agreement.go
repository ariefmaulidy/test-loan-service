package document

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

type AgreementGenerator interface {
	GenerateAgreementLetterURL(ctx context.Context, loanRequestID int64) (string, error)
}

type DummyAgreementGenerator struct {
	rnd *rand.Rand
}

func NewDummyAgreementGenerator() *DummyAgreementGenerator {
	return &DummyAgreementGenerator{rnd: rand.New(rand.NewSource(time.Now().UnixNano()))}
}

func (d *DummyAgreementGenerator) GenerateAgreementLetterURL(ctx context.Context, loanRequestID int64) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	return fmt.Sprintf("https://dummyimage.com/800x1200/000/fff.png&text=Agreement+Loan+%d+Ref+%d", loanRequestID, d.rnd.Intn(999999)), nil
}
