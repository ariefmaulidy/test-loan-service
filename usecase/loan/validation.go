package loan

import (
	"errors"
	"strings"

	loanrequest "github.com/ariefmaulidy/test-loan-service/domain/loan_request"
)

func validateApprovedLoanRequestParam(p *ApprovedLoanRequestParam) error {
	if p.CurrentUserID == 0 {
		return errors.New("current_user_id is required")
	}
	if p.LoanRequestData.ID == 0 {
		return errors.New("loan_request id is required")
	}
	if p.ProofFile == nil {
		return errors.New("picture proof file is required")
	}
	if strings.TrimSpace(p.ProofContentType) == "" {
		return errors.New("picture proof content type is required")
	}
	if !IsAllowedApprovedPictureProofMime(p.ProofContentType) {
		return errors.New("invalid picture proof content type")
	}
	return nil
}

func validateDisbursedLoanRequestParam(p *DisbursedLoanRequestParam) error {
	if p.CurrentUserID == 0 {
		return errors.New("current_user_id is required")
	}
	if p.LoanRequestData.ID == 0 {
		return errors.New("loan_request id is required")
	}
	if p.SignedAgreementFile == nil {
		return errors.New("signed agreement file is required")
	}
	if strings.TrimSpace(p.SignedAgreementContentType) == "" {
		return errors.New("signed agreement content type is required")
	}
	if !IsAllowedDisbursedSignedAgreementMime(p.SignedAgreementContentType) {
		return errors.New("invalid signed agreement content type, must be pdf or jpeg")
	}
	return nil
}

func validateInvestLoanRequestParam(p *InvestLoanRequestParam) error {
	if p.CurrentUserID == 0 {
		return errors.New("current_user_id is required")
	}
	if p.LoanRequestID == 0 {
		return errors.New("loan_request_id is required")
	}
	if p.Amount <= 0 {
		return errors.New("amount must be greater than zero")
	}
	return nil
}

func validateCreateLoanRequest(loanReq *loanrequest.LoanRequest) error {
	if loanReq.BorrowerID == 0 {
		return errors.New("borrower_id is required")
	}
	if loanReq.Principal == 0 {
		return errors.New("principal is required")
	}
	if loanReq.Rate == 0 {
		loanReq.Rate = 10
	}
	if loanReq.ROI == 0 {
		loanReq.ROI = 8
	}
	if loanReq.State == "" {
		loanReq.State = loanrequest.Proposed
	}
	return nil
}
