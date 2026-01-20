package service

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/PiaoAdmin/pmall/app/product/biz/dal/mongo"
	"github.com/PiaoAdmin/pmall/app/product/biz/dal/mysql"
	"github.com/PiaoAdmin/pmall/app/product/biz/utils"
	product "github.com/PiaoAdmin/pmall/rpc_gen/product"
	"github.com/cloudwego/kitex/pkg/klog"
)

// TestCreateProduct 测试创建三款手机：小米17、小米17 Ultra、红米K90
func TestCreateProduct(t *testing.T) {
	// 切换到项目根目录，以便正确读取conf配置文件
	oldDir, _ := os.Getwd()
	projectRoot := filepath.Join(oldDir, "..", "..")
	if err := os.Chdir(projectRoot); err != nil {
		t.Fatalf("无法切换到项目根目录: %v", err)
	}
	defer os.Chdir(oldDir)

	// 初始化数据库连接
	mysql.Init()
	mongo.Init()

	ctx := context.Background()

	// 测试数据：小米17、小米17 Ultra、红米K90
	testCases := []struct {
		name    string
		request *product.CreateProductRequest
	}{
		{
			name: "小米17",
			request: &product.CreateProductRequest{
				Spu: &product.ProductSPU{
					BrandId:     utils.BrandXiaomi,
					CategoryId:  utils.CategoryPhoneXiaomi,
					Name:        "小米17",
					SubTitle:    "骁龙8 Gen4 | 徕卡光学镜头 | 120W澎湃秒充",
					MainImage:   "https://cdn.example.com/xiaomi17/main.jpg",
					ServiceBits: utils.GetCommonServices() | utils.ServiceFast24Ship | utils.ServiceFreeRepair,
				},
				Skus: []*product.ProductSKU{
					{
						SkuCode:     "MI17-12-256-BLACK",
						Name:        "小米17 12GB+256GB 曜石黑",
						SubTitle:    "12GB+256GB 曜石黑",
						MainImage:   "https://cdn.example.com/xiaomi17/black.jpg",
						Price:       "3999.00",
						MarketPrice: "4299.00",
						Stock:       1000,
						SkuSpecData: `{"颜色":"曜石黑","内存":"12GB","存储":"256GB"}`,
					},
					{
						SkuCode:     "MI17-12-512-BLACK",
						Name:        "小米17 12GB+512GB 曜石黑",
						SubTitle:    "12GB+512GB 曜石黑",
						MainImage:   "https://cdn.example.com/xiaomi17/black.jpg",
						Price:       "4399.00",
						MarketPrice: "4699.00",
						Stock:       800,
						SkuSpecData: `{"颜色":"曜石黑","内存":"12GB","存储":"512GB"}`,
					},
					{
						SkuCode:     "MI17-16-512-WHITE",
						Name:        "小米17 16GB+512GB 雪山白",
						SubTitle:    "16GB+512GB 雪山白",
						MainImage:   "https://cdn.example.com/xiaomi17/white.jpg",
						Price:       "4799.00",
						MarketPrice: "5099.00",
						Stock:       600,
						SkuSpecData: `{"颜色":"雪山白","内存":"16GB","存储":"512GB"}`,
					},
					{
						SkuCode:     "MI17-16-1TB-WHITE",
						Name:        "小米17 16GB+1TB 雪山白",
						SubTitle:    "16GB+1TB 雪山白",
						MainImage:   "https://cdn.example.com/xiaomi17/white.jpg",
						Price:       "5299.00",
						MarketPrice: "5599.00",
						Stock:       400,
						SkuSpecData: `{"颜色":"雪山白","内存":"16GB","存储":"1TB"}`,
					},
				},
				Detail: &product.ProductDetail{
					Description: `<h1>小米17 - 性能旗舰</h1>
<p>第三代骁龙8 Gen4移动平台，性能飙升35%</p>
<p>徕卡专业光学镜头，50MP主摄+50MP超广角+50MP长焦</p>
<p>120W有线澎湃秒充+50W无线秒充，19分钟充满</p>
<p>6.73英寸AMOLED 2K超清屏，120Hz自适应刷新率</p>`,
					Images: []string{
						"https://cdn.example.com/xiaomi17/detail1.jpg",
						"https://cdn.example.com/xiaomi17/detail2.jpg",
						"https://cdn.example.com/xiaomi17/detail3.jpg",
						"https://cdn.example.com/xiaomi17/detail4.jpg",
					},
					Videos: []string{
						"https://cdn.example.com/xiaomi17/video1.mp4",
					},
					MarketTagJson: `["限时立减300","12期免息","以旧换新最高抵1000"]`,
					TechTagJson:   `{"处理器":"骁龙8 Gen4","屏幕":"6.73英寸 2K AMOLED","电池":"5000mAh","充电":"120W有线+50W无线","相机":"徕卡三摄 50MP+50MP+50MP","网络":"5G双卡双待","系统":"MIUI 16"}`,
					FaqJson:       `[{"q":"支持双卡吗？","a":"支持5G双卡双待"},{"q":"充电多久充满？","a":"120W快充约19分钟充满"}]`,
				},
			},
		},
		{
			name: "小米17 Ultra",
			request: &product.CreateProductRequest{
				Spu: &product.ProductSPU{
					BrandId:     utils.BrandXiaomi,
					CategoryId:  utils.CategoryPhoneXiaomi,
					Name:        "小米17 Ultra",
					SubTitle:    "专业影像旗舰 | 徕卡1英寸大底主摄 | 卫星通信",
					MainImage:   "https://cdn.example.com/xiaomi17ultra/main.jpg",
					ServiceBits: utils.GetPremiumServices(),
				},
				Skus: []*product.ProductSKU{
					{
						SkuCode:     "MI17U-16-512-TITAN",
						Name:        "小米17 Ultra 16GB+512GB 钛金属",
						SubTitle:    "16GB+512GB 钛金属",
						MainImage:   "https://cdn.example.com/xiaomi17ultra/titan.jpg",
						Price:       "6499.00",
						MarketPrice: "6999.00",
						Stock:       500,
						SkuSpecData: `{"颜色":"钛金属","内存":"16GB","存储":"512GB"}`,
					},
					{
						SkuCode:     "MI17U-16-1TB-TITAN",
						Name:        "小米17 Ultra 16GB+1TB 钛金属",
						SubTitle:    "16GB+1TB 钛金属",
						MainImage:   "https://cdn.example.com/xiaomi17ultra/titan.jpg",
						Price:       "6999.00",
						MarketPrice: "7499.00",
						Stock:       300,
						SkuSpecData: `{"颜色":"钛金属","内存":"16GB","存储":"1TB"}`,
					},
					{
						SkuCode:     "MI17U-24-1TB-BLACK",
						Name:        "小米17 Ultra 24GB+1TB 陨石黑",
						SubTitle:    "24GB+1TB 陨石黑",
						MainImage:   "https://cdn.example.com/xiaomi17ultra/black.jpg",
						Price:       "7499.00",
						MarketPrice: "7999.00",
						Stock:       200,
						SkuSpecData: `{"颜色":"陨石黑","内存":"24GB","存储":"1TB"}`,
					},
				},
				Detail: &product.ProductDetail{
					Description: `<h1>小米17 Ultra - 专业影像旗舰</h1>
<p>徕卡专业光学系统，1英寸超大底主摄，感光能力提升76%</p>
<p>可变光圈f/1.42-f/4.0，支持全像素对焦</p>
<p>卫星通信功能，无信号也能发送求救信息</p>
<p>第二代骁龙8 Gen4领先版，性能再进一步</p>
<p>钛合金中框，航天级材质，轻量化设计</p>`,
					Images: []string{
						"https://cdn.example.com/xiaomi17ultra/detail1.jpg",
						"https://cdn.example.com/xiaomi17ultra/detail2.jpg",
						"https://cdn.example.com/xiaomi17ultra/detail3.jpg",
						"https://cdn.example.com/xiaomi17ultra/detail4.jpg",
						"https://cdn.example.com/xiaomi17ultra/detail5.jpg",
					},
					Videos: []string{
						"https://cdn.example.com/xiaomi17ultra/video1.mp4",
						"https://cdn.example.com/xiaomi17ultra/video2.mp4",
					},
					MarketTagJson: `["旗舰影像","卫星通信","钛合金机身","24期免息"]`,
					TechTagJson:   `{"处理器":"骁龙8 Gen4领先版","屏幕":"6.73英寸 2K LTPO AMOLED","电池":"5300mAh","充电":"120W有线+80W无线","相机":"徕卡四摄 1英寸主摄+超广角+潜望长焦+微距","网络":"5G+卫星通信","系统":"MIUI 16 Pro","材质":"钛合金中框"}`,
					FaqJson:       `[{"q":"卫星通信在哪些地区可用？","a":"目前支持中国大陆及周边海域"},{"q":"1英寸大底有什么优势？","a":"感光面积更大，夜景拍摄更清晰，虚化效果更自然"}]`,
				},
			},
		},
		{
			name: "红米K90",
			request: &product.CreateProductRequest{
				Spu: &product.ProductSPU{
					BrandId:     utils.BrandRedmi,
					CategoryId:  utils.CategoryPhoneRedmi,
					Name:        "红米K90",
					SubTitle:    "性能小金刚 | 天玑9300+ | 6000mAh超大电池",
					MainImage:   "https://cdn.example.com/redmik90/main.jpg",
					ServiceBits: utils.GetCommonServices() | utils.ServiceFast24Ship,
				},
				Skus: []*product.ProductSKU{
					{
						SkuCode:     "K90-8-256-BLUE",
						Name:        "红米K90 8GB+256GB 冰川蓝",
						SubTitle:    "8GB+256GB 冰川蓝",
						MainImage:   "https://cdn.example.com/redmik90/blue.jpg",
						Price:       "1999.00",
						MarketPrice: "2299.00",
						Stock:       2000,
						SkuSpecData: `{"颜色":"冰川蓝","内存":"8GB","存储":"256GB"}`,
					},
					{
						SkuCode:     "K90-12-256-BLUE",
						Name:        "红米K90 12GB+256GB 冰川蓝",
						SubTitle:    "12GB+256GB 冰川蓝",
						MainImage:   "https://cdn.example.com/redmik90/blue.jpg",
						Price:       "2199.00",
						MarketPrice: "2499.00",
						Stock:       1800,
						SkuSpecData: `{"颜色":"冰川蓝","内存":"12GB","存储":"256GB"}`,
					},
					{
						SkuCode:     "K90-12-512-BLACK",
						Name:        "红米K90 12GB+512GB 曜夜黑",
						SubTitle:    "12GB+512GB 曜夜黑",
						MainImage:   "https://cdn.example.com/redmik90/black.jpg",
						Price:       "2499.00",
						MarketPrice: "2799.00",
						Stock:       1500,
						SkuSpecData: `{"颜色":"曜夜黑","内存":"12GB","存储":"512GB"}`,
					},
					{
						SkuCode:     "K90-16-512-GREEN",
						Name:        "红米K90 16GB+512GB 青山绿",
						SubTitle:    "16GB+512GB 青山绿",
						MainImage:   "https://cdn.example.com/redmik90/green.jpg",
						Price:       "2799.00",
						MarketPrice: "3099.00",
						Stock:       1200,
						SkuSpecData: `{"颜色":"青山绿","内存":"16GB","存储":"512GB"}`,
					},
				},
				Detail: &product.ProductDetail{
					Description: `<h1>红米K90 - 性能小金刚</h1>
<p>联发科天玑9300+旗舰芯片，性能强悍</p>
<p>6000mAh超大电池+120W快充，续航无忧</p>
<p>144Hz电竞屏，游戏体验拉满</p>
<p>5000万像素光学防抖主摄，记录精彩瞬间</p>`,
					Images: []string{
						"https://cdn.example.com/redmik90/detail1.jpg",
						"https://cdn.example.com/redmik90/detail2.jpg",
						"https://cdn.example.com/redmik90/detail3.jpg",
					},
					Videos: []string{
						"https://cdn.example.com/redmik90/video1.mp4",
					},
					MarketTagJson: `["超大电池","性价比之王","学生优惠","6期免息"]`,
					TechTagJson:   `{"处理器":"天玑9300+","屏幕":"6.67英寸 144Hz LCD","电池":"6000mAh","充电":"120W有线快充","相机":"50MP主摄+8MP超广角+2MP微距","网络":"5G双卡双待","系统":"MIUI 15"}`,
					FaqJson:       `[{"q":"续航能用多久？","a":"6000mAh超大电池，中度使用可达2天"},{"q":"支持红外遥控吗？","a":"支持万能红外遥控，可控制家电"}]`,
				},
			},
		},
	}

	// 依次创建三款手机
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			service := NewCreateProductService(ctx)
			resp, err := service.Run(tc.request)
			if err != nil {
				t.Errorf("创建商品 %s 失败: %v", tc.name, err)
				return
			}
			if resp == nil || resp.SpuId == 0 {
				t.Errorf("创建商品 %s 返回空响应", tc.name)
				return
			}
			klog.Infof("✅ 成功创建商品 %s, SPU ID: %d", tc.name, resp.SpuId)
			t.Logf("✅ 成功创建商品 %s, SPU ID: %d", tc.name, resp.SpuId)
		})
	}
}
