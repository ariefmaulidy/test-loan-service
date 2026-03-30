package loan

import (
	"io"
	"slices"

	loanrequest "github.com/ariefmaulidy/test-loan-service/domain/loan_request"
)

var AllowedApprovedPictureProofMimes = []string{
	"image/jpeg",
	"image/png",
	"image/webp",
}

var AllowedDisbursedSignedAgreementMimes = []string{
	"application/pdf",
	"image/jpeg",
	"image/jpg",
}

func IsAllowedApprovedPictureProofMime(mime string) bool {
	return slices.Contains(AllowedApprovedPictureProofMimes, mime)
}

func IsAllowedDisbursedSignedAgreementMime(mime string) bool {
	return slices.Contains(AllowedDisbursedSignedAgreementMimes, mime)
}

type GetLoanRequestParam struct {
	CurrentUserID int64    `json:"current_user_id"`
	Role          string   `json:"role"`
	State         []string `json:"state"`
}

type ApprovedLoanRequestParam struct {
	CurrentUserID   int64                   `json:"current_user_id"`
	LoanRequestData loanrequest.LoanRequest `json:"loan_request_data"`

	ProofFile        io.Reader `json:"-"`
	ProofFilename    string    `json:"proof_filename,omitempty"`
	ProofSize        int64     `json:"proof_size,omitempty"`
	ProofContentType string    `json:"proof_content_type,omitempty"`
}

type DisbursedLoanRequestParam struct {
	CurrentUserID   int64                   `json:"current_user_id"`
	LoanRequestData loanrequest.LoanRequest `json:"loan_request_data"`

	SignedAgreementFile        io.Reader `json:"-"`
	SignedAgreementFilename    string    `json:"signed_agreement_filename,omitempty"`
	SignedAgreementSize        int64     `json:"signed_agreement_size,omitempty"`
	SignedAgreementContentType string    `json:"signed_agreement_content_type,omitempty"`
}

type InvestLoanRequestParam struct {
	CurrentUserID int64   `json:"current_user_id"`
	LoanRequestID int64   `json:"loan_request_id"`
	Amount        float64 `json:"amount"`
}
