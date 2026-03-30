package loanhandler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	loanrequest "github.com/ariefmaulidy/test-loan-service/domain/loan_request"
	loanuc "github.com/ariefmaulidy/test-loan-service/usecase/loan"
)

type Handler struct {
	uc *loanuc.Usecase
}

func New(uc *loanuc.Usecase) *Handler {
	return &Handler{uc: uc}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/v1/loan-requests", h.handleLoanRequests)
	mux.HandleFunc("/v1/loan-requests/approve", h.handleApproveLoanRequest)
	mux.HandleFunc("/v1/loan-requests/invest", h.handleInvestInLoan)
	mux.HandleFunc("/v1/loan-requests/disburse", h.handleDisburseLoanRequest)
}

func (h *Handler) handleLoanRequests(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CreateLoanRequest(w, r)
	case http.MethodGet:
		h.GetLoanRequestList(w, r)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *Handler) handleApproveLoanRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	h.ApprovedLoanRequest(w, r)
}

func (h *Handler) CreateLoanRequest(w http.ResponseWriter, r *http.Request) {
	var body loanrequest.LoanRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json body")
		return
	}

	created, err := h.uc.CreateLoanRequest(r.Context(), body)
	if err != nil {
		writeUsecaseError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, created)
}

func (h *Handler) GetLoanRequestList(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	currentUserID, err := strconv.ParseInt(q.Get("current_user_id"), 10, 64)
	if err != nil || currentUserID == 0 {
		writeError(w, http.StatusBadRequest, "current_user_id is required")
		return
	}
	role := strings.TrimSpace(q.Get("role"))
	if role == "" {
		writeError(w, http.StatusBadRequest, "role is required")
		return
	}

	var stateFilter []string
	for _, raw := range q["state"] {
		for _, part := range strings.Split(raw, ",") {
			if t := strings.TrimSpace(part); t != "" {
				stateFilter = append(stateFilter, t)
			}
		}
	}

	req := loanuc.GetLoanRequestParam{
		CurrentUserID: currentUserID,
		Role:          role,
		State:         stateFilter,
	}

	list, err := h.uc.GetLoanRequestList(r.Context(), req)
	if err != nil {
		writeUsecaseError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, list)
}

func (h *Handler) ApprovedLoanRequest(w http.ResponseWriter, r *http.Request) {
	const maxMemory = 32 << 20
	if err := r.ParseMultipartForm(maxMemory); err != nil {
		writeError(w, http.StatusBadRequest, "invalid multipart form")
		return
	}

	currentUserID, err := parseFormInt64(r, "current_user_id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "current_user_id is required")
		return
	}
	loanID, err := parseFormInt64(r, "loan_request_id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "loan_request_id is required")
		return
	}

	fh, hdr, err := r.FormFile("proof")
	if err != nil {
		writeError(w, http.StatusBadRequest, "proof file is required")
		return
	}
	defer fh.Close()

	contentType := hdr.Header.Get("Content-Type")
	var proofReader io.Reader = fh
	if contentType == "" {
		var sniff [512]byte
		n, _ := io.ReadFull(fh, sniff[:])
		contentType = http.DetectContentType(sniff[:n])
		proofReader = io.MultiReader(bytes.NewReader(sniff[:n]), fh)
	}

	param := loanuc.ApprovedLoanRequestParam{
		CurrentUserID: currentUserID,
		LoanRequestData: loanrequest.LoanRequest{
			ID: loanID,
		},
		ProofFile:        proofReader,
		ProofFilename:    hdr.Filename,
		ProofSize:        hdr.Size,
		ProofContentType: contentType,
	}

	if err := h.uc.ApprovedLoanRequest(r.Context(), param); err != nil {
		writeUsecaseError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) handleInvestInLoan(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	h.InvestInLoan(w, r)
}

func (h *Handler) handleDisburseLoanRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	h.DisbursedLoanRequest(w, r)
}

func (h *Handler) InvestInLoan(w http.ResponseWriter, r *http.Request) {
	var body loanuc.InvestLoanRequestParam
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json body")
		return
	}

	created, err := h.uc.InvestInLoan(r.Context(), body)
	if err != nil {
		writeUsecaseError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, created)
}

func (h *Handler) DisbursedLoanRequest(w http.ResponseWriter, r *http.Request) {
	const maxMemory = 32 << 20
	if err := r.ParseMultipartForm(maxMemory); err != nil {
		writeError(w, http.StatusBadRequest, "invalid multipart form")
		return
	}

	currentUserID, err := parseFormInt64(r, "current_user_id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "current_user_id is required")
		return
	}
	loanID, err := parseFormInt64(r, "loan_request_id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "loan_request_id is required")
		return
	}

	fh, hdr, err := r.FormFile("signed_agreement")
	if err != nil {
		writeError(w, http.StatusBadRequest, "signed_agreement file is required")
		return
	}
	defer fh.Close()

	contentType := hdr.Header.Get("Content-Type")
	var fileReader io.Reader = fh
	if contentType == "" {
		var sniff [512]byte
		n, _ := io.ReadFull(fh, sniff[:])
		contentType = http.DetectContentType(sniff[:n])
		fileReader = io.MultiReader(bytes.NewReader(sniff[:n]), fh)
	}

	param := loanuc.DisbursedLoanRequestParam{
		CurrentUserID: currentUserID,
		LoanRequestData: loanrequest.LoanRequest{
			ID: loanID,
		},
		SignedAgreementFile:        fileReader,
		SignedAgreementFilename:    hdr.Filename,
		SignedAgreementSize:        hdr.Size,
		SignedAgreementContentType: contentType,
	}

	if err := h.uc.DisbursedLoanRequest(r.Context(), param); err != nil {
		writeUsecaseError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func parseFormInt64(r *http.Request, key string) (int64, error) {
	s := r.FormValue(key)
	if s == "" {
		return 0, errors.New("missing")
	}
	return strconv.ParseInt(s, 10, 64)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func writeUsecaseError(w http.ResponseWriter, err error) {
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		writeError(w, 499, err.Error())
		return
	}
	writeError(w, http.StatusBadRequest, err.Error())
}
