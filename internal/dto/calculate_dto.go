package dto

type CalculateRequest struct {
	Amount float64 `json:"amount" binding:"required"`
}

type CalculationResult struct {
	Tenor              int     `json:"tenor"`
	MonthlyInstallment float64 `json:"monthly_installment"`
	TotalMargin        float64 `json:"total_margin"`
	TotalPayment       float64 `json:"total_payment"`
}

type CalculateResponse struct {
	Calculations []CalculationResult `json:"calculations"`
}
