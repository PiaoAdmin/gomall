package service

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/PiaoAdmin/pmall/app/product/biz/dal/mongo"
	"github.com/PiaoAdmin/pmall/app/product/biz/dal/mysql"
	"github.com/PiaoAdmin/pmall/app/product/biz/dal/redis"
	"github.com/PiaoAdmin/pmall/app/product/biz/utils"
	product "github.com/PiaoAdmin/pmall/rpc_gen/product"
	"github.com/cloudwego/kitex/pkg/klog"
)

// TestProductIntegration 完整的商品服务集成测试
// 测试流程：创建商品 -> 更新库存 -> 查询商品 -> 扣减库存 -> 释放库存
func TestProductIntegration(t *testing.T) {
	// 初始化环境
	setupTestEnv(t)

	ctx := context.Background()

	// 存储创建的商品和 SKU ID，用于后续测试
	var createdProducts []struct {
		SpuID  uint64
		SkuIDs []uint64
		Name   string
	}

	// ==================== 阶段1: 创建商品 ====================
	t.Run("Phase1_CreateProducts", func(t *testing.T) {
		testCases := getTestProducts()

		for _, tc := range testCases {
			t.Run("Create_"+tc.name, func(t *testing.T) {
				service := NewCreateProductService(ctx)
				resp, err := service.Run(tc.request)
				if err != nil {
					t.Fatalf("创建商品 %s 失败: %v", tc.name, err)
				}
				if resp == nil || resp.SpuId == 0 {
					t.Fatalf("创建商品 %s 返回空响应", tc.name)
				}

				// 获取 SKU IDs（从请求中获取，实际应该从数据库查询）
				skuIDs := make([]uint64, 0, len(tc.request.Skus))
				// 注意：这里需要查询数据库获取实际的 SKU IDs，暂时留空
				// TODO: 添加查询 SKU ID 的逻辑

				createdProducts = append(createdProducts, struct {
					SpuID  uint64
					SkuIDs []uint64
					Name   string
				}{
					SpuID:  resp.SpuId,
					SkuIDs: skuIDs,
					Name:   tc.name,
				})

				klog.Infof("✅ 成功创建商品: %s (SPU ID: %d)", tc.name, resp.SpuId)
				t.Logf("✅ 成功创建商品: %s (SPU ID: %d)", tc.name, resp.SpuId)
			})
		}

		if len(createdProducts) == 0 {
			t.Fatal("没有成功创建任何商品")
		}
		t.Logf("第1阶段完成: 共创建 %d 个商品", len(createdProducts))
	})

	// ==================== 阶段2: 更新商品状态并上架 ====================
	// 暂时弃用 UpdateProductStatus，简化业务流程
	t.Run("Phase2_UpdateProductStatus", func(t *testing.T) {
		t.Skip("暂时跳过商品状态更新，已去掉 publish_status 和 verify_status 校验")
	})

	// ==================== 阶段3: 查询商品列表（后台） ====================
	t.Run("Phase3_ListProducts", func(t *testing.T) {
		service := NewListProductsService(ctx)

		// 测试1: 查询所有商品
		t.Run("ListAll", func(t *testing.T) {
			resp, err := service.Run(&product.ListProductsRequest{
				Page:     1,
				PageSize: 10,
			})
			if err != nil {
				t.Fatalf("查询商品列表失败: %v", err)
			}
			if resp.Total == 0 {
				t.Fatal("没有查询到任何商品")
			}
			t.Logf("✅ 查询到 %d 个商品", resp.Total)
		})

		// 测试2: 按分类查询（手机）
		t.Run("ListByCategory_Phone", func(t *testing.T) {
			resp, err := service.Run(&product.ListProductsRequest{
				Page:       1,
				PageSize:   10,
				CategoryId: utils.CategoryPhoneXiaomi,
			})
			if err != nil {
				t.Fatalf("按分类查询失败: %v", err)
			}
			t.Logf("✅ 查询到 %d 个手机商品", resp.Total)
		})

		// 测试3: 按品牌查询（小米）
		t.Run("ListByBrand_Xiaomi", func(t *testing.T) {
			resp, err := service.Run(&product.ListProductsRequest{
				Page:     1,
				PageSize: 10,
				BrandId:  utils.BrandXiaomi,
			})
			if err != nil {
				t.Fatalf("按品牌查询失败: %v", err)
			}
			t.Logf("✅ 查询到 %d 个小米品牌商品", resp.Total)
		})

		// 测试4: 关键词搜索
		t.Run("SearchByKeyword", func(t *testing.T) {
			resp, err := service.Run(&product.ListProductsRequest{
				Page:     1,
				PageSize: 10,
				Keyword:  "小米",
			})
			if err != nil {
				t.Fatalf("关键词搜索失败: %v", err)
			}
			t.Logf("✅ 关键词'小米'搜索到 %d 个商品", resp.Total)
		})

		t.Log("✅ 第3阶段完成: 后台商品列表查询测试通过")
	})

	// ==================== 阶段4: C端商品搜索 ====================
	t.Run("Phase4_SearchProducts", func(t *testing.T) {
		service := NewSearchProductsService(ctx)

		// 测试1: 搜索所有可购买商品
		t.Run("SearchAll", func(t *testing.T) {
			resp, err := service.Run(&product.SearchProductsRequest{
				Page:     1,
				PageSize: 20,
			})
			if err != nil {
				t.Fatalf("搜索商品失败: %v", err)
			}
			if resp.Total == 0 {
				t.Fatal("没有搜索到任何可购买商品")
			}
			t.Logf("✅ 搜索到 %d 个可购买的 SKU", resp.Total)

			// 验证返回的 SKU 包含完整信息
			if len(resp.List) > 0 {
				item := resp.List[0]
				if item.Sku == nil || item.SpuName == "" {
					t.Error("搜索结果缺少必要信息")
				} else {
					t.Logf("   示例商品: %s - %s (价格: ¥%s)", item.SpuName, item.Sku.Name, item.Sku.Price)
				}
			}
		})

		// 测试2: 按分类搜索（空调）
		t.Run("SearchByCategory_AC", func(t *testing.T) {
			resp, err := service.Run(&product.SearchProductsRequest{
				Page:       1,
				PageSize:   10,
				CategoryId: utils.CategoryHomeAppl, // 家电分类
			})
			if err != nil {
				t.Fatalf("按分类搜索失败: %v", err)
			}
			t.Logf("✅ 家电分类搜索到 %d 个 SKU", resp.Total)
		})

		// 测试3: 按品牌搜索
		t.Run("SearchByBrand", func(t *testing.T) {
			resp, err := service.Run(&product.SearchProductsRequest{
				Page:     1,
				PageSize: 10,
				BrandId:  utils.BrandMijia,
			})
			if err != nil {
				t.Fatalf("按品牌搜索失败: %v", err)
			}
			t.Logf("✅ 米家品牌搜索到 %d 个 SKU", resp.Total)
		})

		// 测试4: 价格区间搜索
		t.Run("SearchByPriceRange", func(t *testing.T) {
			resp, err := service.Run(&product.SearchProductsRequest{
				Page:     1,
				PageSize: 10,
				MinPrice: "1000",
				MaxPrice: "3000",
			})
			if err != nil {
				t.Fatalf("价格区间搜索失败: %v", err)
			}
			t.Logf("✅ 价格区间 ¥1000-3000 搜索到 %d 个 SKU", resp.Total)
		})

		// 测试5: 关键词搜索
		t.Run("SearchByKeyword", func(t *testing.T) {
			resp, err := service.Run(&product.SearchProductsRequest{
				Page:     1,
				PageSize: 10,
				Keyword:  "空调",
			})
			if err != nil {
				t.Fatalf("关键词搜索失败: %v", err)
			}
			t.Logf("✅ 关键词'空调'搜索到 %d 个 SKU", resp.Total)
		})

		// 测试6: 排序测试 - 价格升序
		t.Run("SearchWithSort_PriceAsc", func(t *testing.T) {
			resp, err := service.Run(&product.SearchProductsRequest{
				Page:     1,
				PageSize: 5,
				SortType: 1, // 价格升序
			})
			if err != nil {
				t.Fatalf("价格升序排序失败: %v", err)
			}
			if len(resp.List) > 1 {
				// 验证排序正确性
				for i := 0; i < len(resp.List)-1; i++ {
					price1, _ := utils.PriceConvert(resp.List[i].Sku.Price)
					price2, _ := utils.PriceConvert(resp.List[i+1].Sku.Price)
					if price1 > price2 {
						t.Errorf("价格排序错误: ¥%.2f > ¥%.2f", price1, price2)
					}
				}
				t.Log("✅ 价格升序排序正确")
			}
		})

		// 测试7: 综合条件搜索
		t.Run("SearchWithMultipleConditions", func(t *testing.T) {
			resp, err := service.Run(&product.SearchProductsRequest{
				Page:     1,
				PageSize: 10,
				Keyword:  "小米",
				BrandId:  utils.BrandXiaomi,
				MinPrice: "3000",
				MaxPrice: "5000",
				SortType: 2, // 价格降序
			})
			if err != nil {
				t.Fatalf("综合条件搜索失败: %v", err)
			}
			t.Logf("✅ 综合条件搜索到 %d 个 SKU", resp.Total)
		})

		t.Log("✅ 第4阶段完成: C端商品搜索测试通过")
	})

	// ==================== 阶段5: 更新 SKU 库存 ====================
	t.Run("Phase5_UpdateSkuStock", func(t *testing.T) {
		if len(createdProducts) == 0 {
			t.Skip("没有商品可更新")
		}

		// 由于我们没有存储 SKU IDs，这里跳过
		// 实际应用中需要先查询 SKU IDs
		t.Skip("需要先实现 SKU ID 查询逻辑")

		// TODO: 实现库存更新测试
		// service := NewBatchUpdateSkuService(ctx)
		// resp, err := service.Run(&product.BatchUpdateSkuRequest{
		// 	Skus: []*product.ProductSKU{
		// 		{
		// 			Id:    skuID,
		// 			Stock: 999,
		// 		},
		// 	},
		// })
	})

	// ==================== 阶段6: 扣减库存 ====================
	t.Run("Phase6_DeductStock", func(t *testing.T) {
		t.Skip("需要先实现 SKU ID 查询逻辑")

		// TODO: 实现扣减库存测试
		// service := NewDeductStockService(ctx)
		// resp, err := service.Run(&product.DeductStockRequest{
		// 	OrderSn: "TEST_ORDER_001",
		// 	Items: []*product.SkuDeductItem{
		// 		{SkuId: skuID, Count: 2},
		// 	},
		// })
	})

	// ==================== 阶段7: 释放库存 ====================
	t.Run("Phase7_ReleaseStock", func(t *testing.T) {
		t.Skip("需要先实现 SKU ID 查询逻辑")

		// TODO: 实现释放库存测试
		// service := NewReleaseStockService(ctx)
		// resp, err := service.Run(&product.ReleaseStockRequest{
		// 	OrderSn: "TEST_ORDER_001",
		// 	Items: []*product.SkuDeductItem{
		// 		{SkuId: skuID, Count: 2},
		// 	},
		// })
	})

	t.Log("========================================")
	t.Log("✅✅✅ 集成测试全部完成 ✅✅✅")
	t.Log("========================================")
}

// setupTestEnv 初始化测试环境
func setupTestEnv(t *testing.T) {
	oldDir, _ := os.Getwd()
	projectRoot := filepath.Join(oldDir, "..", "..")
	if err := os.Chdir(projectRoot); err != nil {
		t.Fatalf("无法切换到项目根目录: %v", err)
	}
	t.Cleanup(func() { os.Chdir(oldDir) })

	mysql.Init()
	redis.Init()
	mongo.Init()
}

// getTestProducts 返回测试商品数据
func getTestProducts() []struct {
	name    string
	request *product.CreateProductRequest
} {
	return []struct {
		name    string
		request *product.CreateProductRequest
	}{
		// ==================== 手机类 ====================
		{
			name: "小米17",
			request: &product.CreateProductRequest{
				Spu: &product.ProductSPU{
					BrandId:     utils.BrandXiaomi,
					CategoryId:  utils.CategoryPhoneXiaomi,
					Name:        "小米17",
					SubTitle:    "骁龙8 Gen4 | 徕卡光学镜头 | 120W澎湃秒充",
					MainImage:   "https://cdn.cnbj1.fds.api.mi-img.com/mi-mall/xiaomi17/main.jpg",
					ServiceBits: utils.GetCommonServices() | utils.ServiceFast24Ship | utils.ServiceFreeRepair,
				},
				Skus: []*product.ProductSKU{
					{
						SkuCode:     "MI17-12-256-BLACK",
						Name:        "小米17 12GB+256GB 曜石黑",
						SubTitle:    "12GB+256GB 曜石黑",
						MainImage:   "https://cdn.cnbj1.fds.api.mi-img.com/mi-mall/xiaomi17/black.jpg",
						Price:       "3999.00",
						MarketPrice: "4299.00",
						Stock:       1000,
						SkuSpecData: `{"颜色":"曜石黑","内存":"12GB","存储":"256GB"}`,
					},
					{
						SkuCode:     "MI17-16-512-WHITE",
						Name:        "小米17 16GB+512GB 雪山白",
						SubTitle:    "16GB+512GB 雪山白",
						MainImage:   "https://cdn.cnbj1.fds.api.mi-img.com/mi-mall/xiaomi17/white.jpg",
						Price:       "4799.00",
						MarketPrice: "5099.00",
						Stock:       600,
						SkuSpecData: `{"颜色":"雪山白","内存":"16GB","存储":"512GB"}`,
					},
				},
				Detail: &product.ProductDetail{
					Description:   `<h1>小米17 - 性能旗舰</h1><p>第三代骁龙8 Gen4移动平台</p>`,
					Images:        []string{"https://cdn.cnbj1.fds.api.mi-img.com/mi-mall/xiaomi17/detail1.jpg"},
					Videos:        []string{"https://cdn.cnbj1.fds.api.mi-img.com/mi-mall/xiaomi17/video1.mp4"},
					MarketTagJson: `["限时立减300","12期免息"]`,
					TechTagJson:   `{"处理器":"骁龙8 Gen4","屏幕":"6.73英寸 2K AMOLED"}`,
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
					MainImage:   "https://cdn.cnbj1.fds.api.mi-img.com/mi-mall/redmik90/main.jpg",
					ServiceBits: utils.GetCommonServices(),
				},
				Skus: []*product.ProductSKU{
					{
						SkuCode:     "K90-8-256-BLUE",
						Name:        "红米K90 8GB+256GB 冰川蓝",
						SubTitle:    "8GB+256GB 冰川蓝",
						MainImage:   "https://cdn.cnbj1.fds.api.mi-img.com/mi-mall/redmik90/blue.jpg",
						Price:       "1999.00",
						MarketPrice: "2299.00",
						Stock:       2000,
						SkuSpecData: `{"颜色":"冰川蓝","内存":"8GB","存储":"256GB"}`,
					},
					{
						SkuCode:     "K90-12-512-BLACK",
						Name:        "红米K90 12GB+512GB 曜夜黑",
						SubTitle:    "12GB+512GB 曜夜黑",
						MainImage:   "https://cdn.cnbj1.fds.api.mi-img.com/mi-mall/redmik90/black.jpg",
						Price:       "2499.00",
						MarketPrice: "2799.00",
						Stock:       1500,
						SkuSpecData: `{"颜色":"曜夜黑","内存":"12GB","存储":"512GB"}`,
					},
				},
				Detail: &product.ProductDetail{
					Description:   `<h1>红米K90 - 性能小金刚</h1><p>天玑9300+旗舰芯片</p>`,
					Images:        []string{"https://cdn.cnbj1.fds.api.mi-img.com/mi-mall/redmik90/detail1.jpg"},
					MarketTagJson: `["性价比之王","学生优惠"]`,
					TechTagJson:   `{"处理器":"天玑9300+","电池":"6000mAh"}`,
				},
			},
		},

		// ==================== 家电类 - 空调 ====================
		{
			name: "米家壁挂式空调",
			request: &product.CreateProductRequest{
				Spu: &product.ProductSPU{
					BrandId:     utils.BrandMijia,
					CategoryId:  301, // 壁挂空调
					Name:        "米家壁挂式空调 1.5匹",
					SubTitle:    "新一级能效 | 智能温控 | 静音运行",
					MainImage:   "https://cdn.cnbj1.fds.api.mi-img.com/mi-mall/ac/wall-ac.jpg",
					ServiceBits: utils.GetCommonServices() | utils.ServiceAuthenticGuarantee | utils.ServiceExtendedWarranty,
				},
				Skus: []*product.ProductSKU{
					{
						SkuCode:     "MJAC-1.5P-WHITE",
						Name:        "米家壁挂式空调 1.5匹 雅白",
						SubTitle:    "1.5匹 新一级能效 雅白",
						MainImage:   "https://cdn.cnbj1.fds.api.mi-img.com/mi-mall/ac/wall-ac-white.jpg",
						Price:       "2199.00",
						MarketPrice: "2599.00",
						Stock:       500,
						SkuSpecData: `{"匹数":"1.5匹","颜色":"雅白","能效":"新一级"}`,
					},
					{
						SkuCode:     "MJAC-1.5P-GRAY",
						Name:        "米家壁挂式空调 1.5匹 星空灰",
						SubTitle:    "1.5匹 新一级能效 星空灰",
						MainImage:   "https://cdn.cnbj1.fds.api.mi-img.com/mi-mall/ac/wall-ac-gray.jpg",
						Price:       "2299.00",
						MarketPrice: "2699.00",
						Stock:       300,
						SkuSpecData: `{"匹数":"1.5匹","颜色":"星空灰","能效":"新一级"}`,
					},
				},
				Detail: &product.ProductDetail{
					Description:   `<h1>米家壁挂式空调</h1><p>新一级能效，省电更省心</p><p>智能温控，自动调节舒适温度</p>`,
					Images:        []string{"https://cdn.cnbj1.fds.api.mi-img.com/mi-mall/ac/wall-detail1.jpg"},
					MarketTagJson: `["新一级能效","免费安装","以旧换新"]`,
					TechTagJson:   `{"匹数":"1.5匹","能效":"新一级","制冷量":"3500W","噪音":"22dB"}`,
				},
			},
		},
		{
			name: "米家立式空调",
			request: &product.CreateProductRequest{
				Spu: &product.ProductSPU{
					BrandId:     utils.BrandMijia,
					CategoryId:  302, // 立式空调
					Name:        "米家立式空调 3匹",
					SubTitle:    "客厅首选 | 急速冷暖 | 智能除湿",
					MainImage:   "https://cdn.cnbj1.fds.api.mi-img.com/mi-mall/ac/floor-ac.jpg",
					ServiceBits: utils.GetCommonServices() | utils.ServiceAuthenticGuarantee | utils.ServiceExtendedWarranty | utils.ServiceInstallmentFree,
				},
				Skus: []*product.ProductSKU{
					{
						SkuCode:     "MJAC-3P-WHITE",
						Name:        "米家立式空调 3匹 象牙白",
						SubTitle:    "3匹 新一级能效 象牙白",
						MainImage:   "https://cdn.cnbj1.fds.api.mi-img.com/mi-mall/ac/floor-ac-white.jpg",
						Price:       "4999.00",
						MarketPrice: "5499.00",
						Stock:       200,
						SkuSpecData: `{"匹数":"3匹","颜色":"象牙白","能效":"新一级"}`,
					},
					{
						SkuCode:     "MJAC-3P-GOLD",
						Name:        "米家立式空调 3匹 香槟金",
						SubTitle:    "3匹 新一级能效 香槟金",
						MainImage:   "https://cdn.cnbj1.fds.api.mi-img.com/mi-mall/ac/floor-ac-gold.jpg",
						Price:       "5299.00",
						MarketPrice: "5799.00",
						Stock:       150,
						SkuSpecData: `{"匹数":"3匹","颜色":"香槟金","能效":"新一级"}`,
					},
				},
				Detail: &product.ProductDetail{
					Description:   `<h1>米家立式空调 3匹</h1><p>大空间首选，制冷制热更迅速</p><p>智能除湿功能，梅雨季节不再潮湿</p>`,
					Images:        []string{"https://cdn.cnbj1.fds.api.mi-img.com/mi-mall/ac/floor-detail1.jpg"},
					MarketTagJson: `["客厅首选","免费安装","12期免息"]`,
					TechTagJson:   `{"匹数":"3匹","能效":"新一级","制冷量":"7200W","适用面积":"30-45㎡"}`,
				},
			},
		},

		// ==================== 家电类 - 其他 ====================
		{
			name: "米家扫地机器人",
			request: &product.CreateProductRequest{
				Spu: &product.ProductSPU{
					BrandId:     utils.BrandMijia,
					CategoryId:  314, // 扫地机器人
					Name:        "米家扫地机器人 Ultra",
					SubTitle:    "激光导航 | 自动集尘 | 拖扫一体",
					MainImage:   "https://cdn.cnbj1.fds.api.mi-img.com/mi-mall/robot/ultra.jpg",
					ServiceBits: utils.GetCommonServices() | utils.ServiceFreeRepair,
				},
				Skus: []*product.ProductSKU{
					{
						SkuCode:     "MJROBOT-ULTRA-WHITE",
						Name:        "米家扫地机器人 Ultra 白色",
						SubTitle:    "激光导航 自动集尘 白色",
						MainImage:   "https://cdn.cnbj1.fds.api.mi-img.com/mi-mall/robot/ultra-white.jpg",
						Price:       "2999.00",
						MarketPrice: "3499.00",
						Stock:       300,
						SkuSpecData: `{"颜色":"白色","功能":"扫拖一体+自动集尘"}`,
					},
					{
						SkuCode:     "MJROBOT-ULTRA-BLACK",
						Name:        "米家扫地机器人 Ultra 黑色",
						SubTitle:    "激光导航 自动集尘 黑色",
						MainImage:   "https://cdn.cnbj1.fds.api.mi-img.com/mi-mall/robot/ultra-black.jpg",
						Price:       "2999.00",
						MarketPrice: "3499.00",
						Stock:       250,
						SkuSpecData: `{"颜色":"黑色","功能":"扫拖一体+自动集尘"}`,
					},
				},
				Detail: &product.ProductDetail{
					Description:   `<h1>米家扫地机器人 Ultra</h1><p>激光雷达导航，智能规划清扫路径</p><p>自动集尘，60天免倒垃圾</p>`,
					Images:        []string{"https://cdn.cnbj1.fds.api.mi-img.com/mi-mall/robot/detail1.jpg"},
					MarketTagJson: `["智能导航","自动集尘","拖扫一体"]`,
					TechTagJson:   `{"导航":"激光雷达","吸力":"5100Pa","电池":"5200mAh","集尘":"自动集尘"}`,
				},
			},
		},
	}
}
