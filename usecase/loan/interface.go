package loan

import (
	"context"
	"io"

	"github.com/ariefmaulidy/test-loan-service/domain/investment"
	loanrequest "github.com/ariefmaulidy/test-loan-service/domain/loan_request"
)

type UsecaseItf interface {
	CreateLoanRequest(ctx context.Context, loanReq loanrequest.LoanRequest) (result loanrequest.LoanRequest, err error)
	GetLoanRequestList(ctx context.Context, req GetLoanRequestParam) (result []loanrequest.LoanRequest, err error)
	ApprovedLoanRequest(ctx context.Context, param ApprovedLoanRequestParam) (err error)
	DisbursedLoanRequest(ctx context.Context, param DisbursedLoanRequestParam) (err error)
	UpdateLoanRequest(ctx context.Context, loanReq loanrequest.LoanRequest) (err error)

	InvestInLoan(ctx context.Context, param InvestLoanRequestParam) (result investment.Investment, err error)
	CreateInvestment(ctx context.Context, inv investment.Investment) (result investment.Investment, err error)
	GetInvestmentList(ctx context.Context, req investment.GetInvestmentParam) (result []investment.Investment, err error)

	UploadFile(ctx context.Context, filename string, r io.Reader, size int64, contentType string) (publicURL string, err error)
}
