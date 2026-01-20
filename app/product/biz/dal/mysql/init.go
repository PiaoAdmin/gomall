package mysql

import (
	"github.com/PiaoAdmin/pmall/app/product/biz/model"
	"github.com/PiaoAdmin/pmall/app/product/conf"
	"github.com/cloudwego/kitex/pkg/klog"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	dsn := conf.GetConf().MySQL.DSN
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	if conf.GetEnv() == "test" {
		DB.AutoMigrate(
			&model.ProductSPU{},
			&model.ProductSKU{},
			&model.ProductCategory{},
			&model.ProductBrand{},
		)
	}
	klog.Info("Successfully connected to MySQL")
}
