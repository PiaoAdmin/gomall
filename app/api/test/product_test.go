package test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	perrors "github.com/PiaoAdmin/pmall/common/errs"
)

// ========== Product Flow Test ==========

// TestProductFlow performs a complete product workflow test
func TestProductFlow(t *testing.T) {
	baseURL := getTestServer(t)
	client := &http.Client{Timeout: 5 * time.Second}

	suffix := time.Now().UnixNano()

	// 1) List categories
	categoriesResp := getJSON[map[string]any](t, client, baseURL+"/categories", nil)
	if categoriesResp.Code != uint64(perrors.Success.Code) {
		t.Fatalf("list categories failed: code=%d msg=%s", categoriesResp.Code, categoriesResp.Message)
	}
	t.Logf("Categories listed successfully")

	// 2) List brands with pagination
	brandsResp := getJSON[map[string]any](t, client, baseURL+"/brands?page=1&page_size=10", nil)
	if brandsResp.Code != uint64(perrors.Success.Code) {
		t.Fatalf("list brands failed: code=%d msg=%s", brandsResp.Code, brandsResp.Message)
	}
	t.Logf("Brands listed successfully")

	// 3) Create a product
	createBody := map[string]any{
		"spu": map[string]any{
			"brand_id":     1,
			"category_id":  1,
			"name":         fmt.Sprintf("Test Product %d", suffix%1000000),
			"sub_title":    "Test product subtitle",
			"main_image":   "https://example.com/image.jpg",
			"sort":         100,
			"service_bits": 0,
		},
		"skus": []map[string]any{
			{
				"sku_code":      fmt.Sprintf("SKU-%d-001", suffix%1000000),
				"name":          "Default SKU",
				"sub_title":     "Default SKU subtitle",
				"main_image":    "https://example.com/sku.jpg",
				"price":         "9999",
				"market_price":  "12999",
				"stock":         100,
				"sku_spec_data": `{"color":"red","size":"L"}`,
			},
		},
		"detail": map[string]any{
			"description":     "Test product description",
			"images":          []string{"https://example.com/detail1.jpg"},
			"videos":          []string{},
			"market_tag_json": `{"hot":true}`,
			"tech_tag_json":   `{"new":true}`,
		},
	}
	createResp := postJSON[map[string]any](t, client, baseURL+"/admin/products", createBody, nil)
	if createResp.Code != uint64(perrors.Success.Code) {
		t.Fatalf("create product failed: code=%d msg=%s", createResp.Code, createResp.Message)
	}

	spuID := int64(createResp.Data["spu_id"].(float64))
	t.Logf("Product created successfully with SPU ID: %d", spuID)

	// 4) Get home products
	homeResp := getJSON[map[string]any](t, client, baseURL+"/products/home?page=1&page_size=10", nil)
	if homeResp.Code != uint64(perrors.Success.Code) {
		t.Fatalf("get home products failed: code=%d msg=%s", homeResp.Code, homeResp.Message)
	}
	t.Logf("Home products retrieved successfully")

	// 5) Search products
	searchResp := getJSON[map[string]any](t, client, baseURL+"/products/search?page=1&page_size=10&keyword=Test", nil)
	if searchResp.Code != uint64(perrors.Success.Code) {
		t.Fatalf("search products failed: code=%d msg=%s", searchResp.Code, searchResp.Message)
	}
	t.Logf("Products searched successfully")

	// 6) Get product detail
	detailResp := getJSON[map[string]any](t, client, fmt.Sprintf("%s/products/%d", baseURL, spuID), nil)
	if detailResp.Code != uint64(perrors.Success.Code) {
		t.Fatalf("get product detail failed: code=%d msg=%s", detailResp.Code, detailResp.Message)
	}
	t.Logf("Product detail retrieved successfully")

	// 7) Batch update SKU (get SKU ID from detail response first)
	product := detailResp.Data["product"].(map[string]any)
	skus := product["skus"].([]any)
	if len(skus) > 0 {
		sku := skus[0].(map[string]any)
		skuID := int64(sku["id"].(float64))

		updateBody := map[string]any{
			"items": []map[string]any{
				{
					"sku_id": skuID,
					"price":  "8888",
					"stock":  200,
				},
			},
		}
		updateResp := postJSON[map[string]any](t, client, baseURL+"/admin/skus/batch", updateBody, nil)
		if updateResp.Code != uint64(perrors.Success.Code) {
			t.Fatalf("batch update sku failed: code=%d msg=%s", updateResp.Code, updateResp.Message)
		}
		t.Logf("SKU updated successfully")
	}
}

// ========== List Categories Tests ==========

func TestListCategories(t *testing.T) {
	baseURL := getTestServer(t)
	client := &http.Client{Timeout: 5 * time.Second}

	resp := getJSON[map[string]any](t, client, baseURL+"/categories", nil)
	if resp.Code != uint64(perrors.Success.Code) {
		t.Fatalf("list categories failed: code=%d msg=%s", resp.Code, resp.Message)
	}
	t.Logf("%v", resp.Data)
	if resp.Data["categories"] == nil {
		t.Fatalf("categories data is nil")
	}
}

func TestListCategoriesWithParent(t *testing.T) {
	baseURL := getTestServer(t)
	client := &http.Client{Timeout: 5 * time.Second}

	// Get root categories
	resp := getJSON[map[string]any](t, client, baseURL+"/categories?parent_id=0", nil)
	if resp.Code != uint64(perrors.Success.Code) {
		t.Fatalf("list root categories failed: code=%d msg=%s", resp.Code, resp.Message)
	}

	categories := resp.Data["categories"]
	if categories == nil {
		t.Logf("no categories found, which is acceptable for new database")
		return
	}

	// If there are categories, try to get children of first one
	catList := categories.([]any)
	if len(catList) > 0 {
		firstCat := catList[0].(map[string]any)
		parentID := int64(firstCat["id"].(float64))

		childResp := getJSON[map[string]any](t, client, fmt.Sprintf("%s/categories?parent_id=%d", baseURL, parentID), nil)
		if childResp.Code != uint64(perrors.Success.Code) {
			t.Fatalf("list child categories failed: code=%d msg=%s", childResp.Code, childResp.Message)
		}
	}
}

// ========== List Brands Tests ==========

func TestListBrands(t *testing.T) {
	baseURL := getTestServer(t)
	client := &http.Client{Timeout: 5 * time.Second}

	resp := getJSON[map[string]any](t, client, baseURL+"/brands", nil)
	if resp.Code != uint64(perrors.Success.Code) {
		t.Fatalf("list brands failed: code=%d msg=%s", resp.Code, resp.Message)
	}

	if resp.Data["brands"] == nil {
		t.Fatalf("brands data is nil")
	}
}

func TestListBrandsWithPagination(t *testing.T) {
	baseURL := getTestServer(t)
	client := &http.Client{Timeout: 5 * time.Second}

	testCases := []struct {
		page     int
		pageSize int
	}{
		{1, 10},
		{1, 20},
		{2, 10},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("page=%d_size=%d", tc.page, tc.pageSize), func(t *testing.T) {
			url := fmt.Sprintf("%s/brands?page=%d&page_size=%d", baseURL, tc.page, tc.pageSize)
			resp := getJSON[map[string]any](t, client, url, nil)
			if resp.Code != uint64(perrors.Success.Code) {
				t.Fatalf("list brands failed: code=%d msg=%s", resp.Code, resp.Message)
			}
		})
	}
}

// ========== Get Home Products Tests ==========

func TestGetHomeProducts(t *testing.T) {
	baseURL := getTestServer(t)
	client := &http.Client{Timeout: 5 * time.Second}

	resp := getJSON[map[string]any](t, client, baseURL+"/products/home?page=1&page_size=10", nil)
	if resp.Code != uint64(perrors.Success.Code) {
		t.Fatalf("get home products failed: code=%d msg=%s", resp.Code, resp.Message)
	}
	t.Logf("%v", resp.Data)
	if resp.Data["list"] == nil {
		t.Fatalf("product list is nil")
	}
}

func TestGetHomeProductsWithFilters(t *testing.T) {
	baseURL := getTestServer(t)
	client := &http.Client{Timeout: 5 * time.Second}

	testCases := []struct {
		name   string
		params string
	}{
		{"with category", "page=1&page_size=10&category_id=1"},
		{"with brand", "page=1&page_size=10&brand_id=1"},
		{"with both filters", "page=1&page_size=10&category_id=1&brand_id=1"},
		{"large page size", "page=1&page_size=50"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("%s/products/home?%s", baseURL, tc.params)
			resp := getJSON[map[string]any](t, client, url, nil)
			t.Logf("%v", resp.Data)
			if resp.Code != uint64(perrors.Success.Code) {
				t.Fatalf("get home products failed: code=%d msg=%s", resp.Code, resp.Message)
			}
		})
	}
}

// ========== Search Products Tests ==========

func TestSearchProducts(t *testing.T) {
	baseURL := getTestServer(t)
	client := &http.Client{Timeout: 5 * time.Second}

	resp := getJSON[map[string]any](t, client, baseURL+"/products/search?page=1&page_size=10", nil)
	if resp.Code != uint64(perrors.Success.Code) {
		t.Fatalf("search products failed: code=%d msg=%s", resp.Code, resp.Message)
	}

	if resp.Data["list"] == nil {
		t.Fatalf("search result list is nil")
	}
}

func TestSearchProductsWithKeyword(t *testing.T) {
	baseURL := getTestServer(t)
	client := &http.Client{Timeout: 5 * time.Second}

	keywords := []string{"手机", "电脑", "Test", "Product"}

	for _, keyword := range keywords {
		t.Run(keyword, func(t *testing.T) {
			url := fmt.Sprintf("%s/products/search?page=1&page_size=10&keyword=%s", baseURL, keyword)
			resp := getJSON[map[string]any](t, client, url, nil)
			if resp.Code != uint64(perrors.Success.Code) {
				t.Fatalf("search products failed: code=%d msg=%s", resp.Code, resp.Message)
			}
		})
	}
}

func TestSearchProductsWithPriceRange(t *testing.T) {
	baseURL := getTestServer(t)
	client := &http.Client{Timeout: 5 * time.Second}

	testCases := []struct {
		name     string
		minPrice int64
		maxPrice int64
	}{
		{"low range", 0, 5000},
		{"mid range", 5000, 10000},
		{"high range", 10000, 99999},
		{"only min", 1000, 0},
		{"only max", 0, 5000},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("%s/products/search?page=1&page_size=10&min_price=%d&max_price=%d",
				baseURL, tc.minPrice, tc.maxPrice)
			resp := getJSON[map[string]any](t, client, url, nil)
			if resp.Code != uint64(perrors.Success.Code) {
				t.Fatalf("search products failed: code=%d msg=%s", resp.Code, resp.Message)
			}
		})
	}
}

func TestSearchProductsWithSort(t *testing.T) {
	baseURL := getTestServer(t)
	client := &http.Client{Timeout: 5 * time.Second}

	sortTypes := []struct {
		name     string
		sortType int
	}{
		{"default", 0},
		{"price asc", 1},
		{"price desc", 2},
		{"sales desc", 3},
	}

	for _, st := range sortTypes {
		t.Run(st.name, func(t *testing.T) {
			url := fmt.Sprintf("%s/products/search?page=1&page_size=10&sort_type=%d", baseURL, st.sortType)
			resp := getJSON[map[string]any](t, client, url, nil)

			if resp.Code != uint64(perrors.Success.Code) {
				t.Fatalf("search products failed: code=%d msg=%s", resp.Code, resp.Message)
			}
			t.Logf("Sort type: %s", st.name)
			if resp.Data["list"] != nil {
				list := resp.Data["list"].([]any)
				t.Logf("Found %d products", len(list))
				for i, item := range list {
					product := item.(map[string]any)
					t.Logf("  [%d] SKU ID: %v, Name: %v, Price: %v",
						i+1, product["sku_id"], product["name"], product["price"])
				}
			} else {
				t.Logf("No products found")
			}
		})
	}
}

// ========== Get Product Detail Tests ==========

func TestGetProductDetail(t *testing.T) {
	baseURL := getTestServer(t)
	client := &http.Client{Timeout: 5 * time.Second}

	// Create a test product first
	suffix := time.Now().UnixNano()

	createBody := map[string]any{
		"spu": map[string]any{
			"brand_id":     1,
			"category_id":  1,
			"name":         fmt.Sprintf("Detail Test Product %d", suffix%1000000),
			"sub_title":    "Detail test subtitle",
			"main_image":   "https://example.com/image.jpg",
			"sort":         100,
			"service_bits": 0,
		},
		"skus": []map[string]any{
			{
				"sku_code":      fmt.Sprintf("SKU-%d-DT", suffix%1000000),
				"name":          "Detail Test SKU",
				"sub_title":     "SKU subtitle",
				"main_image":    "https://example.com/sku.jpg",
				"price":         "5999",
				"market_price":  "7999",
				"stock":         50,
				"sku_spec_data": `{"color":"blue"}`,
			},
		},
		"detail": map[string]any{
			"description":     "Detailed description",
			"images":          []string{"https://example.com/d1.jpg"},
			"videos":          []string{},
			"market_tag_json": `{}`,
			"tech_tag_json":   `{}`,
		},
	}
	createResp := postJSON[map[string]any](t, client, baseURL+"/admin/products", createBody, nil)
	if createResp.Code != uint64(perrors.Success.Code) {
		t.Fatalf("create product failed: code=%d msg=%s", createResp.Code, createResp.Message)
	}

	spuID := int64(createResp.Data["spu_id"].(float64))

	// Get product detail
	detailResp := getJSON[map[string]any](t, client, fmt.Sprintf("%s/products/%d", baseURL, spuID), nil)
	if detailResp.Code != uint64(perrors.Success.Code) {
		t.Fatalf("get product detail failed: code=%d msg=%s", detailResp.Code, detailResp.Message)
	}

	if detailResp.Data["product"] == nil {
		t.Fatalf("product detail is nil")
	}

	product := detailResp.Data["product"].(map[string]any)
	if product["id"] != float64(spuID) {
		t.Errorf("product ID mismatch: expected %d, got %v", spuID, product["id"])
	}
}

func TestGetProductDetailNonexistent(t *testing.T) {
	baseURL := getTestServer(t)
	client := &http.Client{Timeout: 5 * time.Second}

	// Try to get a product that doesn't exist
	nonexistentID := int64(999999999)
	resp := getJSON[any](t, client, fmt.Sprintf("%s/products/%d", baseURL, nonexistentID), nil)

	if resp.Code == uint64(perrors.Success.Code) {
		t.Fatalf("expected get nonexistent product to fail, but succeeded")
	}
}

// ========== Create Product Tests ==========

func TestCreateProduct(t *testing.T) {
	baseURL := getTestServer(t)
	client := &http.Client{Timeout: 5 * time.Second}

	suffix := time.Now().UnixNano()

	createBody := map[string]any{
		"spu": map[string]any{
			"brand_id":     1,
			"category_id":  1,
			"name":         fmt.Sprintf("New Product %d", suffix%1000000),
			"sub_title":    "Product subtitle",
			"main_image":   "https://example.com/main.jpg",
			"sort":         100,
			"service_bits": 0,
		},
		"skus": []map[string]any{
			{
				"sku_code":      fmt.Sprintf("SKU-%d-NEW", suffix%1000000),
				"name":          "New SKU",
				"sub_title":     "SKU subtitle",
				"main_image":    "https://example.com/sku.jpg",
				"price":         "9999",
				"market_price":  "12999",
				"stock":         100,
				"sku_spec_data": `{"color":"red","size":"M"}`,
			},
		},
		"detail": map[string]any{
			"description":     "Product description",
			"images":          []string{"https://example.com/img1.jpg", "https://example.com/img2.jpg"},
			"videos":          []string{"https://example.com/video1.mp4"},
			"market_tag_json": `{"hot":true,"new":true}`,
			"tech_tag_json":   `{"5G":true}`,
		},
	}

	resp := postJSON[map[string]any](t, client, baseURL+"/admin/products", createBody, nil)
	if resp.Code != uint64(perrors.Success.Code) {
		t.Fatalf("create product failed: code=%d msg=%s", resp.Code, resp.Message)
	}

	if resp.Data["spu_id"] == nil {
		t.Fatalf("spu_id is nil in response")
	}
}

func TestCreateProductWithMultipleSKUs(t *testing.T) {
	baseURL := getTestServer(t)
	client := &http.Client{Timeout: 5 * time.Second}

	suffix := time.Now().UnixNano()

	createBody := map[string]any{
		"spu": map[string]any{
			"brand_id":     1,
			"category_id":  1,
			"name":         fmt.Sprintf("Multi SKU Product %d", suffix%1000000),
			"sub_title":    "Product with multiple SKUs",
			"main_image":   "https://example.com/main.jpg",
			"sort":         100,
			"service_bits": 0,
		},
		"skus": []map[string]any{
			{
				"sku_code":      fmt.Sprintf("SKU-%d-R-S", suffix%1000000),
				"name":          "Red Small",
				"sub_title":     "Red color, small size",
				"main_image":    "https://example.com/red-s.jpg",
				"price":         "8999",
				"market_price":  "10999",
				"stock":         50,
				"sku_spec_data": `{"color":"red","size":"S"}`,
			},
			{
				"sku_code":      fmt.Sprintf("SKU-%d-R-M", suffix%1000000),
				"name":          "Red Medium",
				"sub_title":     "Red color, medium size",
				"main_image":    "https://example.com/red-m.jpg",
				"price":         "9999",
				"market_price":  "11999",
				"stock":         100,
				"sku_spec_data": `{"color":"red","size":"M"}`,
			},
			{
				"sku_code":      fmt.Sprintf("SKU-%d-B-M", suffix%1000000),
				"name":          "Blue Medium",
				"sub_title":     "Blue color, medium size",
				"main_image":    "https://example.com/blue-m.jpg",
				"price":         "9999",
				"market_price":  "11999",
				"stock":         80,
				"sku_spec_data": `{"color":"blue","size":"M"}`,
			},
		},
		"detail": map[string]any{
			"description":     "Multi-SKU product description",
			"images":          []string{"https://example.com/img1.jpg"},
			"videos":          []string{},
			"market_tag_json": `{}`,
			"tech_tag_json":   `{}`,
		},
	}

	resp := postJSON[map[string]any](t, client, baseURL+"/admin/products", createBody, nil)
	if resp.Code != uint64(perrors.Success.Code) {
		t.Fatalf("create product with multiple SKUs failed: code=%d msg=%s", resp.Code, resp.Message)
	}
}

// func TestCreateProductMissingFields(t *testing.T) {
// 	baseURL := getTestServer(t)
// 	client := &http.Client{Timeout: 5 * time.Second}
//
// 	testCases := []struct {
// 		name string
// 		body map[string]any
// 	}{
// 		{
// 			"missing spu",
// 			map[string]any{
// 				"skus":   []map[string]any{},
// 				"detail": map[string]any{},
// 			},
// 		},
// 		{
// 			"missing skus",
// 			map[string]any{
// 				"spu": map[string]any{
// 					"brand_id":    1,
// 					"category_id": 1,
// 					"name":        "Test",
// 				},
// 				"detail": map[string]any{},
// 			},
// 		},
// 		{
// 			"empty skus array",
// 			map[string]any{
// 				"spu": map[string]any{
// 					"brand_id":    1,
// 					"category_id": 1,
// 					"name":        "Test",
// 				},
// 				"skus":   []map[string]any{},
// 				"detail": map[string]any{},
// 			},
// 		},
// 	}
//
// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			resp := postJSON[any](t, client, baseURL+"/admin/products", tc.body, nil)
// 			if resp.Code == uint64(perrors.Success.Code) {
// 				t.Fatalf("expected create product with %s to fail, but succeeded", tc.name)
// 			}
// 		})
// 	}
// }
//
// ========== Batch Update SKU Tests ==========

func TestBatchUpdateSku(t *testing.T) {
	baseURL := getTestServer(t)
	client := &http.Client{Timeout: 5 * time.Second}

	// Create product first
	suffix := time.Now().UnixNano()

	createBody := map[string]any{
		"spu": map[string]any{
			"brand_id":     1,
			"category_id":  1,
			"name":         fmt.Sprintf("Update Test Product %d", suffix%1000000),
			"sub_title":    "For update testing",
			"main_image":   "https://example.com/main.jpg",
			"sort":         100,
			"service_bits": 0,
		},
		"skus": []map[string]any{
			{
				"sku_code":      fmt.Sprintf("SKU-%d-UP", suffix%1000000),
				"name":          "Update Test SKU",
				"sub_title":     "SKU for update",
				"main_image":    "https://example.com/sku.jpg",
				"price":         "5000",
				"market_price":  "7000",
				"stock":         50,
				"sku_spec_data": `{}`,
			},
		},
		"detail": map[string]any{
			"description":     "Test",
			"images":          []string{},
			"videos":          []string{},
			"market_tag_json": `{}`,
			"tech_tag_json":   `{}`,
		},
	}
	createResp := postJSON[map[string]any](t, client, baseURL+"/admin/products", createBody, nil)
	if createResp.Code != uint64(perrors.Success.Code) {
		t.Fatalf("create product failed: code=%d msg=%s", createResp.Code, createResp.Message)
	}

	spuID := int64(createResp.Data["spu_id"].(float64))

	// Get product detail to get SKU ID
	detailResp := getJSON[map[string]any](t, client, fmt.Sprintf("%s/products/%d", baseURL, spuID), nil)
	if detailResp.Code != uint64(perrors.Success.Code) {
		t.Fatalf("get product detail failed: code=%d msg=%s", detailResp.Code, detailResp.Message)
	}

	product := detailResp.Data["product"].(map[string]any)
	skus := product["skus"].([]any)
	if len(skus) == 0 {
		t.Fatalf("no SKUs found in created product")
	}

	sku := skus[0].(map[string]any)
	skuID := int64(sku["id"].(float64))

	// Batch update SKU
	updateBody := map[string]any{
		"items": []map[string]any{
			{
				"sku_id": skuID,
				"price":  "6000",
				"stock":  100,
			},
		},
	}
	updateResp := postJSON[map[string]any](t, client, baseURL+"/admin/skus/batch", updateBody, nil)
	if updateResp.Code != uint64(perrors.Success.Code) {
		t.Fatalf("batch update sku failed: code=%d msg=%s", updateResp.Code, updateResp.Message)
	}

	// Verify update
	detailResp2 := getJSON[map[string]any](t, client, fmt.Sprintf("%s/products/%d", baseURL, spuID), nil)
	if detailResp2.Code != uint64(perrors.Success.Code) {
		t.Fatalf("get product detail after update failed: code=%d msg=%s", detailResp2.Code, detailResp2.Message)
	}

	product2 := detailResp2.Data["product"].(map[string]any)
	skus2 := product2["skus"].([]any)

	// Debug: Print full SKU data to understand structure
	t.Logf("Number of SKUs after update: %d", len(skus2))
	if len(skus2) == 0 {
		t.Fatalf("no SKUs found after update")
	}

	sku2 := skus2[0].(map[string]any)
	t.Logf("SKU data after update: %+v", sku2)

	// Check if price field exists and has correct value
	if sku2["price"] == nil {
		t.Fatalf("price is nil after update, full SKU data: %+v", sku2)
	}

	// Price should be string type
	priceStr, ok := sku2["price"].(string)
	if !ok {
		t.Fatalf("price is not string type, got type %T with value %v", sku2["price"], sku2["price"])
	}
	if priceStr != "6000.00" {
		t.Errorf("price not updated: expected \"6000.00\", got %q", priceStr)
	}

	// Check stock
	if sku2["stock"] == nil {
		t.Fatalf("stock is nil after update")
	}
	if int64(sku2["stock"].(float64)) != 100 {
		t.Errorf("stock not updated: expected 100, got %v", sku2["stock"])
	}

	t.Logf("SKU update verification passed: price=%s, stock=%v", priceStr, sku2["stock"])
}

func TestBatchUpdateMultipleSKUs(t *testing.T) {
	baseURL := getTestServer(t)
	client := &http.Client{Timeout: 5 * time.Second}

	// Create product with multiple SKUs
	suffix := time.Now().UnixNano()

	createBody := map[string]any{
		"spu": map[string]any{
			"brand_id":     1,
			"category_id":  1,
			"name":         fmt.Sprintf("Batch Update Product %d", suffix%1000000),
			"sub_title":    "For batch update testing",
			"main_image":   "https://example.com/main.jpg",
			"sort":         100,
			"service_bits": 0,
		},
		"skus": []map[string]any{
			{
				"sku_code":      fmt.Sprintf("SKU-%d-1", suffix%1000000),
				"name":          "SKU 1",
				"sub_title":     "First SKU",
				"main_image":    "https://example.com/sku1.jpg",
				"price":         "1000",
				"market_price":  "1500",
				"stock":         10,
				"sku_spec_data": `{"variant":"1"}`,
			},
			{
				"sku_code":      fmt.Sprintf("SKU-%d-2", suffix%1000000),
				"name":          "SKU 2",
				"sub_title":     "Second SKU",
				"main_image":    "https://example.com/sku2.jpg",
				"price":         "2000",
				"market_price":  "2500",
				"stock":         20,
				"sku_spec_data": `{"variant":"2"}`,
			},
		},
		"detail": map[string]any{
			"description":     "Test",
			"images":          []string{},
			"videos":          []string{},
			"market_tag_json": `{}`,
			"tech_tag_json":   `{}`,
		},
	}
	createResp := postJSON[map[string]any](t, client, baseURL+"/admin/products", createBody, nil)
	if createResp.Code != uint64(perrors.Success.Code) {
		t.Fatalf("create product failed: code=%d msg=%s", createResp.Code, createResp.Message)
	}

	spuID := int64(createResp.Data["spu_id"].(float64))

	// Get SKU IDs
	detailResp := getJSON[map[string]any](t, client, fmt.Sprintf("%s/products/%d", baseURL, spuID), nil)
	product := detailResp.Data["product"].(map[string]any)
	skus := product["skus"].([]any)

	skuID1 := int64(skus[0].(map[string]any)["id"].(float64))
	skuID2 := int64(skus[1].(map[string]any)["id"].(float64))

	// Batch update both SKUs
	updateBody := map[string]any{
		"items": []map[string]any{
			{
				"sku_id": skuID1,
				"price":  "1200",
				"stock":  15,
			},
			{
				"sku_id": skuID2,
				"price":  "2200",
				"stock":  25,
			},
		},
	}
	updateResp := postJSON[map[string]any](t, client, baseURL+"/admin/skus/batch", updateBody, nil)
	if updateResp.Code != uint64(perrors.Success.Code) {
		t.Fatalf("batch update multiple skus failed: code=%d msg=%s", updateResp.Code, updateResp.Message)
	}

	if int(updateResp.Data["updated_count"].(float64)) != 2 {
		t.Errorf("expected updated_count=2, got %v", updateResp.Data["updated_count"])
	}
}

func TestBatchUpdateSkuNonexistent(t *testing.T) {
	baseURL := getTestServer(t)
	client := &http.Client{Timeout: 5 * time.Second}

	updateBody := map[string]any{
		"items": []map[string]any{
			{
				"sku_id": 999999999,
				"price":  "9999",
				"stock":  100,
			},
		},
	}

	resp := postJSON[any](t, client, baseURL+"/admin/skus/batch", updateBody, nil)
	if resp.Code == uint64(perrors.Success.Code) {
		t.Fatalf("expected batch update nonexistent SKU to fail, but succeeded")
	}
}

// ========== Edge Cases and Validation Tests ==========

func TestGetHomeProductsInvalidPagination(t *testing.T) {
	baseURL := getTestServer(t)
	client := &http.Client{Timeout: 5 * time.Second}

	testCases := []struct {
		name   string
		params string
	}{
		{"negative page", "page=-1&page_size=10"},
		{"zero page", "page=0&page_size=10"},
		{"negative page_size", "page=1&page_size=-10"},
		{"zero page_size", "page=1&page_size=0"},
		{"huge page_size", "page=1&page_size=10000"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("%s/products/home?%s", baseURL, tc.params)
			resp := getJSON[any](t, client, url, nil)
			// Should either handle gracefully or return error
			t.Logf("%s: code=%d msg=%s", tc.name, resp.Code, resp.Message)
		})
	}
}

func TestSearchProductsEmptyKeyword(t *testing.T) {
	baseURL := getTestServer(t)
	client := &http.Client{Timeout: 5 * time.Second}

	resp := getJSON[map[string]any](t, client, baseURL+"/products/search?page=1&page_size=10&keyword=", nil)
	if resp.Code != uint64(perrors.Success.Code) {
		t.Fatalf("search with empty keyword failed: code=%d msg=%s", resp.Code, resp.Message)
	}
}

func TestSearchProductsInvalidSortType(t *testing.T) {
	baseURL := getTestServer(t)
	client := &http.Client{Timeout: 5 * time.Second}

	resp := getJSON[any](t, client, baseURL+"/products/search?page=1&page_size=10&sort_type=999", nil)
	// Should handle invalid sort type gracefully
	t.Logf("Invalid sort type response: code=%d msg=%s", resp.Code, resp.Message)
}

func TestGetProductDetailInvalidID(t *testing.T) {
	baseURL := getTestServer(t)
	client := &http.Client{Timeout: 5 * time.Second}

	testCases := []struct {
		name string
		id   string
	}{
		{"negative ID", "-1"},
		{"zero ID", "0"},
		{"non-numeric ID", "abc"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp := getJSON[any](t, client, fmt.Sprintf("%s/products/%s", baseURL, tc.id), nil)
			if resp.Code == uint64(perrors.Success.Code) {
				t.Fatalf("Warning: invalid product ID %s accepted", tc.id)
			}
			t.Logf("Get product detail with %s ID response: code=%d msg=%s", tc.name, resp.Code, resp.Message)
		})
	}
}
