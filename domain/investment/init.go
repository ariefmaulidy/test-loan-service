package investment

type InvestmentDomain struct {
	Resource *Resource
}

type Resource struct {
	db *DummyDB
}

func NewDummyDB() *DummyDB {
	return &DummyDB{
		Investments: make(map[int64]*Investment),
	}
}

func NewInvestmentDomain(db *DummyDB) *InvestmentDomain {
	return &InvestmentDomain{
		Resource: &Resource{db: db},
	}
}
