package dto

type ChargeDetail struct {
	ID              uint          `json:"id"`
	AccountId       uint          `json:"account_id"`
	Account         AccountDetail `json:"account,omitempty"`
	Type            uint8         `json:"type"`
	Category        Category      `json:"category"`
	Money           float64       `json:"money"`
	Description     string        `json:"description,omitempty"`
	RepaidDetail    RepaidDetail  `json:"repaid_detail,omitempty"`
	RepayAccount    AccountDetail `json:"repay_account,omitempty"`
	TransferAccount AccountDetail `json:"transfer_account,omitempty"`
	CreateAt        string        `json:"create_at,omitempty"`
	UpdateAt        string        `json:"update_at,omitempty"`
}

type RepaidDetail struct {
	ID          uint     `json:"id"`
	Money       float64  `json:"money,omitempty"`
	Category    Category `json:"category"`
	Description string   `json:"description,omitempty"`
	CreateAt    string   `json:"create_at,omitempty"`
}

type UnpaidDetail struct {
	RepaidDetail
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
	RepaidDetails     []RepaidDetail `json:"repaid_details"`
}
