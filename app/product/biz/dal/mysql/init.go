package mysql

import (
	"fmt"
	"github.com/PiaoAdmin/gomall/app/product/biz/model"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB  *gorm.DB
	err error
)

func Init() {
	//dsn := fmt.Sprintf(conf.GetConf().MySQL.DSN, os.Getenv("MYSQL_USER"), os.Getenv("MYSQL_PASSWORD"), os.Getenv("MYSQL_HOST"), os.Getenv("MYSQL_PORT"))
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/product?charset=utf8mb4&parseTime=True&loc=Local", os.Getenv("MYSQL_USER"), os.Getenv("MYSQL_PASSWORD"), os.Getenv("MYSQL_HOST"), os.Getenv("MYSQL_PORT"))
	fmt.Printf(dsn)
	DB, err = gorm.Open(mysql.Open(dsn),
		&gorm.Config{
			PrepareStmt:            true,
			SkipDefaultTransaction: true,
		},
	)
	if err != nil {
		panic(err)
	}
	err := DB.AutoMigrate(
		model.Product{},
		model.Sku{},
		model.Category{},
	)
	if err != nil {
		return
	}
	err2 := DB.AutoMigrate(model.Sku{})
	if err2 != nil {
		return
	}
}
