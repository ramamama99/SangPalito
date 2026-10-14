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
