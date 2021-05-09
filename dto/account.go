package dto

// AccountDetail 账户模型
type AccountDetail struct {
	ID        uint    `json:"id"`
	Name      string  `json:"name"`
	HasCredit bool    `json:"has_credit,omitempty"`
	Cash      float64 `json:"cash,omitempty"`
	Credit    float64 `json:"credit,omitempty"`
	Sort      uint8   `json:"sort,omitempty"`
	CreateAt  string  `json:"create_at,omitempty"`
	UpdateAt  string  `json:"update_at,omitempty"`
	ChangeAt  string  `json:"change_at,omitempty"'`
}

type AccountEditDetail struct {
	ID        uint    `json:"id"`
	Name      string  `json:"name"`
	HasCredit bool    `json:"has_credit"`
	Cash      float64 `json:"cash"`
	Credit    float64 `json:"credit"`
	Sort      uint8   `json:"sort"`
}
