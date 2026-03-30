package investment

import (
	"context"
	"errors"
	"sync"
)

type DummyDB struct {
	Investments map[int64]*Investment
	mu          sync.RWMutex
}

func (db *DummyDB) nextInvestmentID() int64 {
	var max int64
	for id := range db.Investments {
		if id > max {
			max = id
		}
	}
	return max + 1
}

func (resource *Resource) CreateInvestmentDB(ctx context.Context, inv Investment) (created Investment, err error) {
	select {
	case <-ctx.Done():
		return created, ctx.Err()
	default:
	}

	resource.db.mu.Lock()
	defer resource.db.mu.Unlock()

	if inv.ID != 0 {
		if _, exists := resource.db.Investments[inv.ID]; exists {
			return created, errors.New("investment already exists")
		}
	} else {
		inv.ID = resource.db.nextInvestmentID()
	}

	resource.db.Investments[inv.ID] = &inv
	return inv, nil
}

func (resource *Resource) GetInvestmentListDB(ctx context.Context, req GetInvestmentParam) (result []Investment, err error) {
	select {
	case <-ctx.Done():
		return result, ctx.Err()
	default:
	}

	resource.db.mu.RLock()
	defer resource.db.mu.RUnlock()

	for _, inv := range resource.db.Investments {
		if req.LoanRequestID != 0 && inv.LoanRequestID != req.LoanRequestID {
			continue
		}

		result = append(result, *inv)
	}

	return result, nil
}
