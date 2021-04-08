package initial

import (
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewDb() (*gorm.DB, error) {
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: viper.GetString("mysql.username") + ":" +
			viper.GetString("mysql.password") +"@tcp(" +
			viper.GetString("mysql.host") + ":" +
			viper.GetString("mysql.port") + ")/" +
			viper.GetString("mysql.database") + "?charset=" +
			viper.GetString("mysql.charset") + "&parseTime=" +
			viper.GetString("mysql.parseTime") + "&loc=" +
			viper.GetString("mysql.location"), // DSN data source name
		DefaultStringSize: 191, // string 类型字段的默认长度
		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{})

	return db, err
}
