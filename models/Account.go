package models

// Account 账户模型
type Account struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"uniqueIndex;not null"`
	HasCredit uint8  `gorm:"default:0;not null"`
	Cash      int64  `gorm:"default:0;not null"`
	Credit    int64  `gorm:"default:0;not null"`
	Sort      uint8  `gorm:"index;default:0;not null"`
	CreateAt  int    `gorm:"autoCreateTime;default:0;not null"`
	UpdateAt  int    `gorm:"autoUpdateTime;default:0;not null"`
	ChangeAt  int    `gorm:"default:0;not null"`
}
