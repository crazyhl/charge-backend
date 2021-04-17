package dto

type ChargeDetail struct {
	ID          uint    `json:"id"`
	AccountId   uint    `json:"account_id"`
	Type        uint8   `json:"type"`
	CategoryId  uint    `json:"category_id"`
	Money       float64 `json:"money"`
	Description string  `json:"description"`
	Repay       bool    `json:"repay"`
	RepayId     uint    `json:"repay_id"`
	CreateAt    string  `json:"create_at"`
	UpdateAt    string  `json:"update_at"`
	RepayAt     string  `json:"repay_at"`
}

type ChargeEditDetail struct {
	ID          uint    `json:"id"`
	AccountId   uint    `json:"account_id"`
	Type        uint8   `json:"type"`
	CategoryId  uint    `json:"category_id"`
	Money       float64 `json:"money"`
	Description string  `json:"description"`
	Repay       bool    `json:"repay"`
	RepayId     uint    `json:"repay_id"`
	RepayAt     int     `json:"repay_at"`
}
