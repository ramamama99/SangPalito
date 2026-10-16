package dto

type AnalyticsSummaryResponse struct {
	TotalVisits   int64   `json:"total_visits"`
	TotalPatients int64   `json:"total_patients"`
	AvgWaitMin    float64 `json:"avg_wait_min"`
	AvgServiceMin float64 `json:"avg_service_min"`
}

type HourlyDistribution struct {
	Hour  int   `json:"hour"`
	Total int64 `json:"total"`
}

type BusyHoursResponse struct {
	Data []HourlyDistribution `json:"data"`
}

type TopDiagnosisItem struct {
	ICD10Code   string `json:"icd10_code"`
	Description string `json:"description"`
	Total       int64  `json:"total"`
}

type TopDiagnosesResponse struct {
	Data []TopDiagnosisItem `json:"data"`
}

type DailyPatientCount struct {
	Date    string `json:"date"`
	DayName string `json:"day_name"`
	Count   int64  `json:"count"`
}

type WeeklyPatientsResponse struct {
	TotalThisWeek  int64               `json:"total_this_week"`
	TotalLastWeek  int64               `json:"total_last_week"`
	GrowthPct      float64             `json:"growth_pct"`
	DailyBreakdown []DailyPatientCount `json:"daily_breakdown"`
	StatusCounts   map[string]int64    `json:"status_counts"`
}

type DoctorPerformanceItem struct {
	DoctorID          string  `json:"doctor_id"`
	DoctorName        string  `json:"doctor_name"`
	Specialty         string  `json:"specialty"`
	PolyclinicName    string  `json:"polyclinic_name"`
	TotalHandled      int64   `json:"total_handled"`
	AvgServiceMin     float64 `json:"avg_service_min"`
	CompletionRatePct float64 `json:"completion_rate_pct"`
	SatisfactionScore float64 `json:"satisfaction_score"`
}

type DoctorPerformanceResponse struct {
	Data []DoctorPerformanceItem `json:"data"`
}
