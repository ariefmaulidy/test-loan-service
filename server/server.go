package server

import (
	"log"

	"github.com/ariefmaulidy/test-loan-service/domain/document"
	"github.com/ariefmaulidy/test-loan-service/domain/investment"
	loanrequest "github.com/ariefmaulidy/test-loan-service/domain/loan_request"
	"github.com/ariefmaulidy/test-loan-service/domain/notification"
	"github.com/ariefmaulidy/test-loan-service/domain/storage"
	"github.com/ariefmaulidy/test-loan-service/usecase/loan"
)

type CommonDependency struct {
	LoanUsecase *loan.Usecase
}

func InitDependencies() *CommonDependency {
	loanDB := loanrequest.NewDummyDB()
	loanDom := loanrequest.NewLoanRequestDomain(loanDB)
	if loanDom == nil {
		log.Fatal("failed to init loan domain")
	}

	invDB := investment.NewDummyDB()
	invDom := investment.NewInvestmentDomain(invDB)
	if invDom == nil {
		log.Fatal("failed to init investment domain")
	}

	uploader := storage.NewDummyUploader()
	if uploader == nil {
		log.Fatal("failed to init uploader")
	}

	emailSender := notification.NewDummyEmailSender()
	if emailSender == nil {
		log.Fatal("failed to init email sender")
	}

	agreementGen := document.NewDummyAgreementGenerator()
	if agreementGen == nil {
		log.Fatal("failed to init agreement generator")
	}

	uc := loan.NewUsecase(loanDom, invDom, uploader, emailSender, agreementGen)

	return &CommonDependency{
		LoanUsecase: uc,
	}
}
