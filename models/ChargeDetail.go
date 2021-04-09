package models

// ChargeDetail 记账明细
type ChargeDetail struct {
	ID          uint  `gorm:"primaryKey"`
	AccountId   uint  `gorm:"index;"`
	Type        uint8 `gorm:"index"`
	CategoryId  uint  `gorm:"index"`
	Money       int64
	Description string `gorm:"default:''"`
	Repay       uint8
	RepayId     uint
	CreateAt    int `gorm:"autoCreateTime;default:0;not null"`
	UpdateAt    int `gorm:"autoUpdateTime;default:0;not null"`
	RepayAt     int `gorm:"default:0;not null"`
}
