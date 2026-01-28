package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	perrors "github.com/PiaoAdmin/pmall/common/errs"
)

// respEnvelope is a small helper to decode the unified response wrapper.
type respEnvelope[T any] struct {
	Code    uint64 `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

type loginData struct {
	Token string `json:"token"`
}

// getTestServer returns the test server base URL
func getTestServer(t *testing.T) string {
	return testBaseURL
}

// doJSON is a generic helper for making HTTP requests with JSON body and response
func doJSON[T any](t *testing.T, client *http.Client, method, url string, body any, headers map[string]string) respEnvelope[T] {
	t.Helper()

	var bodyReader *bytes.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("marshal body: %v", err)
		}
		bodyReader = bytes.NewReader(b)
	}

	var req *http.Request
	var err error
	if bodyReader != nil {
		req, err = http.NewRequest(method, url, bodyReader)
	} else {
		req, err = http.NewRequest(method, url, nil)
	}
	if err != nil {
		t.Fatalf("new request: %v", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("do request: %v", err)
	}
	defer resp.Body.Close()

	var env respEnvelope[T]
	if err := json.NewDecoder(resp.Body).Decode(&env); err != nil {
		t.Fatalf("decode resp: %v", err)
	}
	return env
}

// postJSON performs a POST request with JSON body
func postJSON[T any](t *testing.T, client *http.Client, url string, body any, headers map[string]string) respEnvelope[T] {
	return doJSON[T](t, client, http.MethodPost, url, body, headers)
}

// getJSON performs a GET request
func getJSON[T any](t *testing.T, client *http.Client, url string, headers map[string]string) respEnvelope[T] {
	return doJSON[T](t, client, http.MethodGet, url, nil, headers)
}

// createTestUser registers a new user and returns username and password
func createTestUser(t *testing.T, client *http.Client, baseURL string, suffix int64) (username, password string) {
	t.Helper()
	username = fmt.Sprintf("tester_%d", suffix%1000000)
	password = "Passw0rd!"

	regBody := map[string]any{
		"username":         username,
		"password":         password,
		"password_confirm": password,
		"email":            fmt.Sprintf("%s@example.com", username),
	}
	regResp := postJSON[map[string]any](t, client, baseURL+"/register", regBody, nil)
	if regResp.Code != uint64(perrors.Success.Code) {
		t.Fatalf("registration failed: code=%d msg=%s", regResp.Code, regResp.Message)
	}
	return username, password
}

// loginTestUser logs in with username and password, returns token
func loginTestUser(t *testing.T, client *http.Client, baseURL, username, password string) string {
	t.Helper()
	loginBody := map[string]any{
		"username": username,
		"password": password,
	}
	loginResp := postJSON[loginData](t, client, baseURL+"/login", loginBody, nil)
	if loginResp.Code != uint64(perrors.Success.Code) {
		t.Fatalf("login failed: code=%d msg=%s", loginResp.Code, loginResp.Message)
	}
	if loginResp.Data.Token == "" {
		t.Fatalf("login token is empty")
	}
	return loginResp.Data.Token
}

// createAndLoginTestUser registers and logs in a new user, returns username, password and token
func createAndLoginTestUser(t *testing.T, client *http.Client, baseURL string, suffix int64) (username, password, token string) {
	t.Helper()
	username, password = createTestUser(t, client, baseURL, suffix)
	token = loginTestUser(t, client, baseURL, username, password)
	return username, password, token
}

// createTestProduct creates a test product and returns SPU ID and SKU ID
func createTestProduct(t *testing.T, client *http.Client, baseURL string, suffix int64) (spuID, skuID uint64) {
	t.Helper()

	createBody := map[string]any{
		"spu": map[string]any{
			"brand_id":     1,
			"category_id":  1,
			"name":         fmt.Sprintf("Test Product %d", suffix%1000000),
			"sub_title":    "Test product for testing",
			"main_image":   "https://example.com/test-product.jpg",
			"sort":         100,
			"service_bits": 0,
		},
		"skus": []map[string]any{
			{
				"sku_code":      fmt.Sprintf("TEST-SKU-%d", suffix%1000000),
				"name":          "Test SKU",
				"sub_title":     "Test SKU for testing",
				"main_image":    "https://example.com/test-sku.jpg",
				"price":         "9999",
				"market_price":  "12999",
				"stock":         1000,
				"sku_spec_data": `{"color":"blue","size":"M"}`,
			},
		},
		"detail": map[string]any{
			"description":     "Test product for automated testing",
			"images":          []string{"https://example.com/detail1.jpg"},
			"videos":          []string{},
			"market_tag_json": `{"test":true}`,
			"tech_tag_json":   `{"automated":true}`,
		},
	}

	createResp := postJSON[map[string]any](t, client, baseURL+"/admin/products", createBody, nil)
	if createResp.Code != uint64(perrors.Success.Code) {
		t.Fatalf("create test product failed: code=%d msg=%s", createResp.Code, createResp.Message)
	}

	spuID = uint64(createResp.Data["spu_id"].(float64))

	// Get product detail to retrieve SKU ID
	detailResp := getJSON[map[string]any](t, client, fmt.Sprintf("%s/products/%d", baseURL, spuID), nil)
	if detailResp.Code != uint64(perrors.Success.Code) {
		t.Fatalf("get product detail failed: code=%d msg=%s", detailResp.Code, detailResp.Message)
	}

	product := detailResp.Data["product"].(map[string]any)
	skus := product["skus"].([]any)
	if len(skus) == 0 {
		t.Fatalf("created product has no SKUs")
	}

	sku := skus[0].(map[string]any)
	skuID = uint64(sku["id"].(float64))

	return spuID, skuID
}

type productDetailResp struct {
	Product struct {
		SaleCount int32 `json:"sale_count"`
		Skus      []struct {
			Id    uint64 `json:"id"`
			Stock int32  `json:"stock"`
		} `json:"skus"`
	} `json:"product"`
}

func getProductStockAndSales(t *testing.T, client *http.Client, baseURL string, spuID, skuID uint64) (int32, int32) {
	t.Helper()
	detailURL := fmt.Sprintf("%s/products/%d", baseURL, spuID)
	detailEnv := getJSON[productDetailResp](t, client, detailURL, nil)
	if detailEnv.Code != uint64(perrors.Success.Code) {
		t.Fatalf("product detail failed: code=%d msg=%s", detailEnv.Code, detailEnv.Message)
	}
	stock := int32(-1)
	for _, sku := range detailEnv.Data.Product.Skus {
		if sku.Id == skuID {
			stock = sku.Stock
			break
		}
	}
	if stock < 0 {
		t.Fatalf("sku not found in product detail: sku_id=%d", skuID)
	}
	return stock, detailEnv.Data.Product.SaleCount
}
