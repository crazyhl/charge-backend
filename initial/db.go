package initial

import (
	"charge/container"
	"charge/models"
	"github.com/spf13/viper"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDb() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(viper.GetString("db.path")+viper.GetString("db.name")), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   logger.Default.LogMode(logger.Info),
	})
	return db, err
}

func AutoMigrate() error {
	db := container.GetContainer().GetDb()
	err := db.AutoMigrate(
		models.Account{},
		models.Category{},
		models.ChargeDetail{},
		models.ChargeSummaryMonth{},
	)

	return err
}
