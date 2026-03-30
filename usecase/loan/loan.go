package loan

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ariefmaulidy/test-loan-service/domain/investment"
	loanrequest "github.com/ariefmaulidy/test-loan-service/domain/loan_request"
)

func (u *Usecase) CreateLoanRequest(ctx context.Context, loanReq loanrequest.LoanRequest) (result loanrequest.LoanRequest, err error) {
	select {
	case <-ctx.Done():
		return result, ctx.Err()
	default:
	}

	if err := validateCreateLoanRequest(&loanReq); err != nil {
		return result, err
	}

	result, err = u.loanDomain.CreateLoanRequest(ctx, loanReq)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (u *Usecase) GetLoanRequestList(ctx context.Context, req GetLoanRequestParam) (result []loanrequest.LoanRequest, err error) {
	select {
	case <-ctx.Done():
		return result, ctx.Err()
	default:
	}

	if req.CurrentUserID == 0 {
		return result, errors.New("current_user_id is required")
	}

	var loanRequestParam loanrequest.GetLoanRequestParam
	// filter loan request by role
	// borrower: only get loan request that belongs to the current user
	// investor: only get loan request that has been approved, invested, or disbursed
	// staff: get all loan requests
	switch req.Role {
	case "borrower":
		loanRequestParam = loanrequest.GetLoanRequestParam{
			BorrowerID: req.CurrentUserID,
			State:      req.State,
		}
	case "investor":
		loanRequestParam = loanrequest.GetLoanRequestParam{
			State: []string{loanrequest.Approved, loanrequest.Invested, loanrequest.Disbursed},
		}
	case "staff":
		loanRequestParam = loanrequest.GetLoanRequestParam{
			State: req.State,
		}
	default:
		return result, errors.New("invalid role")
	}

	result, err = u.loanDomain.GetLoanRequestList(ctx, loanRequestParam)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (u *Usecase) ApprovedLoanRequest(ctx context.Context, param ApprovedLoanRequestParam) (err error) {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	if err := validateApprovedLoanRequestParam(&param); err != nil {
		return err
	}

	url, err := u.UploadFile(ctx, param.ProofFilename, param.ProofFile, param.ProofSize, param.ProofContentType)
	if err != nil {
		return err
	}

	loan, err := u.loanDomain.GetLoanRequestByID(ctx, param.LoanRequestData.ID)
	if err != nil {
		return err
	}

	if loan.State != loanrequest.Proposed {
		return errors.New("only proposed loan request can be approved")
	}

	loan.State = loanrequest.Approved
	loan.ApprovedDate = time.Now()
	loan.Metadata.PictureProofURL = strings.TrimSpace(url)
	loan.Metadata.ApproverEmployeeID = param.CurrentUserID

	return u.loanDomain.UpdateLoanRequest(ctx, loan)
}

func (u *Usecase) DisbursedLoanRequest(ctx context.Context, param DisbursedLoanRequestParam) (err error) {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	if err := validateDisbursedLoanRequestParam(&param); err != nil {
		return err
	}

	url, err := u.UploadFile(ctx, param.SignedAgreementFilename, param.SignedAgreementFile, param.SignedAgreementSize, param.SignedAgreementContentType)
	if err != nil {
		return err
	}

	loan, err := u.loanDomain.GetLoanRequestByID(ctx, param.LoanRequestData.ID)
	if err != nil {
		return err
	}

	if loan.State != loanrequest.Invested {
		return errors.New("only fully invested loan can be disbursed")
	}

	loan.State = loanrequest.Disbursed
	loan.DisbursedDate = time.Now()
	loan.Metadata.SignedAgreementLetterURL = strings.TrimSpace(url)
	loan.Metadata.DisburserEmployeeID = param.CurrentUserID

	return u.loanDomain.UpdateLoanRequest(ctx, loan)
}

func (u *Usecase) InvestInLoan(ctx context.Context, param InvestLoanRequestParam) (created investment.Investment, err error) {
	select {
	case <-ctx.Done():
		return created, ctx.Err()
	default:
	}

	if err := validateInvestLoanRequestParam(&param); err != nil {
		return created, err
	}

	loan, err := u.loanDomain.GetLoanRequestByID(ctx, param.LoanRequestID)
	if err != nil {
		return created, err
	}

	if loan.State == loanrequest.Invested {
		return created, errors.New("loan is already fully invested")
	}
	if loan.State != loanrequest.Approved {
		return created, errors.New("loan is not open for investment")
	}

	existing, err := u.investmentDomain.GetInvestmentList(ctx, investment.GetInvestmentParam{
		LoanRequestID: param.LoanRequestID,
	})
	if err != nil {
		return created, err
	}

	var funded float64
	for _, inv := range existing {
		funded += inv.Amount
	}

	if funded+param.Amount > loan.Principal {
		return created, errors.New("total investment would exceed loan principal")
	}

	created, err = u.investmentDomain.CreateInvestment(ctx, investment.Investment{
		InvestorID:    param.CurrentUserID,
		LoanRequestID: param.LoanRequestID,
		Amount:        param.Amount,
	})
	if err != nil {
		return created, err
	}

	if funded+param.Amount >= loan.Principal {
		agreementURL, err := u.agreementGen.GenerateAgreementLetterURL(ctx, loan.ID)
		if err != nil {
			return created, err
		}
		loan.State = loanrequest.Invested
		loan.Metadata.AgreementLetterURL = agreementURL
		if err := u.loanDomain.UpdateLoanRequest(ctx, loan); err != nil {
			return created, err
		}

		allInvestments, _ := u.investmentDomain.GetInvestmentList(ctx, investment.GetInvestmentParam{
			LoanRequestID: loan.ID,
		})
		for _, inv := range allInvestments {
			_ = u.emailSender.SendEmail(ctx,
				fmt.Sprintf("investor-%d@example.com", inv.InvestorID),
				fmt.Sprintf("Loan #%d Fully Invested - Agreement Letter", loan.ID),
				fmt.Sprintf("Dear Investor %d, the loan #%d is now fully invested. Please find the agreement letter here: %s", inv.InvestorID, loan.ID, agreementURL),
			)
		}
	}

	return created, nil
}
