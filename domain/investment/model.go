package investment

type Investment struct {
	ID            int64   `json:"id"`
	InvestorID    int64   `json:"investor_id"`
	LoanRequestID int64   `json:"loan_request_id"`
	Amount        float64 `json:"amount"`
}

type GetInvestmentParam struct {
	LoanRequestID int64  `json:"loan_request_id"`
	State         string `json:"state"`
}
