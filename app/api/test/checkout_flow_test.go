package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"
)

type checkoutRespEnvelope[T any] struct {
	Code    uint64 `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

type loginResp struct {
	Token string `json:"token"`
}

type searchListResp struct {
	List []struct {
		SkuId     uint64 `json:"sku_id"`
		SpuId     uint64 `json:"spu_id"`
		Stock     int32  `json:"stock"`
		SaleCount int32  `json:"sale_count"`
	} `json:"list"`
	Total int64 `json:"total"`
}

type checkoutResp struct {
	OrderId string `json:"order_id"`
}

// 使用可通过 Luhn 校验的测试卡号
const validCreditCard = "4111111111111111"

func TestCheckoutFlow(t *testing.T) {
	client := &http.Client{Timeout: 10 * time.Second}
	base := testBaseURL

	// 登录获取 token
	authHeader := loginForAuth(t, client, base)

	// 创建测试商品，确保库存可控
	spuID, skuID := createTestProduct(t, client, base, time.Now().UnixNano())
	stockBefore, saleBefore := getProductStockAndSales(t, client, base, spuID, skuID)
	if stockBefore < 1 {
		t.Skip("sku stock < 1, skip stock/sale assertion")
	}

	// 清空购物车，避免历史数据影响
	clearEnv := postJSONEnvelope[map[string]any](t, client, base+"/cart/clear", map[string]any{}, authHeader)
	if clearEnv.Code != 20000 {
		t.Fatalf("clear cart failed: code=%d msg=%s", clearEnv.Code, clearEnv.Message)
	}

	addBody := map[string]any{
		"sku_id":   skuID,
		"quantity": 1,
	}
	// 加入购物车
	addEnv := postJSONEnvelope[map[string]any](t, client, base+"/cart/add", addBody, authHeader)
	if addEnv.Code != 20000 {
		t.Fatalf("add to cart failed: code=%d msg=%s", addEnv.Code, addEnv.Message)
	}

	checkoutBody := map[string]any{
		"shipping_address": map[string]any{
			"name":           "piao",
			"street_address": "test street",
			"city":           "test city",
			"zip_code":       12345,
		},
		"credit_card": validCreditCard,
	}
	// 直接 checkout 下单并支付
	checkoutEnv := postJSONEnvelope[checkoutResp](t, client, base+"/checkout", checkoutBody, authHeader)
	if checkoutEnv.Code != 20000 {
		t.Fatalf("checkout failed: code=%d msg=%s", checkoutEnv.Code, checkoutEnv.Message)
	}
	if checkoutEnv.Data.OrderId == "" {
		t.Fatalf("checkout order_id empty")
	}

	// 校验库存-1、销量+1
	stockAfter, saleAfter := getProductStockAndSales(t, client, base, spuID, skuID)
	if stockAfter != stockBefore-1 {
		t.Fatalf("stock not decreased: before=%d after=%d", stockBefore, stockAfter)
	}
	if saleAfter != saleBefore+1 {
		t.Fatalf("sale_count not increased: before=%d after=%d", saleBefore, saleAfter)
	}
}

func TestCheckoutEdgeCases(t *testing.T) {
	client := &http.Client{Timeout: 10 * time.Second}
	base := testBaseURL

	t.Run("missing auth", func(t *testing.T) {
		// 缺少鉴权 header
		env := postJSONEnvelope[map[string]any](t, client, base+"/checkout", map[string]any{}, nil)
		if env.Code != 40005 {
			t.Fatalf("expected auth failed code=40005, got=%d msg=%s", env.Code, env.Message)
		}
	})

	authHeader := loginForAuth(t, client, base)
	skuID, _, _, _ := fetchFirstSkuInfo(t, client, base)

	// 清空购物车，制造空购物车场景
	clearEnv := postJSONEnvelope[map[string]any](t, client, base+"/cart/clear", map[string]any{}, authHeader)
	if clearEnv.Code != 20000 {
		t.Fatalf("clear cart failed: code=%d msg=%s", clearEnv.Code, clearEnv.Message)
	}

	t.Run("empty cart", func(t *testing.T) {
		// 空购物车下单
		checkoutBody := map[string]any{
			"shipping_address": map[string]any{
				"name":           "piao",
				"street_address": "test street",
				"city":           "test city",
				"zip_code":       12345,
			},
		}
		env := postJSONEnvelope[map[string]any](t, client, base+"/checkout", checkoutBody, authHeader)
		if env.Code != 40007 {
			t.Fatalf("expected empty cart code=40007, got=%d msg=%s", env.Code, env.Message)
		}
	})

	addBody := map[string]any{
		"sku_id":   skuID,
		"quantity": 1,
	}
	addEnv := postJSONEnvelope[map[string]any](t, client, base+"/cart/add", addBody, authHeader)
	if addEnv.Code != 20000 {
		t.Fatalf("add to cart failed: code=%d msg=%s", addEnv.Code, addEnv.Message)
	}

	t.Run("missing shipping address", func(t *testing.T) {
		// 缺少收货地址
		env := postJSONEnvelope[map[string]any](t, client, base+"/checkout", map[string]any{}, authHeader)
		if env.Code != 40002 {
			t.Fatalf("expected invalid params code=40002, got=%d msg=%s", env.Code, env.Message)
		}
	})

	t.Run("bind error on bad shipping_address type", func(t *testing.T) {
		// shipping_address 字段类型错误
		env := postRawJSONEnvelope[map[string]any](t, client, base+"/checkout", `{"shipping_address":"bad"}`, authHeader)
		if env.Code != 40001 {
			t.Fatalf("expected bind error code=40001, got=%d msg=%s", env.Code, env.Message)
		}
	})
}

func TestCheckoutStockInsufficient(t *testing.T) {
	client := &http.Client{Timeout: 10 * time.Second}
	base := testBaseURL

	// 登录获取 token
	authHeader := loginForAuth(t, client, base)
	skuID, spuID, stockBefore, saleBefore := fetchFirstSkuInfo(t, client, base)
	if stockBefore <= 0 {
		t.Skip("sku stock <= 0, skip stock insufficient test")
	}

	// 清空购物车，保证只测一个 SKU 的库存不足
	clearEnv := postJSONEnvelope[map[string]any](t, client, base+"/cart/clear", map[string]any{}, authHeader)
	if clearEnv.Code != 20000 {
		t.Fatalf("clear cart failed: code=%d msg=%s", clearEnv.Code, clearEnv.Message)
	}

	addBody := map[string]any{
		"sku_id":   skuID,
		"quantity": int(stockBefore) + 1,
	}
	addEnv := postJSONEnvelope[map[string]any](t, client, base+"/cart/add", addBody, authHeader)
	if addEnv.Code != 20000 {
		t.Fatalf("add to cart failed: code=%d msg=%s", addEnv.Code, addEnv.Message)
	}

	checkoutBody := map[string]any{
		"shipping_address": map[string]any{
			"name":           "piao",
			"street_address": "test street",
			"city":           "test city",
			"zip_code":       12345,
		},
		"credit_card": validCreditCard,
	}
	// 下单数量超过库存，期望库存不足
	checkoutEnv := postJSONEnvelope[checkoutResp](t, client, base+"/checkout", checkoutBody, authHeader)
	if checkoutEnv.Code != 40002 {
		t.Fatalf("expected stock not enough code=40002, got=%d msg=%s", checkoutEnv.Code, checkoutEnv.Message)
	}

	stockAfter, saleAfter := getProductStockAndSales(t, client, base, spuID, skuID)
	if stockAfter != stockBefore {
		t.Fatalf("stock changed on insufficient checkout: before=%d after=%d", stockBefore, stockAfter)
	}
	if saleAfter != saleBefore {
		t.Fatalf("sale_count changed on insufficient checkout: before=%d after=%d", saleBefore, saleAfter)
	}
}

func TestCheckoutPaymentFailRollbackStock(t *testing.T) {
	client := &http.Client{Timeout: 10 * time.Second}
	base := testBaseURL

	// 创建并登录测试用户，避免与其他用例购物车冲突
	suffix := time.Now().UnixNano()
	_, _, token := createAndLoginTestUser(t, client, base, suffix)
	authHeader := map[string]string{"Authorization": fmt.Sprintf("Bearer %s", token)}

	// 查询一个可购买 SKU
	skuID, spuID, stockBefore, _ := fetchFirstSkuInfo(t, client, base)
	if stockBefore <= 0 {
		t.Skip("sku stock <= 0, skip payment fail test")
	}

	addBody := map[string]any{
		"sku_id":   skuID,
		"quantity": 1,
	}
	addEnv := postJSONEnvelope[map[string]any](t, client, base+"/cart/add", addBody, authHeader)
	if addEnv.Code != 20000 {
		t.Fatalf("add to cart failed: code=%d msg=%s", addEnv.Code, addEnv.Message)
	}

	checkoutBody := map[string]any{
		"shipping_address": map[string]any{
			"name":           "piao",
			"street_address": "test street",
			"city":           "test city",
			"zip_code":       12345,
		},
		"credit_card": "123",
	}
	// 使用非法信用卡，触发支付失败
	checkoutEnv := postJSONEnvelope[checkoutResp](t, client, base+"/checkout", checkoutBody, authHeader)
	if checkoutEnv.Code != 40002 {
		t.Fatalf("expected payment failed code=40002, got=%d msg=%s", checkoutEnv.Code, checkoutEnv.Message)
	}

	// 支付失败后库存应回滚
	stockAfter, _ := getProductStockAndSales(t, client, base, spuID, skuID)
	if stockAfter != stockBefore {
		t.Fatalf("stock not rolled back after payment failure: before=%d after=%d", stockBefore, stockAfter)
	}
}

func postJSONEnvelope[T any](t *testing.T, client *http.Client, url string, body any, headers map[string]string) checkoutRespEnvelope[T] {
	t.Helper()
	buf, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal body: %v", err)
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(buf))
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("do request: %v", err)
	}
	defer resp.Body.Close()

	var env checkoutRespEnvelope[T]
	if err := json.NewDecoder(resp.Body).Decode(&env); err != nil {
		t.Fatalf("decode resp: %v", err)
	}
	return env
}

func postRawJSONEnvelope[T any](t *testing.T, client *http.Client, url string, body string, headers map[string]string) checkoutRespEnvelope[T] {
	t.Helper()
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader([]byte(body)))
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("do request: %v", err)
	}
	defer resp.Body.Close()

	var env checkoutRespEnvelope[T]
	if err := json.NewDecoder(resp.Body).Decode(&env); err != nil {
		t.Fatalf("decode resp: %v", err)
	}
	return env
}

func getJSONEnvelope[T any](t *testing.T, client *http.Client, url string, headers map[string]string) checkoutRespEnvelope[T] {
	t.Helper()
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("do request: %v", err)
	}
	defer resp.Body.Close()

	var env checkoutRespEnvelope[T]
	if err := json.NewDecoder(resp.Body).Decode(&env); err != nil {
		t.Fatalf("decode resp: %v", err)
	}
	return env
}

func loginForAuth(t *testing.T, client *http.Client, base string) map[string]string {
	t.Helper()
	loginBody := map[string]any{
		"username": "piao",
		"password": "123456",
	}
	loginEnv := postJSONEnvelope[loginResp](t, client, base+"/login", loginBody, nil)
	if loginEnv.Code != 20000 {
		t.Fatalf("login failed: code=%d msg=%s", loginEnv.Code, loginEnv.Message)
	}
	if loginEnv.Data.Token == "" {
		t.Fatalf("login token empty")
	}
	return map[string]string{"Authorization": fmt.Sprintf("Bearer %s", loginEnv.Data.Token)}
}

func fetchFirstSkuInfo(t *testing.T, client *http.Client, base string) (uint64, uint64, int32, int32) {
	t.Helper()
	query := url.Values{}
	query.Set("page", "1")
	query.Set("page_size", "1")
	searchURL := base + "/products/search?" + query.Encode()
	searchEnv := getJSONEnvelope[searchListResp](t, client, searchURL, nil)
	if searchEnv.Code != 20000 {
		t.Fatalf("search failed: code=%d msg=%s", searchEnv.Code, searchEnv.Message)
	}
	if len(searchEnv.Data.List) == 0 {
		t.Skip("no sku found in search results")
	}
	item := searchEnv.Data.List[0]
	return item.SkuId, item.SpuId, item.Stock, item.SaleCount
}
