<<<<<<< HEAD
/*
 * @Author: liaosijie
 * @Date: 2025-02-18 16:47:32
 * @Last Modified by: liaosijie
 * @Last Modified time: 2025-02-18 17:08:47
 */

package mysql

import (
	//"douyin-gomall/gomall/app/order/biz/model"
	//"douyin-gomall/gomall/app/order/conf"
	"fmt"
	"os"

	"github.com/PiaoAdmin/gomall/app/order/biz/model"
	"github.com/PiaoAdmin/gomall/app/order/conf"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"k8s.io/klog/v2"
=======
package mysql

import (
	"github.com/PiaoAdmin/gomall/app/order/conf"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
>>>>>>> b6e73c27fce12b01552c5334097a847176b8f26a
)

var (
	DB  *gorm.DB
	err error
)

func Init() {
<<<<<<< HEAD
	dsn :=fmt.Sprint(conf.GetConf().MySQL.DSN, os.Getenv("MYSQL_USER"), os.Getenv("MYSQL_PASSWORD"), os.Getenv("MYSQL_HOST"))
	DB, err = gorm.Open(mysql.Open(dsn),
=======
	DB, err = gorm.Open(mysql.Open(conf.GetConf().MySQL.DSN),
>>>>>>> b6e73c27fce12b01552c5334097a847176b8f26a
		&gorm.Config{
			PrepareStmt:            true,
			SkipDefaultTransaction: true,
		},
	)
	if err != nil {
		panic(err)
	}
<<<<<<< HEAD
	if os.Getenv("GO_ENV") != "online"{
		if err := DB.AutoMigrate(&model.Order{},&model.OrderItem{}); err!= nil {
			klog.Error(err)
		}
	}
=======
>>>>>>> b6e73c27fce12b01552c5334097a847176b8f26a
}
