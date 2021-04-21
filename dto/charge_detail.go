package dto

type ChargeDetail struct {
	ID                uint          `json:"id"`
	AccountId         uint          `json:"account_id"`
	Account           AccountDetail `json:"account,omitempty"`
	Type              uint8         `json:"type"`
	CategoryId        uint          `json:"category_id"`
	Category          Category      `json:"category"`
	Money             float64       `json:"money"`
	Description       string        `json:"description,omitempty"`
	RepayDetailId     uint          `json:"repay_detail_id"`
	ReplayDetail      ReplayDetail  `json:"replay_detail,omitempty"`
	RepayAccountId    uint          `json:"repay_account_id"`
	RepayAccount      AccountDetail `json:"repay_account,omitempty"`
	TransferAccountId uint          `json:"transfer_account_id"`
	TransferAccount   AccountDetail `json:"transfer_account,omitempty"`
	CreateAt          string        `json:"create_at,omitempty"`
	UpdateAt          string        `json:"update_at,omitempty"`
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
	ID                uint           `json:"id"`
	AccountId         uint           `json:"account_id"`
	Type              uint8          `json:"type"`
	CategoryId        uint           `json:"category_id"`
	Money             float64        `json:"money"`
	Description       string         `json:"description"`
	RepayAccountId    uint           `json:"repay_account_id"`
	TransferAccountId uint           `json:"transfer_account_id"`
	RepaidDetails     []ReplayDetail `json:"repaid_details"`
}
