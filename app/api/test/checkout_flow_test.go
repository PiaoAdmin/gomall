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
		SkuId uint64 `json:"sku_id"`
	} `json:"list"`
	Total int64 `json:"total"`
}

type checkoutResp struct {
	OrderId string `json:"order_id"`
}

func TestCheckoutFlow(t *testing.T) {
	client := &http.Client{Timeout: 10 * time.Second}
	base := testBaseURL

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
	authHeader := map[string]string{"Authorization": fmt.Sprintf("Bearer %s", loginEnv.Data.Token)}

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
	skuID := searchEnv.Data.List[0].SkuId

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
	}
	checkoutEnv := postJSONEnvelope[checkoutResp](t, client, base+"/checkout", checkoutBody, authHeader)
	if checkoutEnv.Code != 20000 {
		t.Fatalf("checkout failed: code=%d msg=%s", checkoutEnv.Code, checkoutEnv.Message)
	}
	if checkoutEnv.Data.OrderId == "" {
		t.Fatalf("checkout order_id empty")
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
