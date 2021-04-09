package container

import (
	"gorm.io/gorm"
)

type container struct {
	db *gorm.DB
}

func (c *container) SetDb(db *gorm.DB)  {
	c.db = db
}

func (c *container) GetDb() *gorm.DB {
	return c.db
}


// ---------------- 外部使用部分 ------------------------
var c *container

func GetContainer() *container {
	if c == nil {
		c = new(container)
	}

	return c
}