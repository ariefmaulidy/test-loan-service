package loan

import (
	"github.com/ariefmaulidy/test-loan-service/domain/investment"
	loanrequest "github.com/ariefmaulidy/test-loan-service/domain/loan_request"
	"github.com/ariefmaulidy/test-loan-service/domain/document"
	"github.com/ariefmaulidy/test-loan-service/domain/notification"
	"github.com/ariefmaulidy/test-loan-service/domain/storage"
)

type Usecase struct {
	loanDomain       loanrequest.LoanRequestDomainItf
	investmentDomain investment.InvestmentDomainItf
	uploader         storage.Uploader
	emailSender      notification.EmailSender
	agreementGen     document.AgreementGenerator
}

func NewUsecase(
	loanDomain loanrequest.LoanRequestDomainItf,
	investmentDomain investment.InvestmentDomainItf,
	uploader storage.Uploader,
	emailSender notification.EmailSender,
	agreementGen document.AgreementGenerator,
) *Usecase {
	return &Usecase{
		loanDomain:       loanDomain,
		investmentDomain: investmentDomain,
		uploader:         uploader,
		emailSender:      emailSender,
		agreementGen:     agreementGen,
	}
}
