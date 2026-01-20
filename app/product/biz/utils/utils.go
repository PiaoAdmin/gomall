package utils

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/PiaoAdmin/pmall/common/errs"
)

// ==================== 商品状态 ====================
const (
	PublishStatusPublished   int8 = 1 // 上架
	PublishStatusUnpublished int8 = 2 // 下架

	VerifyStatusUnverified  int8 = 0 // 未审核
	VerifyStatusApproved    int8 = 1 // 审核通过
	VerifyStatusNotApproved int8 = 2 // 审核不通过
)

func IsValidPublishStatus(status int8) bool {
	return status == PublishStatusPublished || status == PublishStatusUnpublished
}

func IsValidVerifyStatus(status int8) bool {
	return status == VerifyStatusUnverified || status == VerifyStatusApproved || status == VerifyStatusNotApproved
}

// ==================== 商品服务Bit ====================
const (
	ServiceSevenDayReturn     int64 = 1 << iota // 第0位：7天无理由退货
	ServiceFifteenDayExchange                   // 第1位：15天换货
	ServiceOnSiteInstall                        // 第2位：上门安装
	ServiceExtendedWarranty                     // 第3位：延保服务
	ServiceTradeIn                              // 第4位：以旧换新
	ServiceFast24Ship                           // 第5位：24小时发货
	ServiceShippingInsurance                    // 第6位：运费险
	ServiceAuthenticGuarantee                   // 第7位：正品保障
	ServiceFreeRepair                           // 第8位：免费维修
	ServiceInstallmentFree                      // 第9位：分期免息
)

// HasService 检查商品是否包含某项服务
func HasService(serviceBits int64, service int64) bool {
	return serviceBits&service != 0
}

// AddService 为商品添加服务
func AddService(serviceBits int64, service int64) int64 {
	return serviceBits | service
}

// RemoveService 移除商品服务
func RemoveService(serviceBits int64, service int64) int64 {
	return serviceBits &^ service
}

// GetCommonServices 获取常用服务组合
func GetCommonServices() int64 {
	return ServiceSevenDayReturn | ServiceAuthenticGuarantee | ServiceShippingInsurance
}

// GetPremiumServices 获取高端商品服务组合
func GetPremiumServices() int64 {
	return ServiceSevenDayReturn | ServiceFifteenDayExchange | ServiceExtendedWarranty |
		ServiceTradeIn | ServiceFast24Ship | ServiceShippingInsurance |
		ServiceAuthenticGuarantee | ServiceFreeRepair | ServiceInstallmentFree

}

// ==================== 品牌常量 ====================
const (
	BrandXiaomi uint64 = 1 // 小米
	BrandRedmi  uint64 = 2 // 红米
	BrandMijia  uint64 = 3 // 米家
)

// ==================== 一级分类常量 ====================
const (
	CategoryPhone    uint64 = 1 // 手机
	CategoryTV       uint64 = 2 // 电视
	CategoryHomeAppl uint64 = 3 // 家电
	CategoryLaptop   uint64 = 4 // 笔记本
	CategoryTablet   uint64 = 5 // 平板
	CategoryEarphone uint64 = 6 // 耳机
	CategoryRouter   uint64 = 7 // 路由器
)

// ==================== 二级分类常量 - 手机 ====================
const (
	CategoryPhoneXiaomi uint64 = 101 // 小米手机
	CategoryPhoneRedmi  uint64 = 102 // 红米手机
)

// ==================== 二级分类常量 - 家电 ====================
const (
	CategoryWallAC          uint64 = 301 // 壁挂空调
	CategoryFloorAC         uint64 = 302 // 立式空调
	CategoryCentralACPro    uint64 = 303 // 中央空调Pro
	CategoryRefrigerator    uint64 = 304 // 冰箱
	CategoryDrumWasher      uint64 = 305 // 滚筒洗衣机
	CategoryWaveWasher      uint64 = 306 // 波轮洗衣机
	CategoryHeater          uint64 = 307 // 电暖器
	CategoryDehumidifier    uint64 = 308 // 除湿机
	CategoryFloorWasher     uint64 = 309 // 洗地机
	CategoryWaterPurifier   uint64 = 310 // 净水器
	CategorySteamOven       uint64 = 311 // 微蒸烤
	CategoryHoodStove       uint64 = 312 // 烟灶
	CategoryDishwasher      uint64 = 313 // 洗碗机
	CategoryRobotVacuum     uint64 = 314 // 扫地机器人
	CategoryVacuumCleaner   uint64 = 315 // 吸尘器
	CategoryHumidifier      uint64 = 316 // 加湿器
	CategoryAirPurifier     uint64 = 317 // 空气净化器
	CategoryRiceCooker      uint64 = 318 // 电饭煲
	CategoryInductionCooker uint64 = 319 // 电磁炉
	CategoryKettle          uint64 = 320 // 水壶
	CategoryStandingFan     uint64 = 321 // 落地风扇
	CategoryProjector       uint64 = 322 // 投影仪
	CategoryLighting        uint64 = 323 // 灯具
	CategoryMiteRemover     uint64 = 324 // 除螨仪
)

func IsValidBrand(brandID uint64) bool {
	switch brandID {
	case BrandXiaomi, BrandRedmi, BrandMijia:
		return true
	default:
		return false
	}
}

func IsValidCategory(categoryID uint64) bool {
	// 一级分类 1-7
	if categoryID >= 1 && categoryID <= 7 {
		return true
	}
	// 二级分类 - 手机 101-102
	if categoryID >= 101 && categoryID <= 102 {
		return true
	}
	// 二级分类 - 家电 301-324
	if categoryID >= 301 && categoryID <= 324 {
		return true
	}
	return false
}

func PriceToString(price float64) string {
	return strconv.FormatFloat(price, 'f', 2, 64)
}

func PriceConvert(price string) (float64, error) {
	cleaned := strings.TrimSpace(price)
	cleaned = strings.ReplaceAll(cleaned, ",", "")
	if cleaned == "" {
		return 0, nil
	}
	val, err := strconv.ParseFloat(cleaned, 64)
	if err != nil {
		return 0, errs.New(errs.ErrParam.Code, "invalid price format: "+fmt.Sprintf("%v", err))
	}
	return val, nil
}

func ValidateJsonFormat(jsonStr string) error {
	cleaned := strings.TrimSpace(jsonStr)
	if cleaned == "" {
		return nil
	}
	var js map[string]interface{}
	err := json.Unmarshal([]byte(cleaned), &js)
	if err != nil {
		return err
	}
	return nil
}
