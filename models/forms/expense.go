package forms

import "bfimpl/models/oa"

type ParseExpenseDetailResponse struct {
	Details  []*oa.ExpenseDetail `json:"details"`
	FileName string              `json:"file_name"`
}

type PaidCardInfo struct {
	BankName string
	CardID   string
	//UserName       string
	//ExpenseSummary float64
	PaymentName string
}
