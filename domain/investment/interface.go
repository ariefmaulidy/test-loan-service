package investment

import "context"

type InvestmentDomainItf interface {
	CreateInvestment(ctx context.Context, inv Investment) (created Investment, err error)
	GetInvestmentList(ctx context.Context, req GetInvestmentParam) (result []Investment, err error)
}

type InvestmentResourceItf interface {
	CreateInvestmentDB(ctx context.Context, inv Investment) (created Investment, err error)
	GetInvestmentListDB(ctx context.Context, req GetInvestmentParam) (result []Investment, err error)
}
