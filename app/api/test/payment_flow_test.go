package test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	perrors "github.com/PiaoAdmin/pmall/common/errs"
)

func TestOrderAndPayFlow(t *testing.T) {
	baseURL := getTestServer(t)
	client := &http.Client{Timeout: 10 * time.Second}

	suffix := time.Now().UnixNano()

	// 创建并登录测试用户
	_, _, token := createAndLoginTestUser(t, client, baseURL, suffix)
	authHeader := map[string]string{"Authorization": fmt.Sprintf("Bearer %s", token)}

	// 创建商品用于下单
	spuID, skuID := createTestProduct(t, client, baseURL, suffix)
	stockBefore, saleBefore := getProductStockAndSales(t, client, baseURL, spuID, skuID)

	addCartBody := map[string]any{
		"sku_id":   skuID,
		"quantity": 2,
	}
	// 加入购物车
	addCartResp := postJSON[map[string]any](t, client, baseURL+"/cart/add", addCartBody, authHeader)
	if addCartResp.Code != uint64(perrors.Success.Code) {
		t.Fatalf("add to cart failed: code=%d msg=%s", addCartResp.Code, addCartResp.Message)
	}

	placeBody := map[string]any{
		"email": "buyer@example.com",
		"shipping_address": map[string]any{
			"name":           "Tester",
			"street_address": "123 Test St",
			"city":           "TestCity",
			"zip_code":       100000,
		},
	}
	// 先下单，生成订单
	placeResp := postJSON[map[string]any](t, client, baseURL+"/orders", placeBody, authHeader)
	if placeResp.Code != uint64(perrors.Success.Code) {
		t.Fatalf("place order failed: code=%d msg=%s", placeResp.Code, placeResp.Message)
	}
	orderMap, ok := placeResp.Data["order"].(map[string]any)
	if !ok {
		t.Fatal("unexpected place order response shape")
	}
	orderID, _ := orderMap["order_id"].(string)
	if orderID == "" {
		t.Fatal("empty order_id returned")
	}

	// 下单后库存锁定、销量增加
	stockAfter, saleAfter := getProductStockAndSales(t, client, baseURL, spuID, skuID)
	if stockAfter != stockBefore-2 {
		t.Fatalf("stock not locked on place order: before=%d after=%d", stockBefore, stockAfter)
	}
	if saleAfter != saleBefore+2 {
		t.Fatalf("sale_count not increased on place order: before=%d after=%d", saleBefore, saleAfter)
	}

	payBody := map[string]any{
		"order_id":    orderID,
		"amount":      "1.00",
		"credit_card": validCreditCard,
	}
	// 使用支付接口完成支付
	payResp := postJSON[map[string]any](t, client, baseURL+"/payment/pay", payBody, authHeader)
	if payResp.Code != uint64(perrors.Success.Code) {
		t.Fatalf("pay failed: code=%d msg=%s", payResp.Code, payResp.Message)
	}
	if tradeNo, ok := payResp.Data["trade_no"].(string); !ok || tradeNo == "" {
		t.Fatalf("pay trade_no empty")
	}
}

func TestPaymentInvalidCard(t *testing.T) {
	baseURL := getTestServer(t)
	client := &http.Client{Timeout: 10 * time.Second}

	suffix := time.Now().UnixNano()

	// 创建并登录测试用户
	_, _, token := createAndLoginTestUser(t, client, baseURL, suffix)
	authHeader := map[string]string{"Authorization": fmt.Sprintf("Bearer %s", token)}

	// 创建商品用于下单
	_, skuID := createTestProduct(t, client, baseURL, suffix)

	addCartBody := map[string]any{
		"sku_id":   skuID,
		"quantity": 1,
	}
	// 加入购物车
	addCartResp := postJSON[map[string]any](t, client, baseURL+"/cart/add", addCartBody, authHeader)
	if addCartResp.Code != uint64(perrors.Success.Code) {
		t.Fatalf("add to cart failed: code=%d msg=%s", addCartResp.Code, addCartResp.Message)
	}

	placeBody := map[string]any{
		"email": "buyer@example.com",
		"shipping_address": map[string]any{
			"name":           "Tester",
			"street_address": "123 Test St",
			"city":           "TestCity",
			"zip_code":       100000,
		},
	}
	// 先下单
	placeResp := postJSON[map[string]any](t, client, baseURL+"/orders", placeBody, authHeader)
	if placeResp.Code != uint64(perrors.Success.Code) {
		t.Fatalf("place order failed: code=%d msg=%s", placeResp.Code, placeResp.Message)
	}
	orderMap, ok := placeResp.Data["order"].(map[string]any)
	if !ok {
		t.Fatal("unexpected place order response shape")
	}
	orderID, _ := orderMap["order_id"].(string)
	if orderID == "" {
		t.Fatal("empty order_id returned")
	}

	payBody := map[string]any{
		"order_id":    orderID,
		"amount":      "1.00",
		"credit_card": "123",
	}
	// 使用非法信用卡支付，期望失败
	payResp := postJSON[map[string]any](t, client, baseURL+"/payment/pay", payBody, authHeader)
	if payResp.Code != uint64(perrors.ErrParam.Code) {
		t.Fatalf("expected invalid card code=%d, got=%d msg=%s", perrors.ErrParam.Code, payResp.Code, payResp.Message)
	}
}
