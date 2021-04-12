package models

// Account 账户模型
type Account struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"uniqueIndex;not null"`
	HasCredit uint8  `gorm:"default:0;not null"`
	Cash      int64  `gorm:"default:0;not null"`
	Credit    int64  `gorm:"default:0;not null"`
	Sort      uint8  `gorm:"index;default:0;not null"`
	CreateAt  int64  `gorm:"autoCreateTime"`
	UpdateAt  int64  `gorm:"autoUpdateTime"`
	ChangeAt  int64  `gorm:"default:0;not null"`
}
