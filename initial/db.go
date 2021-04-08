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
		DisableDatetimePrecision: true, // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex: true, // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn: true, // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{})

	return db, err
}
