package models

// ChargeSummaryMonth 记账金额月汇总
type ChargeSummaryMonth struct {
	AccountId   uint    `gorm:"primaryKey;"` // 账户id
	Account     Account // orm 的账户详情
	Date        string  `gorm:"primaryKey"` // 按月的日期
	Year        int     `gorm:"index"'`
	CashIn      int64   // 现金收入
	CashOut     int64   // 现金支出
	CreditIn    int64   // 信用账户还款
	CreditOut   int64   // 信用账户借款
	TransferIn  int64   // 转账入
	TransferOut int64   // 转账出
}
