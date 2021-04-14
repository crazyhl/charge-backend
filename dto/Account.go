package dto

// Account 账户模型
type Account struct {
	ID        uint    `json:"id"`
	Name      string  `json:"name"`
	HasCredit bool    `json:"has_credit"`
	Cash      float64 `json:"cash"`
	Credit    float64 `json:"credit"`
	Sort      uint8   `json:"sort"`
	CreateAt  string  `json:"create_at"`
	UpdateAt  string  `json:"update_at"`
	ChangeAt  string  `json:"change_at"'`
}
