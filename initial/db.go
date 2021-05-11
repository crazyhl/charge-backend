package initial

import (
	"charge/container"
	"charge/models"
	"github.com/spf13/viper"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
	"path/filepath"
)

func NewDb() (*gorm.DB, error) {
	dbPath := viper.GetString("db.path")
	if dbPath == "./" {
		executable, getExecutableErr := os.Executable()
		if getExecutableErr != nil {
			panic(getExecutableErr)
		}
		dbPath = filepath.Dir(executable) + "/"
	}
	db, err := gorm.Open(sqlite.Open(dbPath+viper.GetString("db.name")), &gorm.Config{
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
