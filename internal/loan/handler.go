package loan

import (
	"encoding/json"
	"fmt"
	"loan-service/internal/utils"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type LoanHandler struct {
	loanService IService
}

func NewLoanHandler(loanService IService) *LoanHandler {
	return &LoanHandler{loanService: loanService}
}

func (h *LoanHandler) CreateLoan(w http.ResponseWriter, r *http.Request) {
	var loanRequest CreateLoanRequest
	if err := json.NewDecoder(r.Body).Decode(&loanRequest); err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	err := loanRequest.ValidateCreateLoan()
	if err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	loan := Loan{
		BorrowerID:      loanRequest.BorrowerId,
		PrincipalAmount: loanRequest.PrincipalAmount,
		Rate:            loanRequest.Rate,
		LoanTerm:        loanRequest.LoanTerm,
	}

	if err := h.loanService.CreateLoan(r.Context(), &loan); err != nil {
		utils.SendJSONResponse(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	utils.SendJSONResponse(w, http.StatusCreated, "create loan success", nil)
}

func (h *LoanHandler) GetLoan(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	loan, err := h.loanService.GetLoanByID(r.Context(), id)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusNotFound, err.Error(), nil)
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, "success get loan data", loan)
}

func (h *LoanHandler) ApproveLoan(w http.ResponseWriter, r *http.Request) {
	var approval ApprovalRequest
	err := json.NewDecoder(r.Body).Decode(&approval)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	err = approval.ValidateApprovalRequest()
	if err != nil {
		utils.SendJSONResponse(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	if err := h.loanService.ApproveLoan(r.Context(), approval); err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, "success approve loan", nil)
}

func (h *LoanHandler) InvestLoan(w http.ResponseWriter, r *http.Request) {
	var request InvestmentRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	err := request.ValidateInvestmentRequest()
	if err != nil {
		utils.SendJSONResponse(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	investment := Investment{
		LoanID:     request.LoanId,
		InvestorID: request.InvestorId,
		Amount:     request.Amount,
	}

	if err := h.loanService.InvestLoan(r.Context(), investment); err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, "success invested", nil)
}

func (h *LoanHandler) DisburseLoan(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // Limit upload size to 10MB
	if err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, "failed to parse multipart form", nil)
		return
	}

	file, header, err := r.FormFile("agreement_letter")
	if err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, "failed to retrieve file", nil)
		return
	}
	defer file.Close()

	loanIdRequest := r.FormValue("loan_id")
	fieldOfficerRequest := r.FormValue("field_officer_id")

	if loanIdRequest == "" || fieldOfficerRequest == "" {
		utils.SendJSONResponse(w, http.StatusBadRequest, "missing required fields", nil)
		return
	}

	loanId, err := strconv.ParseInt(loanIdRequest, 10, 64)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, "error convert string to int64", err.Error())
		return
	}

	fieldOfficerId, err := strconv.Atoi(fieldOfficerRequest)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, "error convert string to int", err.Error())
		return
	}

	filePath := fmt.Sprintf("/amr/%s", header.Filename) // Modify storage as needed

	disbursement := Disbursement{
		LoanID:          loanId,
		AgreementLetter: filePath,
		FieldOfficerID:  fieldOfficerId,
	}

	if err := h.loanService.DisburseLoan(r.Context(), disbursement); err != nil {
		utils.SendJSONResponse(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, "Loan successfully disbursed", nil)
}
