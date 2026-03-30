package loanrequest

import "time"

const (
	Proposed  string = "proposed"
	Approved  string = "approved"
	Invested  string = "invested"
	Disbursed string = "disbursed"
)

type LoanRequest struct {
	ID            int64        `json:"id"`
	BorrowerID    int64        `json:"borrower_id"`
	Principal     float64      `json:"principal"`
	State         string       `json:"state"`
	Rate          float64      `json:"rate"`
	ROI           float64      `json:"roi"`
	Metadata      MetadataLoan `json:"metadata"`
	ApprovedDate  time.Time    `json:"approved_date"`
	DisbursedDate time.Time    `json:"disbursed_date"`
}

type MetadataLoan struct {
	// Approved metadata
	PictureProofURL    string `json:"picture_proof_url"`
	ApproverEmployeeID int64  `json:"approver_employee_id"`

	// Invested metadata
	AgreementLetterURL string `json:"agreement_letter_url"`

	// Disbursed metadata
	SignedAgreementLetterURL string `json:"signed_agreement_letter_url"`
	DisburserEmployeeID      int64  `json:"disburser_employee_id"`
}

type GetLoanRequestParam struct {
	BorrowerID int64    `json:"borrower_id"`
	State      []string `json:"state"`
}
