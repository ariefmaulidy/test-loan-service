package loanrequest

import "context"

type LoanRequestDomainItf interface {
	CreateLoanRequest(ctx context.Context, loanReq LoanRequest) (created LoanRequest, err error)
	GetLoanRequestByID(ctx context.Context, id int64) (LoanRequest, error)
	GetLoanRequestList(ctx context.Context, req GetLoanRequestParam) (result []LoanRequest, err error)
	UpdateLoanRequest(ctx context.Context, loanReq LoanRequest) (err error)
}

type LoanRequestResourceItf interface {
	CreateLoanRequestDB(ctx context.Context, loanReq LoanRequest) (created LoanRequest, err error)
	GetLoanRequestByIDDB(ctx context.Context, id int64) (LoanRequest, error)
	GetLoanRequestListDB(ctx context.Context, req GetLoanRequestParam) (result []LoanRequest, err error)
	UpdateLoanRequestDB(ctx context.Context, loanReq LoanRequest) (err error)
}
