package models

// Category 分类
type Category struct {
	ID   uint   `gorm:"primaryKey"`
	Type uint8  `gorm:"not null;uniqueIndex:'uk_type_name'"`
	Name string `gorm:"uniqueIndex:'uk_type_name';not null"`
}
