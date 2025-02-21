/**
 * @Author: ZhangHaoChen
 * @Date:   2/19/25 PM5:07
 */

package model

import (
	"time"

	"github.com/bwmarrin/snowflake"
	"gorm.io/gorm"
	"log"
)

// CustomModel 自定义模型结构体
type CustomModel struct {
	ID        int64          `gorm:"primaryKey;autoIncrement:false"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// BeforeCreate 在创建记录之前调用
func (cm *CustomModel) BeforeCreate(tx *gorm.DB) (err error) {
	cm.ID = CreateId(700) // 这里的1可以替换为实际的mark值
	return
}
func CreateId(mark int64) (id int64) {
	node, err := snowflake.NewNode(mark) // mark表示节点ID，每个人要不一样
	if err != nil {
		log.Fatalf("failed to create snowflake node: %v", err)
	}
	id = node.Generate().Int64()
	return
}
