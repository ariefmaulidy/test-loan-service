package loanrequest

import (
	"context"
)

func (d *LoanRequestDomain) CreateLoanRequest(ctx context.Context, loanReq LoanRequest) (created LoanRequest, err error) {
	select {
	case <-ctx.Done():
		return created, ctx.Err()
	default:
	}

	created, err = d.Resource.CreateLoanRequestDB(ctx, loanReq)
	if err != nil {
		return created, err
	}

	return created, nil
}

func (d *LoanRequestDomain) GetLoanRequestByID(ctx context.Context, id int64) (LoanRequest, error) {
	select {
	case <-ctx.Done():
		return LoanRequest{}, ctx.Err()
	default:
	}

	return d.Resource.GetLoanRequestByIDDB(ctx, id)
}

func (d *LoanRequestDomain) GetLoanRequestList(ctx context.Context, req GetLoanRequestParam) (result []LoanRequest, err error) {
	select {
	case <-ctx.Done():
		return result, ctx.Err()
	default:
	}

	result, err = d.Resource.GetLoanRequestListDB(ctx, req)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (d *LoanRequestDomain) UpdateLoanRequest(ctx context.Context, loanReq LoanRequest) (err error) {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	err = d.Resource.UpdateLoanRequestDB(ctx, loanReq)
	if err != nil {
		return err
	}

	return nil
}
