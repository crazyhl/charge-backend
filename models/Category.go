package models

// 绑定类型
// 收入 0 支出 1 借 2 还 3 转 4

// Category 分类
type Category struct {
	ID   uint   `gorm:"primaryKey"`
	Type uint8  `gorm:"not null;uniqueIndex:'uk_type_name'"`
	Name string `gorm:"uniqueIndex:'uk_type_name';not null"`
}
