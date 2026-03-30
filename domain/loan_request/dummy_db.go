package loanrequest

import (
	"context"
	"errors"
	"slices"
	"sync"
)

type DummyDB struct {
	Loans map[int64]*LoanRequest
	mu    sync.RWMutex
}

func (db *DummyDB) nextLoanID() int64 {
	var max int64
	for id := range db.Loans {
		if id > max {
			max = id
		}
	}
	return max + 1
}

func (resource *Resource) CreateLoanRequestDB(ctx context.Context, loanReq LoanRequest) (created LoanRequest, err error) {
	select {
	case <-ctx.Done():
		return created, ctx.Err()
	default:
	}

	resource.db.mu.Lock()
	defer resource.db.mu.Unlock()

	if loanReq.ID != 0 {
		if _, exists := resource.db.Loans[loanReq.ID]; exists {
			return created, errors.New("loan request already exists")
		}
	} else {
		loanReq.ID = resource.db.nextLoanID()
	}

	resource.db.Loans[loanReq.ID] = &loanReq
	return loanReq, nil
}

func (resource *Resource) GetLoanRequestByIDDB(ctx context.Context, id int64) (LoanRequest, error) {
	select {
	case <-ctx.Done():
		return LoanRequest{}, ctx.Err()
	default:
	}

	resource.db.mu.RLock()
	defer resource.db.mu.RUnlock()

	loan, ok := resource.db.Loans[id]
	if !ok {
		return LoanRequest{}, errors.New("loan request not found")
	}
	return *loan, nil
}

func (resource *Resource) GetLoanRequestListDB(ctx context.Context, req GetLoanRequestParam) (result []LoanRequest, err error) {
	select {
	case <-ctx.Done():
		return result, ctx.Err()
	default:
	}

	resource.db.mu.RLock()
	defer resource.db.mu.RUnlock()

	for _, loan := range resource.db.Loans {
		// filter borrowerID (optional)
		if req.BorrowerID != 0 && loan.BorrowerID != req.BorrowerID {
			continue
		}

		// filter state (optional): kosong = semua state; ada isi = loan harus salah satu dari slice
		if len(req.State) > 0 && !slices.Contains(req.State, loan.State) {
			continue
		}

		result = append(result, *loan)
	}

	return result, nil
}

func (resource *Resource) UpdateLoanRequestDB(ctx context.Context, loanReq LoanRequest) (err error) {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	resource.db.mu.Lock()
	defer resource.db.mu.Unlock()

	if _, ok := resource.db.Loans[loanReq.ID]; !ok {
		return errors.New("loan request not found")
	}
	resource.db.Loans[loanReq.ID] = &loanReq

	return nil
}
