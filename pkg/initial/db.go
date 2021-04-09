package initial

import (
	"github.com/spf13/viper"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewDb() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(viper.GetString("db.path")+viper.GetString("db.name")), &gorm.Config{})
	return db, err
}
