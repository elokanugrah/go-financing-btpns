package dto

type SubmitFinancingRequest struct {
	UserID          int64   `json:"user_id"`
	FacilityLimitID int64   `json:"facility_limit_id"`
	Amount          float64 `json:"amount"`
	Tenor           int     `json:"tenor"`
	StartDate       string  `json:"start_date"`
}

type ScheduleItem struct {
	DueDate           string  `json:"due_date"`
	InstallmentAmount float64 `json:"installment_amount"`
}

type SubmitFinancingResponse struct {
	UserID          int64          `json:"user_id"`
	FacilityLimitID int64          `json:"facility_limit_id"`
	Amount          float64        `json:"amount"`
	Tenor           int            `json:"tenor"`
	StartDate       string         `json:"start_date"`
	MonthlyInstall  float64        `json:"monthly_installment"`
	TotalMargin     float64        `json:"total_margin"`
	TotalPayment    float64        `json:"total_payment"`
	Schedule        []ScheduleItem `json:"schedule"`
}
