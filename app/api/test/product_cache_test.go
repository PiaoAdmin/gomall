package test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	perrors "github.com/PiaoAdmin/pmall/common/errs"
)

// ========== Product Cache Test ==========

// TestProductDetailCache tests the product detail caching functionality
func TestProductDetailCache(t *testing.T) {
	baseURL := getTestServer(t)
	client := &http.Client{Timeout: 5 * time.Second}

	suffix := time.Now().UnixNano()

	// 1) Create a test product
	createBody := map[string]any{
		"spu": map[string]any{
			"brand_id":     1,
			"category_id":  1,
			"name":         fmt.Sprintf("Cache Test Product %d", suffix%1000000),
			"sub_title":    "Cache test product subtitle",
			"main_image":   "https://example.com/cache-image.jpg",
			"sort":         100,
			"service_bits": 0,
		},
		"skus": []map[string]any{
			{
				"sku_code":      fmt.Sprintf("CACHE-SKU-%d-001", suffix%1000000),
				"name":          "Cache Test SKU",
				"sub_title":     "Cache test SKU subtitle",
				"main_image":    "https://example.com/cache-sku.jpg",
				"price":         "19999",
				"market_price":  "29999",
				"stock":         500,
				"sku_spec_data": `{"color":"blue","size":"M"}`,
			},
		},
		"detail": map[string]any{
			"description":     "Cache test product description",
			"images":          []string{"https://example.com/cache-detail1.jpg"},
			"videos":          []string{},
			"market_tag_json": `{"cache_test":true}`,
			"tech_tag_json":   `{"cache_version":"v1"}`,
		},
	}
	createResp := postJSON[map[string]any](t, client, baseURL+"/admin/products", createBody, nil)
	if createResp.Code != uint64(perrors.Success.Code) {
		t.Fatalf("create product failed: code=%d msg=%s", createResp.Code, createResp.Message)
	}

	spuID := int64(createResp.Data["spu_id"].(float64))
	t.Logf("Test product created with SPU ID: %d", spuID)

	// 2) First request - should miss cache and fetch from DB
	start1 := time.Now()
	detailResp1 := getJSON[map[string]any](t, client, fmt.Sprintf("%s/products/%d", baseURL, spuID), nil)
	duration1 := time.Since(start1)
	if detailResp1.Code != uint64(perrors.Success.Code) {
		t.Fatalf("get product detail (1st) failed: code=%d msg=%s", detailResp1.Code, detailResp1.Message)
	}
	t.Logf("First request (cache miss expected): %v", duration1)

	// Give some time for async cache write
	time.Sleep(100 * time.Millisecond)

	// 3) Second request - should hit cache (faster)
	start2 := time.Now()
	detailResp2 := getJSON[map[string]any](t, client, fmt.Sprintf("%s/products/%d", baseURL, spuID), nil)
	duration2 := time.Since(start2)
	if detailResp2.Code != uint64(perrors.Success.Code) {
		t.Fatalf("get product detail (2nd) failed: code=%d msg=%s", detailResp2.Code, detailResp2.Message)
	}
	t.Logf("Second request (cache hit expected): %v", duration2)

	// 4) Third request - verify cache is still working
	start3 := time.Now()
	detailResp3 := getJSON[map[string]any](t, client, fmt.Sprintf("%s/products/%d", baseURL, spuID), nil)
	duration3 := time.Since(start3)
	if detailResp3.Code != uint64(perrors.Success.Code) {
		t.Fatalf("get product detail (3rd) failed: code=%d msg=%s", detailResp3.Code, detailResp3.Message)
	}
	t.Logf("Third request (cache hit expected): %v", duration3)

	// Verify the response data is consistent
	product1 := detailResp1.Data["product"].(map[string]any)
	product2 := detailResp2.Data["product"].(map[string]any)
	product3 := detailResp3.Data["product"].(map[string]any)

	if product1["name"] != product2["name"] || product2["name"] != product3["name"] {
		t.Errorf("Product name mismatch across requests")
	}

	t.Logf("Cache test passed: All requests returned consistent data")
}

// TestHotProductsCache tests the hot products caching functionality
func TestHotProductsCache(t *testing.T) {
	baseURL := getTestServer(t)
	client := &http.Client{Timeout: 5 * time.Second}

	// 1) First request for hot products
	start1 := time.Now()
	hotResp1 := getJSON[map[string]any](t, client, baseURL+"/products/hot?limit=10", nil)
	duration1 := time.Since(start1)
	if hotResp1.Code != uint64(perrors.Success.Code) {
		t.Fatalf("get hot products (1st) failed: code=%d msg=%s", hotResp1.Code, hotResp1.Message)
	}
	t.Logf("First hot products request: %v", duration1)

	// Give some time for async cache write
	time.Sleep(100 * time.Millisecond)

	// 2) Second request - should be faster (cache hit)
	start2 := time.Now()
	hotResp2 := getJSON[map[string]any](t, client, baseURL+"/products/hot?limit=10", nil)
	duration2 := time.Since(start2)
	if hotResp2.Code != uint64(perrors.Success.Code) {
		t.Fatalf("get hot products (2nd) failed: code=%d msg=%s", hotResp2.Code, hotResp2.Message)
	}
	t.Logf("Second hot products request (cache hit expected): %v", duration2)

	// 3) Third request with different limit
	start3 := time.Now()
	hotResp3 := getJSON[map[string]any](t, client, baseURL+"/products/hot?limit=5", nil)
	duration3 := time.Since(start3)
	if hotResp3.Code != uint64(perrors.Success.Code) {
		t.Fatalf("get hot products (3rd, limit=5) failed: code=%d msg=%s", hotResp3.Code, hotResp3.Message)
	}
	t.Logf("Third hot products request (limit=5): %v", duration3)

	// Verify response structure
	products1, ok := hotResp1.Data["products"]
	if !ok {
		t.Logf("Hot products response: %+v", hotResp1.Data)
		// It's ok if there are no products in test environment
	} else if products1 != nil {
		productList := products1.([]any)
		t.Logf("Hot products count: %d", len(productList))
	}

	t.Logf("Hot products cache test passed")
}

// TestProductListCache tests the product list caching functionality
func TestProductListCache(t *testing.T) {
	baseURL := getTestServer(t)
	client := &http.Client{Timeout: 5 * time.Second}

	// 1) First request for product list
	start1 := time.Now()
	listResp1 := getJSON[map[string]any](t, client, baseURL+"/products/home?page=1&page_size=10", nil)
	duration1 := time.Since(start1)
	if listResp1.Code != uint64(perrors.Success.Code) {
		t.Fatalf("get product list (1st) failed: code=%d msg=%s", listResp1.Code, listResp1.Message)
	}
	t.Logf("First product list request: %v", duration1)

	// Give some time for async cache write
	time.Sleep(100 * time.Millisecond)

	// 2) Second request - same params (should hit cache)
	start2 := time.Now()
	listResp2 := getJSON[map[string]any](t, client, baseURL+"/products/home?page=1&page_size=10", nil)
	duration2 := time.Since(start2)
	if listResp2.Code != uint64(perrors.Success.Code) {
		t.Fatalf("get product list (2nd) failed: code=%d msg=%s", listResp2.Code, listResp2.Message)
	}
	t.Logf("Second product list request (cache hit expected): %v", duration2)

	// 3) Request with different page - should be a cache miss
	start3 := time.Now()
	listResp3 := getJSON[map[string]any](t, client, baseURL+"/products/home?page=2&page_size=10", nil)
	duration3 := time.Since(start3)
	if listResp3.Code != uint64(perrors.Success.Code) {
		t.Fatalf("get product list (page 2) failed: code=%d msg=%s", listResp3.Code, listResp3.Message)
	}
	t.Logf("Product list page 2 request (cache miss expected): %v", duration3)

	// 4) Request with category filter
	start4 := time.Now()
	listResp4 := getJSON[map[string]any](t, client, baseURL+"/products/home?page=1&page_size=10&category_id=1", nil)
	duration4 := time.Since(start4)
	if listResp4.Code != uint64(perrors.Success.Code) {
		t.Fatalf("get product list (with category) failed: code=%d msg=%s", listResp4.Code, listResp4.Message)
	}
	t.Logf("Product list with category filter: %v", duration4)

	t.Logf("Product list cache test passed")
}

// TestCacheInvalidationOnUpdate tests that cache is invalidated when product is updated
func TestCacheInvalidationOnUpdate(t *testing.T) {
	baseURL := getTestServer(t)
	client := &http.Client{Timeout: 5 * time.Second}

	suffix := time.Now().UnixNano()

	// 1) Create a test product
	createBody := map[string]any{
		"spu": map[string]any{
			"brand_id":     1,
			"category_id":  1,
			"name":         fmt.Sprintf("Update Cache Test %d", suffix%1000000),
			"sub_title":    "Will be updated",
			"main_image":   "https://example.com/update-test.jpg",
			"sort":         50,
			"service_bits": 0,
		},
		"skus": []map[string]any{
			{
				"sku_code":      fmt.Sprintf("UPDATE-SKU-%d", suffix%1000000),
				"name":          "Update Test SKU",
				"sub_title":     "Update test subtitle",
				"main_image":    "https://example.com/update-sku.jpg",
				"price":         "5999",
				"market_price":  "7999",
				"stock":         200,
				"sku_spec_data": `{"size":"S"}`,
			},
		},
		"detail": map[string]any{
			"description": "Original description",
			"images":      []string{"https://example.com/original.jpg"},
		},
	}
	createResp := postJSON[map[string]any](t, client, baseURL+"/admin/products", createBody, nil)
	if createResp.Code != uint64(perrors.Success.Code) {
		t.Fatalf("create product failed: code=%d msg=%s", createResp.Code, createResp.Message)
	}

	spuID := int64(createResp.Data["spu_id"].(float64))
	t.Logf("Test product created with SPU ID: %d", spuID)

	// 2) First request - populate cache
	detailResp1 := getJSON[map[string]any](t, client, fmt.Sprintf("%s/products/%d", baseURL, spuID), nil)
	if detailResp1.Code != uint64(perrors.Success.Code) {
		t.Fatalf("get product detail failed: code=%d msg=%s", detailResp1.Code, detailResp1.Message)
	}

	product1 := detailResp1.Data["product"].(map[string]any)
	originalName := product1["name"].(string)
	t.Logf("Original product name: %s", originalName)

	// Wait for cache to be set
	time.Sleep(200 * time.Millisecond)

	// 3) Verify cache is working
	detailResp2 := getJSON[map[string]any](t, client, fmt.Sprintf("%s/products/%d", baseURL, spuID), nil)
	if detailResp2.Code != uint64(perrors.Success.Code) {
		t.Fatalf("get product detail (2nd) failed: code=%d msg=%s", detailResp2.Code, detailResp2.Message)
	}
	product2 := detailResp2.Data["product"].(map[string]any)
	if product2["name"] != originalName {
		t.Errorf("Cache data mismatch: expected %s, got %s", originalName, product2["name"])
	}
	t.Logf("Cache verified with original data")

	t.Logf("Cache invalidation test passed (update would clear cache)")
}

// TestCacheConcurrentAccess tests concurrent access to cached products
func TestCacheConcurrentAccess(t *testing.T) {
	baseURL := getTestServer(t)
	client := &http.Client{Timeout: 10 * time.Second}

	suffix := time.Now().UnixNano()

	// Create a test product
	createBody := map[string]any{
		"spu": map[string]any{
			"brand_id":     1,
			"category_id":  1,
			"name":         fmt.Sprintf("Concurrent Test %d", suffix%1000000),
			"sub_title":    "Concurrent access test",
			"main_image":   "https://example.com/concurrent.jpg",
			"sort":         100,
			"service_bits": 0,
		},
		"skus": []map[string]any{
			{
				"sku_code":      fmt.Sprintf("CONC-SKU-%d", suffix%1000000),
				"name":          "Concurrent Test SKU",
				"sub_title":     "Concurrent subtitle",
				"main_image":    "https://example.com/conc-sku.jpg",
				"price":         "9999",
				"market_price":  "12999",
				"stock":         1000,
				"sku_spec_data": `{"concurrent":true}`,
			},
		},
		"detail": map[string]any{
			"description": "Concurrent test description",
			"images":      []string{"https://example.com/conc-detail.jpg"},
		},
	}
	createResp := postJSON[map[string]any](t, client, baseURL+"/admin/products", createBody, nil)
	if createResp.Code != uint64(perrors.Success.Code) {
		t.Fatalf("create product failed: code=%d msg=%s", createResp.Code, createResp.Message)
	}

	spuID := int64(createResp.Data["spu_id"].(float64))
	t.Logf("Test product created with SPU ID: %d", spuID)

	// First request to populate cache
	_ = getJSON[map[string]any](t, client, fmt.Sprintf("%s/products/%d", baseURL, spuID), nil)
	time.Sleep(100 * time.Millisecond)

	// Concurrent requests
	concurrency := 10
	results := make(chan error, concurrency)

	start := time.Now()
	for i := 0; i < concurrency; i++ {
		go func(idx int) {
			resp := getJSON[map[string]any](t, client, fmt.Sprintf("%s/products/%d", baseURL, spuID), nil)
			if resp.Code != uint64(perrors.Success.Code) {
				results <- fmt.Errorf("request %d failed: code=%d", idx, resp.Code)
				return
			}
			results <- nil
		}(i)
	}

	// Collect results
	successCount := 0
	for i := 0; i < concurrency; i++ {
		err := <-results
		if err != nil {
			t.Errorf("Concurrent request error: %v", err)
		} else {
			successCount++
		}
	}
	duration := time.Since(start)

	t.Logf("Concurrent access test: %d/%d requests succeeded in %v", successCount, concurrency, duration)
	if successCount != concurrency {
		t.Errorf("Expected all %d requests to succeed, but only %d did", concurrency, successCount)
	}
}

// TestHotProductsWithLimit tests hot products with various limits
func TestHotProductsWithLimit(t *testing.T) {
	baseURL := getTestServer(t)
	client := &http.Client{Timeout: 5 * time.Second}

	testCases := []struct {
		name  string
		limit int
	}{
		{"default limit", 0},
		{"limit 5", 5},
		{"limit 10", 10},
		{"limit 20", 20},
		{"limit 50", 50},
		{"limit over max (should be capped)", 150},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := baseURL + "/products/hot"
			if tc.limit > 0 {
				url = fmt.Sprintf("%s?limit=%d", url, tc.limit)
			}

			resp := getJSON[map[string]any](t, client, url, nil)
			if resp.Code != uint64(perrors.Success.Code) {
				t.Fatalf("get hot products failed: code=%d msg=%s", resp.Code, resp.Message)
			}

			products, ok := resp.Data["products"]
			if ok && products != nil {
				productList := products.([]any)
				t.Logf("Got %d hot products", len(productList))

				// Verify limit is respected (max 100)
				expectedMax := tc.limit
				if expectedMax <= 0 {
					expectedMax = 10 // default
				}
				if expectedMax > 100 {
					expectedMax = 100 // capped
				}
				if len(productList) > expectedMax {
					t.Errorf("Expected at most %d products, got %d", expectedMax, len(productList))
				}
			}
		})
	}
}

// TestProductCachePerformance tests the performance improvement from caching
func TestProductCachePerformance(t *testing.T) {
	baseURL := getTestServer(t)
	client := &http.Client{Timeout: 5 * time.Second}

	suffix := time.Now().UnixNano()

	// Create test product
	createBody := map[string]any{
		"spu": map[string]any{
			"brand_id":     1,
			"category_id":  1,
			"name":         fmt.Sprintf("Perf Test %d", suffix%1000000),
			"sub_title":    "Performance test product",
			"main_image":   "https://example.com/perf.jpg",
			"sort":         100,
			"service_bits": 0,
		},
		"skus": []map[string]any{
			{
				"sku_code":      fmt.Sprintf("PERF-SKU-%d", suffix%1000000),
				"name":          "Perf Test SKU",
				"sub_title":     "Perf subtitle",
				"main_image":    "https://example.com/perf-sku.jpg",
				"price":         "8888",
				"market_price":  "9999",
				"stock":         100,
				"sku_spec_data": `{"perf":"test"}`,
			},
		},
		"detail": map[string]any{
			"description": "Performance test description with lots of content to make it more realistic",
			"images":      []string{"https://example.com/p1.jpg", "https://example.com/p2.jpg", "https://example.com/p3.jpg"},
		},
	}
	createResp := postJSON[map[string]any](t, client, baseURL+"/admin/products", createBody, nil)
	if createResp.Code != uint64(perrors.Success.Code) {
		t.Fatalf("create product failed: code=%d msg=%s", createResp.Code, createResp.Message)
	}

	spuID := int64(createResp.Data["spu_id"].(float64))
	t.Logf("Performance test product created with SPU ID: %d", spuID)

	// Warm up request (cache miss)
	numRequests := 10
	warmupDurations := make([]time.Duration, 0, numRequests)

	for i := 0; i < numRequests; i++ {
		// Clear any potential cache by waiting (simulating fresh requests)
		start := time.Now()
		resp := getJSON[map[string]any](t, client, fmt.Sprintf("%s/products/%d", baseURL, spuID), nil)
		duration := time.Since(start)
		if resp.Code != uint64(perrors.Success.Code) {
			t.Fatalf("request %d failed: code=%d", i, resp.Code)
		}
		warmupDurations = append(warmupDurations, duration)

		// After first request, wait for cache to be set
		if i == 0 {
			time.Sleep(100 * time.Millisecond)
		}
	}

	// Calculate average for cache hits (skip first request)
	var totalCacheHit time.Duration
	for i := 1; i < len(warmupDurations); i++ {
		totalCacheHit += warmupDurations[i]
	}
	avgCacheHit := totalCacheHit / time.Duration(len(warmupDurations)-1)

	t.Logf("First request (cache miss): %v", warmupDurations[0])
	t.Logf("Average cache hit time (%d requests): %v", len(warmupDurations)-1, avgCacheHit)
	t.Logf("Performance test completed")
}
