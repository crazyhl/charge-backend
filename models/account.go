package models

// Account 账户模型
type Account struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	Name      string `gorm:"uniqueIndex;not null" json:"name"`
	HasCredit uint8  `gorm:"default:0;not null" json:"has_credit"`
	Cash      int64  `gorm:"default:0;not null" json:"cash"`
	Credit    int64  `gorm:"default:0;not null" json:"credit"`
	Sort      uint8  `gorm:"index;default:0;not null" json:"sort"`
	CreateAt  int64  `gorm:"autoCreateTime" json:"create_at"`
	UpdateAt  int64  `gorm:"autoUpdateTime" json:"update_at"`
	ChangeAt  int64  `gorm:"default:0;not null" json:"change_at"'`
}
