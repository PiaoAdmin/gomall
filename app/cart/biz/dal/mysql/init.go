package mysql

import (
	"fmt"
	"os"

	"github.com/PiaoAdmin/gomall/app/cart/biz/dal/model"
	"github.com/PiaoAdmin/gomall/app/cart/conf"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB  *gorm.DB
	err error
)

func Init() {
	dsn := fmt.Sprintf(conf.GetConf().MySQL.DSN,
		os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("MYSQL_HOST"),
	)
	//dsn := "root:jinitaimei114514@tcp(localhost:3306)/cart?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn),
		&gorm.Config{
			PrepareStmt:            true,
			SkipDefaultTransaction: true,
		},
	)
	DB.AutoMigrate(&model.Cart{})
	if err != nil {
		panic(err)
	}
}
