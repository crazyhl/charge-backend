package dto

type ChargeSummaryMonth struct {
	Account AccountDetail `json:"account"` // orm 的账户详情
	Date    string        `json:"date"`    // 按月的日期
	Year    int           `json:"year"`
	CashIn  float64       `json:"cash_in"`  // 现金收入
	CashOut float64       `json:"cash_out"` // 现金支出
}

type ChargeSummaryCategoryDetail struct {
	Category Category `json:"category"`
	Money    float64  `json:"money"`
	Type     uint8    `json:"type"`
}
