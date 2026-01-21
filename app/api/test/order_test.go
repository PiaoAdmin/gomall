package test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	perrors "github.com/PiaoAdmin/pmall/common/errs"
)

// TestOrderFlow covers placing an order, listing it, and cancelling it
func TestOrderFlow(t *testing.T) {
	baseURL := getTestServer(t)
	client := &http.Client{Timeout: 10 * time.Second}

	suffix := time.Now().UnixNano()

	// create and login user
	_, _, token := createAndLoginTestUser(t, client, baseURL, suffix)
	authHeader := map[string]string{"Authorization": fmt.Sprintf("Bearer %s", token)}

	// create product
	_, skuID := createTestProduct(t, client, baseURL, suffix)
	t.Logf("✓ Created test product SKU: %d", skuID)
	// add to cart
	addCartBody := map[string]any{
		"sku_id":   skuID,
		"quantity": 2,
	}
	addCartResp := postJSON[map[string]any](t, client, baseURL+"/cart/add", addCartBody, authHeader)
	if addCartResp.Code != uint64(perrors.Success.Code) {
		t.Fatalf("Add to cart failed: code=%d msg=%s", addCartResp.Code, addCartResp.Message)
	}
	t.Logf("✓ Added to cart: %v", addCartResp.Data)
	// place order
	placeBody := map[string]any{
		"email": "buyer@example.com",
		"shipping_address": map[string]any{
			"name":           "Tester",
			"street_address": "123 Test St",
			"city":           "TestCity",
			"zip_code":       100000,
		},
	}

	placeResp := postJSON[map[string]any](t, client, baseURL+"/orders", placeBody, authHeader)
	if placeResp.Code != uint64(perrors.Success.Code) {
		t.Fatalf("Place order failed: code=%d msg=%s", placeResp.Code, placeResp.Message)
	}
	t.Logf("✓ Place order response: %v", placeResp.Data)

	orderMap, ok := placeResp.Data["order"].(map[string]any)
	if !ok {
		t.Fatal("unexpected place order response shape")
	}
	orderID, _ := orderMap["order_id"].(string)
	if orderID == "" {
		t.Fatal("empty order_id returned")
	}

	// list orders
	listResp := getJSON[map[string]any](t, client, baseURL+"/orders", authHeader)
	if listResp.Code != uint64(perrors.Success.Code) {
		t.Fatalf("List orders failed: code=%d msg=%s", listResp.Code, listResp.Message)
	}
	t.Logf("✓ List orders: %v", listResp.Data)

	// cancel order
	cancelURL := fmt.Sprintf("%s/orders/%s/cancel", baseURL, orderID)
	cancelResp := postJSON[map[string]any](t, client, cancelURL, map[string]any{}, authHeader)
	if cancelResp.Code != uint64(perrors.Success.Code) {
		t.Fatalf("Cancel order failed: code=%d msg=%s", cancelResp.Code, cancelResp.Message)
	}
	t.Logf("✓ Cancel order success: %v", cancelResp.Data)
}
