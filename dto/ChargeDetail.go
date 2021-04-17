package dto

type ChargeDetail struct {
	ID            uint         `json:"id"`
	AccountId     uint         `json:"account_id"`
	Type          uint8        `json:"type"`
	CategoryId    uint         `json:"category_id"`
	Category      string       `json:"category"`
	Money         float64      `json:"money"`
	Description   string       `json:"description,omitempty"`
	Repay         bool         `json:"repay,omitempty"`
	RepayDetailId uint         `json:"repay_detail_id,omitempty"`
	ReplayDetail  ReplayDetail `json:"replay_detail,omitempty"`
	CreateAt      string       `json:"create_at,omitempty"`
	UpdateAt      string       `json:"update_at,omitempty"`
	RepayAt       string       `json:"repay_at,omitempty"`
}

type ReplayDetail struct {
	ID          uint    `json:"id"`
	AccountId   uint    `json:"account_id"`
	Type        uint8   `json:"type"`
	CategoryId  uint    `json:"category_id"`
	Category    string  `json:"category,omitempty"`
	Money       float64 `json:"money"`
	Description string  `json:"description,omitempty"`
	CreateAt    string  `json:"create_at,omitempty"`
}

type ChargeEditDetail struct {
	ID            uint    `json:"id"`
	AccountId     uint    `json:"account_id"`
	Type          uint8   `json:"type"`
	CategoryId    uint    `json:"category_id"`
	Money         float64 `json:"money"`
	Description   string  `json:"description"`
	Repay         bool    `json:"repay"`
	RepayDetailId uint    `json:"repay_detail_id"`
	RepayAt       int     `json:"repay_at"`
}
