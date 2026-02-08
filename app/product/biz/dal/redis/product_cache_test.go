package redis

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

var testRedisClient *redis.Client

func TestMain(m *testing.M) {
	// 初始化测试用 Redis 客户端
	// 使用与实际配置相同的 Redis 地址
	testRedisClient = redis.NewClient(&redis.Options{
		Addr:     "piaohost:6379",
		Password: "123456",
		DB:       15, // 使用独立的DB避免影响正式数据
	})

	// 测试 Redis 连接
	ctx := context.Background()
	if err := testRedisClient.Ping(ctx).Err(); err != nil {
		panic("无法连接 Redis: " + err.Error())
	}

	// 设置测试用的 RedisClient
	RedisClient = testRedisClient

	// 运行测试
	code := m.Run()

	// 清理测试数据
	testRedisClient.FlushDB(ctx)
	testRedisClient.Close()

	os.Exit(code)
}

// TestProductDetailCache 测试商品详情缓存
func TestProductDetailCache(t *testing.T) {
	ctx := context.Background()

	// 准备测试数据
	productID := uint64(10001)
	testProduct := &CachedProductDetail{
		SPU: &CachedSPU{
			ID:            productID,
			Name:          "测试商品",
			CategoryID:    1,
			BrandID:       1,
			PublishStatus: 1,
			MainImage:     "http://example.com/image.jpg",
			SaleCount:     100,
		},
		SKUs: []*CachedSKU{
			{
				ID:      20001,
				SpuID:   productID,
				Name:    "测试SKU",
				Price:   99.99,
				Stock:   50,
				SkuCode: "SKU001",
			},
		},
		Detail: &CachedDetail{
			Description: "这是测试商品描述",
		},
	}

	// 测试1：设置缓存
	err := SetProductDetailCache(ctx, productID, testProduct)
	if err != nil {
		t.Fatalf("设置缓存失败: %v", err)
	}
	t.Log("✓ 缓存设置成功")

	// 测试2：获取缓存
	cachedProduct, err := GetProductDetailFromCache(ctx, productID)
	if err != nil {
		t.Fatalf("获取缓存失败: %v", err)
	}
	if cachedProduct == nil {
		t.Fatal("缓存数据为空")
	}
	if cachedProduct.SPU.Name != testProduct.SPU.Name {
		t.Errorf("商品名称不匹配: got=%s, want=%s", cachedProduct.SPU.Name, testProduct.SPU.Name)
	}
	if len(cachedProduct.SKUs) != len(testProduct.SKUs) {
		t.Errorf("SKU数量不匹配: got=%d, want=%d", len(cachedProduct.SKUs), len(testProduct.SKUs))
	}
	t.Log("✓ 缓存获取成功")

	// 测试3：再次获取（验证缓存一致性）
	cachedProduct2, err := GetProductDetailFromCache(ctx, productID)
	if err != nil {
		t.Fatalf("第二次获取缓存失败: %v", err)
	}
	if cachedProduct2.SPU.ID != cachedProduct.SPU.ID {
		t.Error("两次获取的缓存数据不一致")
	}
	t.Log("✓ 缓存一致性验证通过")

	// 测试4：删除缓存
	err = DeleteProductDetailCache(ctx, productID)
	if err != nil {
		t.Fatalf("删除缓存失败: %v", err)
	}
	t.Log("✓ 缓存删除成功")

	// 测试5：验证删除后获取返回nil
	cachedProduct3, err := GetProductDetailFromCache(ctx, productID)
	if err != nil {
		t.Fatalf("删除后获取缓存失败: %v", err)
	}
	if cachedProduct3 != nil {
		t.Error("删除后缓存应该为nil")
	}
	t.Log("✓ 缓存删除验证通过")
}

// TestHotProductsZSet 测试热门商品 ZSet 排名
func TestHotProductsZSet(t *testing.T) {
	ctx := context.Background()

	// 清理之前的测试数据
	RedisClient.Del(ctx, HotProductsZSetKey)

	// 准备测试数据
	products := []struct {
		id    uint64
		score int
	}{
		{1001, 100},
		{1002, 200},
		{1003, 50},
		{1004, 300},
		{1005, 150},
	}

	// 测试1：添加热门商品分数
	for _, p := range products {
		err := UpdateHotProductScore(ctx, p.id, p.score)
		if err != nil {
			t.Fatalf("更新热门商品分数失败: %v", err)
		}
	}
	t.Log("✓ 热门商品分数设置成功")

	// 测试2：获取 Top 3 热门商品
	topIDs, err := GetTopHotProductIDs(ctx, 3)
	if err != nil {
		t.Fatalf("获取热门商品ID失败: %v", err)
	}
	if len(topIDs) != 3 {
		t.Errorf("期望获取3个商品，实际获取%d个", len(topIDs))
	}
	// 验证排序（分数从高到低）
	// 期望顺序: 1004(300) > 1002(200) > 1005(150)
	expectedOrder := []uint64{1004, 1002, 1005}
	for i, expected := range expectedOrder {
		if topIDs[i] != expected {
			t.Errorf("第%d位商品不匹配: got=%d, want=%d", i+1, topIDs[i], expected)
		}
	}
	t.Log("✓ 热门商品排序正确")

	// 测试3：获取所有热门商品（Top 10）
	allIDs, err := GetTopHotProductIDs(ctx, 10)
	if err != nil {
		t.Fatalf("获取所有热门商品失败: %v", err)
	}
	if len(allIDs) != 5 {
		t.Errorf("期望获取5个商品，实际获取%d个", len(allIDs))
	}
	t.Log("✓ 获取全部热门商品成功")

	// 测试4：增加销量分数
	err = IncrementProductSaleCount(ctx, 1001, 250) // 1001: 100 + 250 = 350
	if err != nil {
		t.Fatalf("增加销量分数失败: %v", err)
	}
	topIDs2, _ := GetTopHotProductIDs(ctx, 1)
	if topIDs2[0] != 1001 {
		t.Errorf("增加分数后，1001应该排第一，实际是%d", topIDs2[0])
	}
	t.Log("✓ 销量增量更新成功")
}

// TestHotProductsListCache 测试热门商品列表缓存
func TestHotProductsListCache(t *testing.T) {
	ctx := context.Background()

	// 准备测试数据
	testProducts := []*HotProductInfo{
		{
			ID:        1001,
			Name:      "热门商品1",
			LowPrice:  19.99,
			MainImage: "http://example.com/1.jpg",
			SaleCount: 500,
		},
		{
			ID:        1002,
			Name:      "热门商品2",
			LowPrice:  29.99,
			MainImage: "http://example.com/2.jpg",
			SaleCount: 400,
		},
		{
			ID:        1003,
			Name:      "热门商品3",
			LowPrice:  39.99,
			MainImage: "http://example.com/3.jpg",
			SaleCount: 300,
		},
	}

	// 测试1：设置热门商品缓存
	err := SetHotProductsCache(ctx, testProducts)
	if err != nil {
		t.Fatalf("设置热门商品缓存失败: %v", err)
	}
	t.Log("✓ 热门商品缓存设置成功")

	// 测试2：获取热门商品缓存
	cachedProducts, err := GetHotProductsFromCache(ctx)
	if err != nil {
		t.Fatalf("获取热门商品缓存失败: %v", err)
	}
	if cachedProducts == nil {
		t.Fatal("缓存数据为空")
	}
	if len(cachedProducts) != len(testProducts) {
		t.Errorf("商品数量不匹配: got=%d, want=%d", len(cachedProducts), len(testProducts))
	}
	t.Log("✓ 热门商品缓存获取成功")

	// 验证数据完整性
	for i, p := range cachedProducts {
		if p.Name != testProducts[i].Name {
			t.Errorf("商品%d名称不匹配", i)
		}
		if p.SaleCount != testProducts[i].SaleCount {
			t.Errorf("商品%d销量不匹配", i)
		}
	}
	t.Log("✓ 热门商品数据完整性验证通过")

	// 测试3：删除热门商品缓存
	err = DeleteHotProductsCache(ctx)
	if err != nil {
		t.Fatalf("删除热门商品缓存失败: %v", err)
	}
	cachedProducts2, _ := GetHotProductsFromCache(ctx)
	if cachedProducts2 != nil {
		t.Error("删除后缓存应该为nil")
	}
	t.Log("✓ 热门商品缓存删除成功")
}

// TestProductListCache 测试商品列表缓存
func TestProductListCache(t *testing.T) {
	ctx := context.Background()

	// 准备测试参数
	page, pageSize := 1, 10
	keyword := ""
	categoryID, brandID := uint64(1), uint64(0)

	testList := &CachedProductList{
		Products: []*CachedSPU{
			{
				ID:        2001,
				Name:      "列表商品1",
				SaleCount: 100,
			},
			{
				ID:        2002,
				Name:      "列表商品2",
				SaleCount: 200,
			},
		},
		Total: 2,
	}

	// 测试1：设置商品列表缓存
	err := SetProductListCache(ctx, page, pageSize, keyword, categoryID, brandID, testList)
	if err != nil {
		t.Fatalf("设置商品列表缓存失败: %v", err)
	}
	t.Log("✓ 商品列表缓存设置成功")

	// 测试2：获取商品列表缓存
	cachedList, err := GetProductListFromCache(ctx, page, pageSize, keyword, categoryID, brandID)
	if err != nil {
		t.Fatalf("获取商品列表缓存失败: %v", err)
	}
	if cachedList == nil {
		t.Fatal("商品列表缓存为空")
	}
	if len(cachedList.Products) != len(testList.Products) {
		t.Errorf("商品数量不匹配: got=%d, want=%d", len(cachedList.Products), len(testList.Products))
	}
	t.Log("✓ 商品列表缓存获取成功")

	// 测试3：不同参数应该获取不到缓存
	cachedList2, err := GetProductListFromCache(ctx, page, pageSize, keyword, uint64(2), brandID) // 不同分类
	if err != nil {
		t.Fatalf("获取不同参数的缓存失败: %v", err)
	}
	if cachedList2 != nil {
		t.Error("不同参数应该获取不到缓存")
	}
	t.Log("✓ 缓存参数隔离验证通过")
}

// TestCacheConcurrency 测试并发读写
func TestCacheConcurrency(t *testing.T) {
	ctx := context.Background()
	productID := uint64(30001)

	testProduct := &CachedProductDetail{
		SPU: &CachedSPU{
			ID:   productID,
			Name: "并发测试商品",
		},
	}

	// 先设置缓存
	SetProductDetailCache(ctx, productID, testProduct)

	// 并发读取
	done := make(chan bool, 100)
	for i := 0; i < 100; i++ {
		go func() {
			_, err := GetProductDetailFromCache(ctx, productID)
			if err != nil {
				t.Errorf("并发读取失败: %v", err)
			}
			done <- true
		}()
	}

	// 等待所有并发完成
	for i := 0; i < 100; i++ {
		<-done
	}
	t.Log("✓ 100次并发读取全部成功")

	// 清理
	DeleteProductDetailCache(ctx, productID)
}

// TestBatchDeleteCache 测试批量删除缓存
func TestBatchDeleteCache(t *testing.T) {
	ctx := context.Background()

	// 准备多个商品缓存
	productIDs := []uint64{40001, 40002, 40003, 40004, 40005}
	for _, id := range productIDs {
		testProduct := &CachedProductDetail{
			SPU: &CachedSPU{
				ID:   id,
				Name: "批量测试商品",
			},
		}
		SetProductDetailCache(ctx, id, testProduct)
	}
	t.Log("✓ 批量设置5个商品缓存成功")

	// 批量删除
	err := BatchDeleteProductDetailCache(ctx, productIDs)
	if err != nil {
		t.Fatalf("批量删除失败: %v", err)
	}
	t.Log("✓ 批量删除成功")

	// 验证全部删除
	for _, id := range productIDs {
		cached, _ := GetProductDetailFromCache(ctx, id)
		if cached != nil {
			t.Errorf("商品%d缓存未被删除", id)
		}
	}
	t.Log("✓ 批量删除验证通过")
}

// TestCacheTTL 测试缓存过期时间
func TestCacheTTL(t *testing.T) {
	ctx := context.Background()
	productID := uint64(50001)

	testProduct := &CachedProductDetail{
		SPU: &CachedSPU{
			ID:   productID,
			Name: "TTL测试商品",
		},
	}

	SetProductDetailCache(ctx, productID, testProduct)

	// 获取 TTL
	key := ProductDetailKeyPrefix + "50001"
	ttl := RedisClient.TTL(ctx, key).Val()

	// TTL 应该在 25-35 分钟之间
	if ttl < 25*time.Minute || ttl > 35*time.Minute {
		t.Errorf("TTL 不在预期范围内: %v", ttl)
	}
	t.Logf("✓ 缓存TTL正确: %v", ttl)

	// 清理
	DeleteProductDetailCache(ctx, productID)
}

// BenchmarkGetProductDetailCache 性能测试
func BenchmarkGetProductDetailCache(b *testing.B) {
	ctx := context.Background()
	productID := uint64(60001)

	testProduct := &CachedProductDetail{
		SPU: &CachedSPU{
			ID:   productID,
			Name: "性能测试商品",
		},
		SKUs: []*CachedSKU{
			{ID: 1, Name: "SKU1"},
			{ID: 2, Name: "SKU2"},
		},
	}
	SetProductDetailCache(ctx, productID, testProduct)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetProductDetailFromCache(ctx, productID)
	}

	// 清理
	DeleteProductDetailCache(ctx, productID)
}

// BenchmarkSetProductDetailCache 性能测试
func BenchmarkSetProductDetailCache(b *testing.B) {
	ctx := context.Background()

	testProduct := &CachedProductDetail{
		SPU: &CachedSPU{
			ID:   70001,
			Name: "性能测试商品",
		},
		SKUs: []*CachedSKU{
			{ID: 1, Name: "SKU1"},
			{ID: 2, Name: "SKU2"},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SetProductDetailCache(ctx, uint64(70001+i), testProduct)
	}
}
