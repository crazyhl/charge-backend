package models

// ChargeDetail 记账明细
type ChargeDetail struct {
	ID          uint     `gorm:"primaryKey"`
	AccountId   uint     `gorm:"index;"` // 账户id
	Account     Account  // orm 的账户详情
	Type        uint8    `gorm:"index"` // 类型 收入 支出 借 还 转
	CategoryId  uint     `gorm:"index"` // 分类
	Category    Category // orm 的分类信息
	Money       int64    // 金额
	Description string   `gorm:"default:''"` // 描述
	// -------------- 下面两个是给借款 type 用的
	RepayDetailId uint          `gorm:"index"` // 对应还款 明细 id
	RepayDetail   *ChargeDetail // 还款详情
	// -------------- 下面两个是给 还款 type 用的
	RepayAccountId uint    `gorm:"index"` // 还款账户id
	RepayAccount   Account // orm 还款账户信息
	// ------------- 下面两个是给转账 type 用的
	TransferAccountId uint    `gorm:"index"` // 转账 账户id
	TransferAccount   Account // orm 的账户信息
	CreateAt          int64   `gorm:"autoCreateTime;default:0;not null"`
	UpdateAt          int64   `gorm:"autoUpdateTime;default:0;not null"`
}
