package loanrequest

type LoanRequestDomain struct {
	Resource *Resource
}

type Resource struct {
	db *DummyDB
}

func NewDummyDB() *DummyDB {
	return &DummyDB{
		Loans: make(map[int64]*LoanRequest),
	}
}

func NewLoanRequestDomain(db *DummyDB) *LoanRequestDomain {
	return &LoanRequestDomain{
		Resource: &Resource{db: db},
	}
}