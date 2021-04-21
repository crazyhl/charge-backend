package dto

type ChargeDetail struct {
	ID                uint              `json:"id"`
	AccountId         uint              `json:"account_id"`
	Account           AccountListDetail `json:"account,omitempty"`
	Type              uint8             `json:"type"`
	CategoryId        uint              `json:"category_id"`
	Category          Category          `json:"category"`
	Money             float64           `json:"money"`
	Description       string            `json:"description,omitempty"`
	Repay             bool              `json:"repay,omitempty"`
	RepayDetailId     uint              `json:"repay_detail_id"`
	ReplayDetail      ReplayDetail      `json:"replay_detail,omitempty"`
	RepayAccountId    uint              `json:"repay_account_id"`
	RepayAccount      AccountListDetail `json:"repay_account,omitempty"`
	Transfer          bool              `json:"transfer"`
	TransferAccountId uint              `json:"transfer_account_id"`
	TransferAccount   AccountListDetail `json:"transfer_account,omitempty"`
	CreateAt          string            `json:"create_at,omitempty"`
	UpdateAt          string            `json:"update_at,omitempty"`
	RepayAt           string            `json:"repay_at,omitempty"`
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
	ID                uint    `json:"id"`
	AccountId         uint    `json:"account_id"`
	Type              uint8   `json:"type"`
	CategoryId        uint    `json:"category_id"`
	Money             float64 `json:"money"`
	Description       string  `json:"description"`
	Repay             bool    `json:"repay"`
	RepayDetailId     uint    `json:"repay_detail_id"`
	RepayAt           int64   `json:"repay_at"`
	Transfer          bool    `json:"transfer"`
	TransferAccountId uint    `json:"transfer_account_id"`
}
