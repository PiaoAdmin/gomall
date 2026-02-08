package test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	perrors "github.com/PiaoAdmin/pmall/common/errs"
)

// TestOrderAutoCancel 测试订单超时自动取消功能
// 流程:
// 1. 创建用户并登录
// 2. 创建商品
// 3. 添加到购物车
// 4. 下单（此时会发送延迟取消消息）
// 5. 等待延迟时间（当前设置30秒）
// 6. 查询订单状态，验证是否自动取消
func TestOrderAutoCancel(t *testing.T) {
	baseURL := getTestServer(t)
	client := &http.Client{Timeout: 60 * time.Second}

	suffix := time.Now().UnixNano()

	// 延迟时间设置（与 order 服务的 OrderCancelDelayTime 保持一致）
	delayTime := 30 * time.Second

	t.Log("=== 订单超时自动取消测试 ===")
	t.Logf("延迟时间: %v", delayTime)

	// Step 1: 创建用户并登录
	_, _, token := createAndLoginTestUser(t, client, baseURL, suffix)
	authHeader := map[string]string{"Authorization": fmt.Sprintf("Bearer %s", token)}
	t.Log("✓ 用户创建并登录成功")

	// Step 2: 创建商品
	_, skuID := createTestProduct(t, client, baseURL, suffix)
	t.Logf("✓ 商品创建成功, SKU ID: %d", skuID)

	// Step 3: 添加到购物车
	addCartBody := map[string]any{
		"sku_id":   skuID,
		"quantity": 1,
	}
	addCartResp := postJSON[map[string]any](t, client, baseURL+"/cart/add", addCartBody, authHeader)
	if addCartResp.Code != uint64(perrors.Success.Code) {
		t.Fatalf("添加购物车失败: code=%d msg=%s", addCartResp.Code, addCartResp.Message)
	}
	t.Log("✓ 商品已添加到购物车")

	// Step 4: 下单
	placeBody := map[string]any{
		"email": "autocancel@example.com",
		"shipping_address": map[string]any{
			"name":           "Auto Cancel Tester",
			"street_address": "123 Test St",
			"city":           "TestCity",
			"zip_code":       100000,
		},
	}

	placeTime := time.Now()
	placeResp := postJSON[map[string]any](t, client, baseURL+"/orders", placeBody, authHeader)
	if placeResp.Code != uint64(perrors.Success.Code) {
		t.Fatalf("下单失败: code=%d msg=%s", placeResp.Code, placeResp.Message)
	}

	orderMap, ok := placeResp.Data["order"].(map[string]any)
	if !ok {
		t.Fatal("下单响应格式错误")
	}
	orderID, _ := orderMap["order_id"].(string)
	if orderID == "" {
		t.Fatal("订单 ID 为空")
	}

	t.Logf("✓ 订单创建成功")
	t.Logf("  订单 ID: %s", orderID)
	t.Logf("  创建时间: %s", placeTime.Format("15:04:05"))
	t.Logf("  预计自动取消时间: %s", placeTime.Add(delayTime).Format("15:04:05"))

	// Step 5: 立即查询订单状态（应该是 placed）
	t.Log("\n立即查询订单状态...")
	initialOrder := getOrderDetail(t, client, baseURL, orderID, authHeader)
	initialStatus := getOrderStatus(initialOrder)
	t.Logf("  当前状态: %s", initialStatus)

	if initialStatus != "placed" {
		t.Logf("  警告: 预期状态为 'placed', 实际为 '%s'", initialStatus)
	}

	// Step 6: 等待延迟时间 + 缓冲
	waitTime := delayTime + 10*time.Second
	t.Logf("\n等待 %v 让订单自动取消...", waitTime)

	// 进度显示
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	done := time.After(waitTime)

	elapsed := 10 * time.Second
waitLoop:
	for {
		select {
		case <-ticker.C:
			t.Logf("  已等待 %v...", elapsed)
			elapsed += 10 * time.Second
		case <-done:
			break waitLoop
		}
	}

	t.Log("等待完成，查询订单最终状态...")

	// Step 7: 再次查询订单状态（应该是 canceled）
	finalOrder := getOrderDetail(t, client, baseURL, orderID, authHeader)
	finalStatus := getOrderStatus(finalOrder)

	t.Logf("\n=== 测试结果 ===")
	t.Logf("订单 ID: %s", orderID)
	t.Logf("初始状态: %s", initialStatus)
	t.Logf("最终状态: %s", finalStatus)
	t.Logf("总耗时: %v", time.Since(placeTime))

	if finalStatus == "canceled" {
		t.Log("\n✓✓✓ 测试通过! 订单已自动取消 ✓✓✓")
	} else {
		t.Errorf("\n✗ 测试失败! 预期状态 'canceled', 实际状态 '%s'", finalStatus)
		t.Log("可能原因:")
		t.Log("1. 取消消费者未启动")
		t.Log("2. 延迟队列配置错误")
		t.Log("3. 订单状态更新失败")
	}
}

// getOrderDetail 获取订单详情
func getOrderDetail(t *testing.T, client *http.Client, baseURL, orderID string, headers map[string]string) map[string]any {
	t.Helper()
	// 通过列表获取订单
	listResp := getJSON[map[string]any](t, client, baseURL+"/orders", headers)
	if listResp.Code != uint64(perrors.Success.Code) {
		t.Logf("获取订单列表失败: code=%d msg=%s", listResp.Code, listResp.Message)
		return nil
	}

	orders, ok := listResp.Data["orders"].([]any)
	if !ok {
		t.Log("订单列表格式错误")
		return nil
	}

	for _, o := range orders {
		order, ok := o.(map[string]any)
		if !ok {
			continue
		}
		if id, _ := order["order_id"].(string); id == orderID {
			return order
		}
	}

	t.Logf("未找到订单: %s", orderID)
	return nil
}

// getOrderStatus 从订单数据中获取状态
func getOrderStatus(order map[string]any) string {
	if order == nil {
		return "unknown"
	}
	status, _ := order["status"].(string)
	if status == "" {
		// 尝试从 order_state 获取
		state, _ := order["order_state"].(string)
		if state != "" {
			return state
		}
		return "unknown"
	}
	return status
}

// TestOrderAutoCancelQuick 快速版本 - 用于手动验证
// 只下单不等待，需要手动查看日志确认取消
func TestOrderAutoCancelQuick(t *testing.T) {
	baseURL := getTestServer(t)
	client := &http.Client{Timeout: 10 * time.Second}

	suffix := time.Now().UnixNano()

	t.Log("=== 快速下单测试（验证延迟消息发送）===")

	// 创建用户并登录
	_, _, token := createAndLoginTestUser(t, client, baseURL, suffix)
	authHeader := map[string]string{"Authorization": fmt.Sprintf("Bearer %s", token)}
	t.Log("✓ 用户就绪")

	// 创建商品并添加到购物车
	_, skuID := createTestProduct(t, client, baseURL, suffix)
	addCartBody := map[string]any{"sku_id": skuID, "quantity": 1}
	postJSON[map[string]any](t, client, baseURL+"/cart/add", addCartBody, authHeader)
	t.Log("✓ 购物车就绪")

	// 下单
	placeBody := map[string]any{
		"email": "quicktest@example.com",
		"shipping_address": map[string]any{
			"name":           "Quick Tester",
			"street_address": "123 Quick St",
			"city":           "QuickCity",
			"zip_code":       100000,
		},
	}

	placeResp := postJSON[map[string]any](t, client, baseURL+"/orders", placeBody, authHeader)
	if placeResp.Code != uint64(perrors.Success.Code) {
		t.Fatalf("下单失败: %s", placeResp.Message)
	}

	orderMap := placeResp.Data["order"].(map[string]any)
	orderID := orderMap["order_id"].(string)

	t.Logf("✓ 订单已创建: %s", orderID)
	t.Logf("  时间: %s", time.Now().Format("15:04:05"))
	t.Log("\n请检查:")
	t.Log("1. order 服务日志 - 应有 'Order cancel delay message published' 日志")
	t.Log("2. RabbitMQ 管理界面 - order_delay_queue 应有消息")
	t.Log("3. 30秒后 - order_cancel_queue 应收到消息")
	t.Log("4. order 服务日志 - 应有取消处理日志")
}
