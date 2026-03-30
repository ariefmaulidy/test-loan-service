package investment

import (
	"context"
)

func (d *InvestmentDomain) CreateInvestment(ctx context.Context, inv Investment) (created Investment, err error) {
	select {
	case <-ctx.Done():
		return created, ctx.Err()
	default:
	}

	created, err = d.Resource.CreateInvestmentDB(ctx, inv)
	if err != nil {
		return created, err
	}

	return created, nil
}

func (d *InvestmentDomain) GetInvestmentList(ctx context.Context, req GetInvestmentParam) (result []Investment, err error) {
	select {
	case <-ctx.Done():
		return result, ctx.Err()
	default:
	}

	result, err = d.Resource.GetInvestmentListDB(ctx, req)
	if err != nil {
		return result, err
	}

	return result, nil
}
