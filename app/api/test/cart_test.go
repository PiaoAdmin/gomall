package test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	perrors "github.com/PiaoAdmin/pmall/common/errs"
)

// ========== Cart Flow Test ==========

// TestCartFlow performs a complete cart workflow test
func TestCartFlow(t *testing.T) {
	baseURL := getTestServer(t)
	client := &http.Client{Timeout: 10000 * time.Second}

	suffix := time.Now().UnixNano()

	// 1) Create and login a test user
	_, _, token := createAndLoginTestUser(t, client, baseURL, suffix)
	t.Logf("✓ User created and logged in")

	authHeader := map[string]string{"Authorization": fmt.Sprintf("Bearer %s", token)}

	// 2) Create a test product to add to cart
	spuID, skuID := createTestProduct(t, client, baseURL, suffix)
	t.Logf("✓ Test product created - SPU ID: %d, SKU ID: %d", spuID, skuID)

	// 3) Add product to cart
	addBody := map[string]any{
		"sku_id":   skuID,
		"quantity": 3,
	}
	addResp := postJSON[map[string]any](t, client, baseURL+"/cart/add", addBody, authHeader)
	if addResp.Code != uint64(perrors.Success.Code) {
		t.Fatalf("❌ Add to cart failed: code=%d msg=%s", addResp.Code, addResp.Message)
	}
	t.Logf("✓ Added to cart successfully")
	printCartResponse(t, "Add to cart response", addResp.Data)

	// 4) Get cart details
	cartResp := getJSON[map[string]any](t, client, baseURL+"/cart", authHeader)
	if cartResp.Code != uint64(perrors.Success.Code) {
		t.Fatalf("❌ Get cart failed: code=%d msg=%s", cartResp.Code, cartResp.Message)
	}
	t.Logf("✓ Cart retrieved successfully")
	printCartDetailsResponse(t, "Cart details", cartResp.Data)

	// Verify the cart has the product we added
	items := cartResp.Data["items"].([]any)
	if len(items) == 0 {
		t.Fatalf("❌ Cart is empty after adding product")
	}
	firstItem := items[0].(map[string]any)
	itemSkuID := uint64(firstItem["sku_id"].(float64))
	itemQuantity := uint64(firstItem["quantity"].(float64))
	if itemSkuID != skuID {
		t.Fatalf("❌ Expected SKU ID %d, got %d", skuID, itemSkuID)
	}
	if itemQuantity != 3 {
		t.Fatalf("❌ Expected quantity 3, got %d", itemQuantity)
	}
	t.Logf("✓ Cart contains correct product: SKU %d, Quantity %d", itemSkuID, itemQuantity)

	// 5) Add more quantity to the same product
	addBody2 := map[string]any{
		"sku_id":   skuID,
		"quantity": 2,
	}
	addResp2 := postJSON[map[string]any](t, client, baseURL+"/cart/add", addBody2, authHeader)
	if addResp2.Code != uint64(perrors.Success.Code) {
		t.Fatalf("❌ Add more quantity failed: code=%d msg=%s", addResp2.Code, addResp2.Message)
	}
	t.Logf("✓ Added more quantity successfully")

	// 6) Verify quantity increased
	cartResp2 := getJSON[map[string]any](t, client, baseURL+"/cart", authHeader)
	if cartResp2.Code != uint64(perrors.Success.Code) {
		t.Fatalf("❌ Get cart failed: code=%d msg=%s", cartResp2.Code, cartResp2.Message)
	}
	items2 := cartResp2.Data["items"].([]any)
	if len(items2) > 0 {
		item := items2[0].(map[string]any)
		newQuantity := uint64(item["quantity"].(float64))
		if newQuantity != 5 { // 3 + 2
			t.Fatalf("❌ Expected quantity 5, got %d", newQuantity)
		}
		t.Logf("✓ Quantity increased correctly to %d", newQuantity)
		printCartDetailsResponse(t, "Updated cart", cartResp2.Data)
	}

	// 7) Remove some quantity from cart
	removeBody := map[string]any{
		"sku_ids": []uint64{skuID, skuID},
	}
	removeResp := postJSON[map[string]any](t, client, baseURL+"/cart/remove", removeBody, authHeader)
	if removeResp.Code != uint64(perrors.Success.Code) {
		t.Fatalf("❌ Remove from cart failed: code=%d msg=%s", removeResp.Code, removeResp.Message)
	}
	t.Logf("✓ Removed from cart successfully")
	printCartResponse(t, "Remove response", removeResp.Data)

	// 8) Verify quantity decreased
	cartResp3 := getJSON[map[string]any](t, client, baseURL+"/cart", authHeader)
	if cartResp3.Code != uint64(perrors.Success.Code) {
		t.Fatalf("❌ Get cart failed: code=%d msg=%s", cartResp3.Code, cartResp3.Message)
	}
	items3 := cartResp3.Data["items"].([]any)
	if len(items3) > 0 {
		item := items3[0].(map[string]any)
		newQuantity := uint64(item["quantity"].(float64))
		if newQuantity != 3 { // 5 - 2
			t.Fatalf("❌ Expected quantity 3, got %d", newQuantity)
		}
		t.Logf("✓ Quantity decreased correctly to %d", newQuantity)
	}

	// 9) Clear cart
	clearResp := postJSON[map[string]any](t, client, baseURL+"/cart/clear", map[string]any{}, authHeader)
	if clearResp.Code != uint64(perrors.Success.Code) {
		t.Fatalf("❌ Clear cart failed: code=%d msg=%s", clearResp.Code, clearResp.Message)
	}
	t.Logf("✓ Cart cleared successfully")

	// 10) Verify cart is empty
	cartResp4 := getJSON[map[string]any](t, client, baseURL+"/cart", authHeader)
	if cartResp4.Code != uint64(perrors.Success.Code) {
		t.Fatalf("❌ Get cart failed: code=%d msg=%s", cartResp4.Code, cartResp.Message)
	}

	if cartResp4.Data["total_amount"] != "0.00" {
		t.Fatalf("❌ Cart should be empty after clearing")
	}
	t.Logf("✓ Cart is empty as expected")
	printCartDetailsResponse(t, "Empty cart", cartResp4.Data)

	t.Logf("✅ Complete cart flow test passed!")
}

// TestCartMultipleProducts tests adding multiple different products to cart
func TestCartMultipleProducts(t *testing.T) {
	baseURL := getTestServer(t)
	client := &http.Client{Timeout: 5 * time.Second}

	suffix := time.Now().UnixNano()

	// Create and login user
	username, _, token := createAndLoginTestUser(t, client, baseURL, suffix)
	t.Logf("✓ User logged in: %s", username)

	authHeader := map[string]string{"Authorization": fmt.Sprintf("Bearer %s", token)}

	// Create multiple test products
	_, skuID1 := createTestProduct(t, client, baseURL, suffix+1)
	_, skuID2 := createTestProduct(t, client, baseURL, suffix+2)
	_, skuID3 := createTestProduct(t, client, baseURL, suffix+3)
	t.Logf("✓ Created 3 test products: SKU1=%d, SKU2=%d, SKU3=%d", skuID1, skuID2, skuID3)

	// Add first product
	addResp1 := postJSON[map[string]any](t, client, baseURL+"/cart/add", map[string]any{
		"sku_id":   skuID1,
		"quantity": 2,
	}, authHeader)
	if addResp1.Code != uint64(perrors.Success.Code) {
		t.Fatalf("❌ Add product 1 failed: code=%d msg=%s", addResp1.Code, addResp1.Message)
	}

	// Add second product
	addResp2 := postJSON[map[string]any](t, client, baseURL+"/cart/add", map[string]any{
		"sku_id":   skuID2,
		"quantity": 5,
	}, authHeader)
	if addResp2.Code != uint64(perrors.Success.Code) {
		t.Fatalf("❌ Add product 2 failed: code=%d msg=%s", addResp2.Code, addResp2.Message)
	}

	// Add third product
	addResp3 := postJSON[map[string]any](t, client, baseURL+"/cart/add", map[string]any{
		"sku_id":   skuID3,
		"quantity": 1,
	}, authHeader)
	if addResp3.Code != uint64(perrors.Success.Code) {
		t.Fatalf("❌ Add product 3 failed: code=%d msg=%s", addResp3.Code, addResp3.Message)
	}

	t.Logf("✓ Added 3 products to cart")

	// Get cart and verify all products
	cartResp := getJSON[map[string]any](t, client, baseURL+"/cart", authHeader)
	if cartResp.Code != uint64(perrors.Success.Code) {
		t.Fatalf("❌ Get cart failed: code=%d msg=%s", cartResp.Code, cartResp.Message)
	}

	items := cartResp.Data["items"].([]any)
	if len(items) != 3 {
		t.Fatalf("❌ Expected 3 items in cart, got %d", len(items))
	}

	t.Logf("✓ Cart contains 3 items")
	printCartDetailsResponse(t, "Multi-product cart", cartResp.Data)

	// Calculate and verify total
	totalAmount := cartResp.Data["total_amount"].(string)
	totalCount := uint64(cartResp.Data["total_quantity"].(float64))

	if totalCount != 8 { // 2 + 5 + 1
		t.Fatalf("❌ Expected total count 8, got %d", totalCount)
	}
	t.Logf("✓ Total count is correct: %d", totalCount)
	t.Logf("✓ Total amount: %s", totalAmount)

	// Remove one product completely
	removeResp := postJSON[map[string]any](t, client, baseURL+"/cart/remove", map[string]any{
		"sku_ids": []uint64{skuID2, skuID2, skuID2, skuID2, skuID2},
	}, authHeader)
	if removeResp.Code != uint64(perrors.Success.Code) {
		t.Fatalf("❌ Remove product failed: code=%d msg=%s", removeResp.Code, removeResp.Message)
	}

	// Verify only 2 products remain
	cartResp2 := getJSON[map[string]any](t, client, baseURL+"/cart", authHeader)
	if cartResp2.Code != uint64(perrors.Success.Code) {
		t.Fatalf("❌ Get cart failed: code=%d msg=%s", cartResp2.Code, cartResp2.Message)
	}

	items2 := cartResp2.Data["items"].([]any)
	if len(items2) != 2 {
		t.Fatalf("❌ Expected 2 items in cart after removal, got %d", len(items2))
	}
	t.Logf("✓ Product removed successfully, 2 items remaining")
	printCartDetailsResponse(t, "Cart after removal", cartResp2.Data)

	t.Logf("✅ Multiple products test passed!")
}

// TestCartWithoutAuth tests cart operations without authentication
func TestCartWithoutAuth(t *testing.T) {
	baseURL := getTestServer(t)
	client := &http.Client{Timeout: 5 * time.Second}

	suffix := time.Now().UnixNano()

	// Create test product
	_, skuID := createTestProduct(t, client, baseURL, suffix)

	// Try to add to cart without auth
	addResp := postJSON[any](t, client, baseURL+"/cart/add", map[string]any{
		"sku_id":   skuID,
		"quantity": 1,
	}, nil)

	if addResp.Code == uint64(perrors.Success.Code) {
		t.Fatalf("❌ Expected add to cart to fail without auth, but succeeded")
	}
	t.Logf("✓ Add to cart without auth failed as expected: code=%d msg=%s", addResp.Code, addResp.Message)

	// Try to get cart without auth
	cartResp := getJSON[any](t, client, baseURL+"/cart", nil)
	if cartResp.Code == uint64(perrors.Success.Code) {
		t.Fatalf("❌ Expected get cart to fail without auth, but succeeded")
	}
	t.Logf("✓ Get cart without auth failed as expected: code=%d msg=%s", cartResp.Code, cartResp.Message)

	t.Logf("✅ Cart without auth test passed!")
}

// TestCartEmptyCart tests getting an empty cart
func TestCartEmptyCart(t *testing.T) {
	baseURL := getTestServer(t)
	client := &http.Client{Timeout: 5 * time.Second}

	suffix := time.Now().UnixNano()

	// Create and login user
	_, _, token := createAndLoginTestUser(t, client, baseURL, suffix)
	authHeader := map[string]string{"Authorization": fmt.Sprintf("Bearer %s", token)}

	// Get cart (should be empty)
	cartResp := getJSON[map[string]any](t, client, baseURL+"/cart", authHeader)
	if cartResp.Code != uint64(perrors.Success.Code) {
		t.Fatalf("❌ Get empty cart failed: code=%d msg=%s", cartResp.Code, cartResp.Message)
	}

	totalAmount := cartResp.Data["total_amount"].(string)

	if totalAmount != "0.00" {
		t.Fatalf("❌ Empty cart total amount should be 0, got %s", totalAmount)
	}

	t.Logf("✓ Empty cart returned correctly")
	printCartDetailsResponse(t, "Empty cart", cartResp.Data)

	t.Logf("✅ Empty cart test passed!")
}

// Helper functions

// printCartResponse prints key fields from add/remove cart response
func printCartResponse(t *testing.T, title string, data map[string]any) {
	t.Helper()
	t.Logf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	t.Logf("📦 %s", title)
	t.Logf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	if msg, ok := data["message"].(string); ok {
		t.Logf("  Message: %s", msg)
	}
	t.Logf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}

// printCartDetailsResponse prints detailed cart information
func printCartDetailsResponse(t *testing.T, title string, data map[string]any) {
	t.Helper()
	t.Logf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	t.Logf("🛒 %s", title)
	t.Logf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	if items, ok := data["items"].([]any); ok {
		t.Logf("  Items Count: %d", len(items))
		for i, item := range items {
			itemMap := item.(map[string]any)
			t.Logf("  ┌─ Item %d ─────────────────────────", i+1)
			if skuID, ok := itemMap["sku_id"].(float64); ok {
				t.Logf("  │ SKU ID:       %d", uint64(skuID))
			}
			if name, ok := itemMap["name"].(string); ok {
				t.Logf("  │ Name:         %s", name)
			}
			if quantity, ok := itemMap["quantity"].(float64); ok {
				t.Logf("  │ Quantity:     %d", uint64(quantity))
			}
			if price, ok := itemMap["price"].(string); ok {
				t.Logf("  │ Price:        %s", price)
			}
			if subtotal, ok := itemMap["subtotal"].(string); ok {
				t.Logf("  │ Subtotal:     %s", subtotal)
			}
			if mainImage, ok := itemMap["main_image"].(string); ok && mainImage != "" {
				t.Logf("  │ Image:        %s", mainImage)
			}
			if stock, ok := itemMap["stock"].(float64); ok {
				t.Logf("  │ Stock:        %d", uint64(stock))
			}
			t.Logf("  └─────────────────────────────────────")
		}
	}

	if totalAmount, ok := data["total_amount"].(string); ok {
		t.Logf("  💰 Total Amount: %s", totalAmount)
	}
	if totalCount, ok := data["total_count"].(float64); ok {
		t.Logf("  📊 Total Count:  %d", uint64(totalCount))
	}

	t.Logf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}
