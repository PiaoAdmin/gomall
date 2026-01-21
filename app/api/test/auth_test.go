package test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	perrors "github.com/PiaoAdmin/pmall/common/errs"
)

// TestAuthFlow performs a real HTTP flow: register -> login -> get user info -> refresh.
func TestAuthFlow(t *testing.T) {
	baseURL := getTestServer(t)

	client := &http.Client{Timeout: 5 * time.Second}

	// Use unique username to avoid conflicts between runs.
	suffix := time.Now().UnixNano()
	username := fmt.Sprintf("AuthFlow_%d", suffix%10000)
	password := "Passw0rd!"

	// 1) Register
	regBody := map[string]any{
		"username":         username,
		"password":         password,
		"password_confirm": password,
		"email":            fmt.Sprintf("%s@example.com", username),
	}
	regResp := postJSON[map[string]any](t, client, baseURL+"/register", regBody, nil)
	if regResp.Code != uint64(perrors.Success.Code) {
		t.Fatalf("register failed: code=%d msg=%s data=%v", regResp.Code, regResp.Message, regResp.Data)
	}

	// 2) Login
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

	authHeader := map[string]string{"Authorization": fmt.Sprintf("Bearer %s", loginResp.Data.Token)}

	// 3) Get user info
	infoResp := getJSON[map[string]any](t, client, baseURL+"/auth/info", authHeader)
	if infoResp.Code != uint64(perrors.Success.Code) {
		t.Fatalf("get user info failed: code=%d msg=%s", infoResp.Code, infoResp.Message)
	}

	// 4) Logout
	logoutResp := postJSON[map[string]any](t, client, baseURL+"/logout", map[string]any{}, authHeader)
	if logoutResp.Code != uint64(perrors.Success.Code) {
		t.Fatalf("logout failed: code=%d msg=%s", logoutResp.Code, logoutResp.Message)
	}

	// 5) Refresh token (should fail since logged out)
	refreshResp := postJSON[any](t, client, baseURL+"/refresh", map[string]any{}, authHeader)
	if refreshResp.Code == uint64(perrors.Success.Code) {
		t.Fatalf("expected refresh after logout to fail, but succeeded")
	}
	if refreshResp.Code != uint64(perrors.ErrNotLogin.Code) {
		t.Logf("expected error code %d, got %d, msg: %s", perrors.ErrNotLogin.Code, refreshResp.Code, refreshResp.Message)
	}

	// 6) Access protected resource after logout (should fail)
	infoResp2 := getJSON[any](t, client, baseURL+"/auth/info", authHeader)
	if infoResp2.Code == uint64(perrors.Success.Code) {
		t.Fatalf("expected access after logout to fail, but succeeded")
	}
	if infoResp2.Code != uint64(perrors.ErrAuthFailed.Code) {
		t.Logf("expected error code %d, got %d, msg: %s", perrors.ErrAuthFailed.Code, infoResp2.Code, infoResp2.Message)
	}
}

// TestRegisterPasswordMismatch tests registration with mismatched passwords
func TestRegisterPasswordMismatch(t *testing.T) {
	baseURL := getTestServer(t)

	client := &http.Client{Timeout: 5 * time.Second}

	suffix := time.Now().UnixNano()
	username := fmt.Sprintf("pwdtok_%d", suffix%1000000)

	regBody := map[string]any{
		"username":         username,
		"password":         "Passw0rd!",
		"password_confirm": "DifferentPass!",
		"email":            fmt.Sprintf("%s@example.com", username),
	}
	regResp := postJSON[any](t, client, baseURL+"/register", regBody, nil)

	// Should fail with ErrParam
	if regResp.Code == uint64(perrors.Success.Code) {
		t.Fatalf("expected registration to fail with password mismatch, but succeeded")
	}
	if regResp.Code != uint64(perrors.ErrParam.Code) {
		t.Logf("expected error code %d, got %d, msg: %s", perrors.ErrParam.Code, regResp.Code, regResp.Message)
	}
}

// TestRegisterDuplicateUser tests registration with existing username
func TestRegisterDuplicateUser(t *testing.T) {
	baseURL := getTestServer(t)

	client := &http.Client{Timeout: 5 * time.Second}

	suffix := time.Now().UnixNano()
	username := fmt.Sprintf("tester_%d", suffix%1000000)
	password := "Passw0rd!"

	// First registration
	regBody := map[string]any{
		"username":         username,
		"password":         password,
		"password_confirm": password,
		"email":            fmt.Sprintf("%s@example.com", username),
	}
	regResp := postJSON[map[string]any](t, client, baseURL+"/register", regBody, nil)
	if regResp.Code != uint64(perrors.Success.Code) {
		t.Fatalf("first registration failed: code=%d msg=%s", regResp.Code, regResp.Message)
	}

	// Second registration with same username
	regResp2 := postJSON[any](t, client, baseURL+"/register", regBody, nil)
	if regResp2.Code == uint64(perrors.Success.Code) {
		t.Fatalf("expected duplicate registration to fail, but succeeded")
	}
}

// TestLoginNonexistentUser tests login with non-existent username
func TestLoginNonexistentUser(t *testing.T) {
	baseURL := getTestServer(t)

	client := &http.Client{Timeout: 5 * time.Second}

	loginBody := map[string]any{
		"username": fmt.Sprintf("nonexist_%d", time.Now().UnixNano()%1000000),
		"password": "Passw0rd!",
	}
	loginResp := postJSON[any](t, client, baseURL+"/login", loginBody, nil)
	if loginResp.Code == uint64(perrors.Success.Code) {
		t.Fatalf("expected login to fail for nonexistent user, but succeeded")
	}
}

// TestLoginWrongPassword tests login with incorrect password
func TestLoginWrongPassword(t *testing.T) {
	baseURL := getTestServer(t)

	client := &http.Client{Timeout: 5 * time.Second}

	suffix := time.Now().UnixNano()
	username := fmt.Sprintf("tester_%d", suffix%1000000)
	password := "Passw0rd!"

	// Register user
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

	// Try login with wrong password
	loginBody := map[string]any{
		"username": username,
		"password": "WrongPassword!",
	}
	loginResp := postJSON[any](t, client, baseURL+"/login", loginBody, nil)

	if loginResp.Code == uint64(perrors.Success.Code) {
		t.Fatalf("expected login to fail with wrong password, but succeeded")
	}
}

// TestUnauthorizedAccess tests accessing protected endpoint without token
func TestUnauthorizedAccess(t *testing.T) {
	baseURL := getTestServer(t)

	client := &http.Client{Timeout: 5 * time.Second}

	// Try to access protected endpoint without token
	infoResp := getJSON[any](t, client, baseURL+"/auth/info", nil)

	if infoResp.Code == uint64(perrors.Success.Code) {
		t.Fatalf("expected unauthorized access to fail, but succeeded")
	}
	if infoResp.Code != uint64(perrors.ErrAuthFailed.Code) {
		t.Logf("expected error code %d, got %d, msg: %s", perrors.ErrAuthFailed.Code, infoResp.Code, infoResp.Message)
	}
}

// TestInvalidToken tests accessing protected endpoint with invalid token
func TestInvalidToken(t *testing.T) {
	baseURL := getTestServer(t)

	client := &http.Client{Timeout: 5 * time.Second}

	// Use invalid token
	authHeader := map[string]string{"Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOjIwMDgzNzE3NDE5NjkwOTI2MDgsImV4cCI6MTc2ODAwMTE0NCwib3JpZ19pYXQiOjE3Njc3ODUxNDR9.W9jB5DD8Yqnrbx5-V6DP7QAujh_He82U4pLDmJRuKak"}
	infoResp := getJSON[any](t, client, baseURL+"/auth/info", authHeader)

	if infoResp.Code == uint64(perrors.Success.Code) {
		t.Fatalf("expected invalid token access to fail, but succeeded")
	}
}

// TestUpdateUser tests updating user information
func TestUpdateUser(t *testing.T) {
	baseURL := getTestServer(t)

	client := &http.Client{Timeout: 5 * time.Second}

	suffix := time.Now().UnixNano()
	username := fmt.Sprintf("tester_%d", suffix%1000000)
	password := "Passw0rd!"

	// Register and login
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

	loginBody := map[string]any{
		"username": username,
		"password": password,
	}
	loginResp := postJSON[loginData](t, client, baseURL+"/login", loginBody, nil)
	if loginResp.Code != uint64(perrors.Success.Code) {
		t.Fatalf("login failed: code=%d msg=%s", loginResp.Code, loginResp.Message)
	}

	authHeader := map[string]string{"Authorization": fmt.Sprintf("Bearer %s", loginResp.Data.Token)}

	// Update user info
	updateBody := map[string]any{
		"email": fmt.Sprintf("updated_%s@example.com", username),
		"phone": "13800138000",
	}
	updateResp := postJSON[map[string]any](t, client, baseURL+"/auth/update", updateBody, authHeader)
	if updateResp.Code != uint64(perrors.Success.Code) {
		t.Fatalf("update user failed: code=%d msg=%s", updateResp.Code, updateResp.Message)
	}

	// Verify update by getting user info
	infoResp := getJSON[map[string]any](t, client, baseURL+"/auth/info", authHeader)
	if infoResp.Code != uint64(perrors.Success.Code) {
		t.Fatalf("get user info failed: code=%d msg=%s", infoResp.Code, infoResp.Message)
	}
}

// TestUpdatePassword tests password change functionality
func TestUpdatePassword(t *testing.T) {
	baseURL := getTestServer(t)

	client := &http.Client{Timeout: 5 * time.Second}

	suffix := time.Now().UnixNano()
	username := fmt.Sprintf("tester_%d", suffix%1000000)
	oldPassword := "OldPassw0rd!"
	newPassword := "NewPassw0rd!"

	// Register
	regBody := map[string]any{
		"username":         username,
		"password":         oldPassword,
		"password_confirm": oldPassword,
		"email":            fmt.Sprintf("%s@example.com", username),
	}
	regResp := postJSON[map[string]any](t, client, baseURL+"/register", regBody, nil)
	if regResp.Code != uint64(perrors.Success.Code) {
		t.Fatalf("registration failed: code=%d msg=%s", regResp.Code, regResp.Message)
	}

	// Login with old password
	loginBody := map[string]any{
		"username": username,
		"password": oldPassword,
	}
	loginResp := postJSON[loginData](t, client, baseURL+"/login", loginBody, nil)
	if loginResp.Code != uint64(perrors.Success.Code) {
		t.Fatalf("login failed: code=%d msg=%s", loginResp.Code, loginResp.Message)
	}

	authHeader := map[string]string{"Authorization": fmt.Sprintf("Bearer %s", loginResp.Data.Token)}

	// Update password
	updatePwdBody := map[string]any{
		"old_password": oldPassword,
		"new_password": newPassword,
	}
	updateResp := postJSON[map[string]any](t, client, baseURL+"/auth/password", updatePwdBody, authHeader)
	if updateResp.Code != uint64(perrors.Success.Code) {
		t.Fatalf("update password failed: code=%d msg=%s", updateResp.Code, updateResp.Message)
	}

	// Try login with old password (should fail)
	loginOldResp := postJSON[any](t, client, baseURL+"/login", loginBody, nil)
	if loginOldResp.Code == uint64(perrors.Success.Code) {
		t.Fatalf("expected login with old password to fail after password change, but succeeded")
	}

	// Login with new password (should succeed)
	loginNewBody := map[string]any{
		"username": username,
		"password": newPassword,
	}
	loginNewResp := postJSON[loginData](t, client, baseURL+"/login", loginNewBody, nil)
	if loginNewResp.Code != uint64(perrors.Success.Code) {
		t.Fatalf("login with new password failed: code=%d msg=%s", loginNewResp.Code, loginNewResp.Message)
	}
}

// TestAccessAfterLogout tests that token is invalid after logout
func TestAccessAfterLogout(t *testing.T) {
	baseURL := getTestServer(t)

	client := &http.Client{Timeout: 5 * time.Second}

	suffix := time.Now().UnixNano()
	username := fmt.Sprintf("tester_%d", suffix%1000000)
	password := "Passw0rd!"

	// Register and login
	regBody := map[string]any{
		"username":         username,
		"password":         password,
		"password_confirm": password,
		"email":            fmt.Sprintf("%s@example.com", username),
	}
	postJSON[map[string]any](t, client, baseURL+"/register", regBody, nil)

	loginBody := map[string]any{
		"username": username,
		"password": password,
	}
	loginResp := postJSON[loginData](t, client, baseURL+"/login", loginBody, nil)
	if loginResp.Code != uint64(perrors.Success.Code) {
		t.Fatalf("login failed: code=%d msg=%s", loginResp.Code, loginResp.Message)
	}

	authHeader := map[string]string{"Authorization": fmt.Sprintf("Bearer %s", loginResp.Data.Token)}

	// Logout
	logoutResp := postJSON[map[string]any](t, client, baseURL+"/logout", map[string]any{}, authHeader)
	if logoutResp.Code != uint64(perrors.Success.Code) {
		t.Fatalf("logout failed: code=%d msg=%s", logoutResp.Code, logoutResp.Message)
	}

	// Try to access protected resource with the same token (should fail)
	infoResp := getJSON[any](t, client, baseURL+"/auth/info", authHeader)
	if infoResp.Code == uint64(perrors.Success.Code) {
		t.Fatalf("expected access after logout to fail, but succeeded")
	}
}

// ========== Registration Edge Cases ==========

// TestRegisterEmptyUsername tests registration with empty username
func TestRegisterEmptyUsername(t *testing.T) {
	baseURL := getTestServer(t)

	client := &http.Client{Timeout: 5 * time.Second}

	regBody := map[string]any{
		"username":         "",
		"password":         "Passw0rd!",
		"password_confirm": "Passw0rd!",
		"email":            "test@example.com",
	}
	regResp := postJSON[any](t, client, baseURL+"/register", regBody, nil)

	if regResp.Code == uint64(perrors.Success.Code) {
		t.Fatalf("expected registration with empty username to fail, but succeeded")
	}
	if regResp.Code != uint64(perrors.ErrParam.Code) {
		t.Fatalf("expected error code %d, got %d", perrors.ErrParam.Code, regResp.Code)
	}
}

// TestRegisterEmptyPassword tests registration with empty password
func TestRegisterEmptyPassword(t *testing.T) {
	baseURL := getTestServer(t)

	client := &http.Client{Timeout: 5 * time.Second}

	suffix := time.Now().UnixNano()
	regBody := map[string]any{
		"username":         fmt.Sprintf("tester_%d", suffix%1000000),
		"password":         "",
		"password_confirm": "",
		"email":            fmt.Sprintf("test_%d@example.com", suffix),
	}
	regResp := postJSON[any](t, client, baseURL+"/register", regBody, nil)

	if regResp.Code == uint64(perrors.Success.Code) {
		t.Fatalf("expected registration with empty password to fail, but succeeded")
	}
	if regResp.Code != uint64(perrors.ErrParam.Code) {
		t.Fatalf("expected error code %d, got %d", perrors.ErrParam.Code, regResp.Code)
	}
}

// TestRegisterInvalidEmail tests registration with invalid email format
func TestRegisterInvalidEmail(t *testing.T) {
	baseURL := getTestServer(t)

	client := &http.Client{Timeout: 5 * time.Second}

	suffix := time.Now().UnixNano()
	regBody := map[string]any{
		"username":         fmt.Sprintf("tester_%d", suffix%1000000),
		"password":         "Passw0rd!",
		"password_confirm": "Passw0rd!",
		"email":            "invalid-email-format",
	}
	regResp := postJSON[any](t, client, baseURL+"/register", regBody, nil)

	if regResp.Code == uint64(perrors.Success.Code) {
		t.Fatalf("expected registration with invalid email to fail, but succeeded")
	}
	if regResp.Code != uint64(perrors.ErrParam.Code) {
		t.Fatalf("expected error code %d, got %d", perrors.ErrParam.Code, regResp.Code)
	}
}

// ========== Login Edge Cases ==========

// TestLoginEmptyUsername tests login with empty username
func TestLoginEmptyUsername(t *testing.T) {
	baseURL := getTestServer(t)

	client := &http.Client{Timeout: 5 * time.Second}

	loginBody := map[string]any{
		"username": "",
		"password": "Passw0rd!",
	}
	loginResp := postJSON[any](t, client, baseURL+"/login", loginBody, nil)

	if loginResp.Code == uint64(perrors.Success.Code) {
		t.Fatalf("expected login with empty username to fail, but succeeded")
	}
}

// TestLoginEmptyPassword tests login with empty password
func TestLoginEmptyPassword(t *testing.T) {
	baseURL := getTestServer(t)

	client := &http.Client{Timeout: 5 * time.Second}

	suffix := time.Now().UnixNano()
	username, _ := createTestUser(t, client, baseURL, suffix)

	loginBody := map[string]any{
		"username": username,
		"password": "",
	}
	loginResp := postJSON[any](t, client, baseURL+"/login", loginBody, nil)

	if loginResp.Code == uint64(perrors.Success.Code) {
		t.Fatalf("expected login with empty password to fail, but succeeded")
	}
}

// ========== Update User Edge Cases ==========

// TestUpdateUserInvalidEmail tests updating user with invalid email format
func TestUpdateUserInvalidEmail(t *testing.T) {
	baseURL := getTestServer(t)

	client := &http.Client{Timeout: 5 * time.Second}

	suffix := time.Now().UnixNano()
	_, _, token := createAndLoginTestUser(t, client, baseURL, suffix)
	authHeader := map[string]string{"Authorization": fmt.Sprintf("Bearer %s", token)}

	// Try to update with invalid email
	updateBody := map[string]any{
		"email": "invalid-email",
	}
	updateResp := postJSON[any](t, client, baseURL+"/auth/update", updateBody, authHeader)

	if updateResp.Code == uint64(perrors.Success.Code) {
		t.Fatalf("expected update with invalid email to fail, but succeeded")
	}
}

// TestUpdateUserWithoutAuth tests updating user without authentication
func TestUpdateUserWithoutAuth(t *testing.T) {
	baseURL := getTestServer(t)

	client := &http.Client{Timeout: 5 * time.Second}

	updateBody := map[string]any{
		"email": "test@example.com",
		"phone": "13800138000",
	}
	updateResp := postJSON[any](t, client, baseURL+"/auth/update", updateBody, nil)

	if updateResp.Code == uint64(perrors.Success.Code) {
		t.Fatalf("expected update without auth to fail, but succeeded")
	}
	if updateResp.Code != uint64(perrors.ErrAuthFailed.Code) {
		t.Logf("expected error code %d, got %d", perrors.ErrAuthFailed.Code, updateResp.Code)
	}
}

// ========== Update Password Edge Cases ==========

// TestUpdatePasswordWrongOldPassword tests password update with wrong old password
func TestUpdatePasswordWrongOldPassword(t *testing.T) {
	baseURL := getTestServer(t)

	client := &http.Client{Timeout: 5 * time.Second}

	suffix := time.Now().UnixNano()
	_, password, token := createAndLoginTestUser(t, client, baseURL, suffix)
	authHeader := map[string]string{"Authorization": fmt.Sprintf("Bearer %s", token)}

	// Try to update password with wrong old password
	updatePwdBody := map[string]any{
		"old_password": "WrongOldPassword!",
		"new_password": "NewPassw0rd!",
	}
	updateResp := postJSON[any](t, client, baseURL+"/auth/password", updatePwdBody, authHeader)

	if updateResp.Code == uint64(perrors.Success.Code) {
		t.Fatalf("expected password update with wrong old password to fail, but succeeded")
	}

	// Verify original password still works
	loginResp := postJSON[loginData](t, client, baseURL+"/login", map[string]any{
		"username": fmt.Sprintf("tester_%d", suffix%1000000),
		"password": password,
	}, nil)
	if loginResp.Code != uint64(perrors.Success.Code) {
		t.Fatalf("original password should still work after failed update")
	}
}

// TestUpdatePasswordEmpty tests password update with empty fields
func TestUpdatePasswordEmpty(t *testing.T) {
	baseURL := getTestServer(t)

	client := &http.Client{Timeout: 5 * time.Second}

	suffix := time.Now().UnixNano()
	_, _, token := createAndLoginTestUser(t, client, baseURL, suffix)
	authHeader := map[string]string{"Authorization": fmt.Sprintf("Bearer %s", token)}

	// Test empty old password
	updatePwdBody1 := map[string]any{
		"old_password": "",
		"new_password": "NewPassw0rd!",
	}
	updateResp1 := postJSON[any](t, client, baseURL+"/auth/password", updatePwdBody1, authHeader)
	if updateResp1.Code == uint64(perrors.Success.Code) {
		t.Fatalf("expected password update with empty old password to fail, but succeeded")
	}

	// Test empty new password
	updatePwdBody2 := map[string]any{
		"old_password": "Passw0rd!",
		"new_password": "",
	}
	updateResp2 := postJSON[any](t, client, baseURL+"/auth/password", updatePwdBody2, authHeader)
	if updateResp2.Code == uint64(perrors.Success.Code) {
		t.Fatalf("expected password update with empty new password to fail, but succeeded")
	}
}

// TestUpdatePasswordWithoutAuth tests password update without authentication
func TestUpdatePasswordWithoutAuth(t *testing.T) {
	baseURL := getTestServer(t)

	client := &http.Client{Timeout: 5 * time.Second}

	updatePwdBody := map[string]any{
		"old_password": "OldPassw0rd!",
		"new_password": "NewPassw0rd!",
	}
	updateResp := postJSON[any](t, client, baseURL+"/auth/password", updatePwdBody, nil)

	if updateResp.Code == uint64(perrors.Success.Code) {
		t.Fatalf("expected password update without auth to fail, but succeeded")
	}
	if updateResp.Code != uint64(perrors.ErrAuthFailed.Code) {
		t.Logf("expected error code %d, got %d", perrors.ErrAuthFailed.Code, updateResp.Code)
	}
}

// TestUpdatePasswordTokenInvalidation tests that old token becomes invalid after password change
func TestUpdatePasswordTokenInvalidation(t *testing.T) {
	baseURL := getTestServer(t)

	client := &http.Client{Timeout: 5 * time.Second}

	suffix := time.Now().UnixNano()
	username := fmt.Sprintf("pwdtok_%d", suffix%1000000) // Use unique prefix and larger range
	oldPassword := "OldPassw0rd!"
	newPassword := "NewPassw0rd!"

	// 1) Register user
	regBody := map[string]any{
		"username":         username,
		"password":         oldPassword,
		"password_confirm": oldPassword,
		"email":            fmt.Sprintf("%s@example.com", username),
	}
	regResp := postJSON[map[string]any](t, client, baseURL+"/register", regBody, nil)
	if regResp.Code != uint64(perrors.Success.Code) {
		t.Fatalf("registration failed: code=%d msg=%s", regResp.Code, regResp.Message)
	}

	// 2) Login and get token
	loginBody := map[string]any{
		"username": username,
		"password": oldPassword,
	}
	loginResp := postJSON[loginData](t, client, baseURL+"/login", loginBody, nil)
	if loginResp.Code != uint64(perrors.Success.Code) {
		t.Fatalf("login failed: code=%d msg=%s", loginResp.Code, loginResp.Message)
	}
	oldToken := loginResp.Data.Token
	oldAuthHeader := map[string]string{"Authorization": fmt.Sprintf("Bearer %s", oldToken)}

	// 3) Verify old token works before password change
	infoResp1 := getJSON[map[string]any](t, client, baseURL+"/auth/info", oldAuthHeader)
	if infoResp1.Code != uint64(perrors.Success.Code) {
		t.Fatalf("old token should work before password change: code=%d msg=%s", infoResp1.Code, infoResp1.Message)
	}

	// 4) Update password using old token
	updatePwdBody := map[string]any{
		"old_password": oldPassword,
		"new_password": newPassword,
	}
	updateResp := postJSON[map[string]any](t, client, baseURL+"/auth/password", updatePwdBody, oldAuthHeader)
	if updateResp.Code != uint64(perrors.Success.Code) {
		t.Fatalf("update password failed: code=%d msg=%s", updateResp.Code, updateResp.Message)
	}

	// 5) Try to access user info with old token after password change (should fail)
	infoResp2 := getJSON[any](t, client, baseURL+"/auth/info", oldAuthHeader)
	if infoResp2.Code == uint64(perrors.Success.Code) {
		t.Fatalf("old token should be invalid after password change, but still works")
	}
	if infoResp2.Code != uint64(perrors.ErrAuthFailed.Code) {
		t.Logf("expected error code %d, got %d, msg: %s", perrors.ErrAuthFailed.Code, infoResp2.Code, infoResp2.Message)
	}

	// Sleep to ensure new token has different orig_iat timestamp (generated in different second)
	time.Sleep(time.Second)

	// 6) Login with new password to get new token
	loginNewBody := map[string]any{
		"username": username,
		"password": newPassword,
	}
	loginNewResp := postJSON[loginData](t, client, baseURL+"/login", loginNewBody, nil)
	if loginNewResp.Code != uint64(perrors.Success.Code) {
		t.Fatalf("login with new password failed: code=%d msg=%s", loginNewResp.Code, loginNewResp.Message)
	}
	newToken := loginNewResp.Data.Token
	newAuthHeader := map[string]string{"Authorization": fmt.Sprintf("Bearer %s", newToken)}

	// 7) Verify new token works
	infoResp3 := getJSON[map[string]any](t, client, baseURL+"/auth/info", newAuthHeader)
	if infoResp3.Code != uint64(perrors.Success.Code) {
		t.Fatalf("new token should work after password change: code=%d msg=%s", infoResp3.Code, infoResp3.Message)
	}

	// 8) Verify we got the correct user info
	if infoResp3.Data["user"] == nil {
		t.Fatalf("user data is nil")
	}
	user := infoResp3.Data["user"].(map[string]any)
	if user["username"] != username {
		t.Errorf("username mismatch: expected %s, got %v", username, user["username"])
	}
}

// ========== Token and Session Tests ==========

// TestMultipleLoginSessions tests multiple login sessions with same user
func TestMultipleLoginSessions(t *testing.T) {
	baseURL := getTestServer(t)

	client := &http.Client{Timeout: 5 * time.Second}

	suffix := time.Now().UnixNano()
	username, password := createTestUser(t, client, baseURL, suffix)

	// Login first session
	token1 := loginTestUser(t, client, baseURL, username, password)
	authHeader1 := map[string]string{"Authorization": fmt.Sprintf("Bearer %s", token1)}

	time.Sleep(time.Second)
	// Login second session
	token2 := loginTestUser(t, client, baseURL, username, password)
	authHeader2 := map[string]string{"Authorization": fmt.Sprintf("Bearer %s", token2)}

	// Both tokens should work
	infoResp1 := getJSON[map[string]any](t, client, baseURL+"/auth/info", authHeader1)
	if infoResp1.Code != uint64(perrors.Success.Code) {
		t.Fatalf("first token should still work: code=%d msg=%s", infoResp1.Code, infoResp1.Message)
	}

	infoResp2 := getJSON[map[string]any](t, client, baseURL+"/auth/info", authHeader2)
	if infoResp2.Code != uint64(perrors.Success.Code) {
		t.Fatalf("second token should work: code=%d msg=%s", infoResp2.Code, infoResp2.Message)
	}

	// Logout from first session
	logoutResp := postJSON[map[string]any](t, client, baseURL+"/logout", map[string]any{}, authHeader1)
	if logoutResp.Code != uint64(perrors.Success.Code) {
		t.Fatalf("logout failed: code=%d msg=%s", logoutResp.Code, logoutResp.Message)
	}

	// First token should no longer work
	infoResp3 := getJSON[any](t, client, baseURL+"/auth/info", authHeader1)
	if infoResp3.Code == uint64(perrors.Success.Code) {
		t.Fatalf("first token should be invalid after logout")
	}

	// Second token should still work
	infoResp4 := getJSON[map[string]any](t, client, baseURL+"/auth/info", authHeader2)
	if infoResp4.Code != uint64(perrors.Success.Code) {
		t.Fatalf("second token should still work: code=%d msg=%s", infoResp4.Code, infoResp4.Message)
	}
}

// TestExpiredTokenFormat tests token with invalid format
func TestMalformedToken(t *testing.T) {
	baseURL := getTestServer(t)

	client := &http.Client{Timeout: 5 * time.Second}

	testCases := []struct {
		name  string
		token string
	}{
		{"empty token", "Bearer "},
		{"no bearer prefix", "some-random-token"},
		{"malformed jwt", "Bearer abc.def.ghi"},
		{"missing parts", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			authHeader := map[string]string{"Authorization": tc.token}
			infoResp := getJSON[any](t, client, baseURL+"/auth/info", authHeader)

			if infoResp.Code == uint64(perrors.Success.Code) {
				t.Fatalf("expected malformed token to fail, but succeeded")
			}
		})
	}
}

// ========== Data Validation Tests ==========

// TestRegisterSpecialCharacters tests registration with special characters in username
func TestRegisterSpecialCharacters(t *testing.T) {
	baseURL := getTestServer(t)

	client := &http.Client{Timeout: 5 * time.Second}

	suffix := time.Now().UnixNano()
	testCases := []struct {
		username string
		shouldOK bool
	}{
		{"user@name", false},     // @ symbol
		{"user name", false},     // space
		{"user<script>", false},  // HTML tags
		{"user'OR'1'='1", false}, // SQL injection attempt
		{"user\nname", false},    // newline
		{"normaluser123", true},  // valid username
	}

	for i, tc := range testCases {
		t.Run(tc.username, func(t *testing.T) {
			regBody := map[string]any{
				"username":         tc.username + fmt.Sprintf("_%d", suffix%1000), // ensure uniqueness
				"password":         "Passw0rd!",
				"password_confirm": "Passw0rd!",
				"email":            fmt.Sprintf("test_%d@example.com", suffix+int64(i)),
			}
			regResp := postJSON[any](t, client, baseURL+"/register", regBody, nil)

			if tc.shouldOK && regResp.Code != uint64(perrors.Success.Code) {
				t.Errorf("expected registration to succeed for '%s', but failed: code=%d msg=%s",
					tc.username, regResp.Code, regResp.Message)
			} else if !tc.shouldOK && regResp.Code == uint64(perrors.Success.Code) {
				t.Errorf("expected registration to fail for '%s', but succeeded", tc.username)
			}
			if !tc.shouldOK && regResp.Code != uint64(perrors.ErrParam.Code) {
				t.Logf("expected error code %d, got %d for username '%s'",
					perrors.ErrParam.Code, regResp.Code, tc.username)
			}
		})
	}
}

// TestUpdateUserPartialFields tests updating only some user fields
func TestUpdateUserPartialFields(t *testing.T) {
	baseURL := getTestServer(t)

	client := &http.Client{Timeout: 5 * time.Second}

	suffix := time.Now().UnixNano()
	username, _, token := createAndLoginTestUser(t, client, baseURL, suffix)
	authHeader := map[string]string{"Authorization": fmt.Sprintf("Bearer %s", token)}

	// Update only email
	updateBody1 := map[string]any{
		"email": fmt.Sprintf("newemail_%d@example.com", suffix),
	}
	updateResp1 := postJSON[map[string]any](t, client, baseURL+"/auth/update", updateBody1, authHeader)
	if updateResp1.Code != uint64(perrors.Success.Code) {
		t.Fatalf("update email only failed: code=%d msg=%s", updateResp1.Code, updateResp1.Message)
	}

	// Update only phone
	updateBody2 := map[string]any{
		"phone": "13900139000",
	}
	updateResp2 := postJSON[map[string]any](t, client, baseURL+"/auth/update", updateBody2, authHeader)
	if updateResp2.Code != uint64(perrors.Success.Code) {
		t.Fatalf("update phone only failed: code=%d msg=%s", updateResp2.Code, updateResp2.Message)
	}

	// Verify both updates persisted
	infoResp := getJSON[map[string]any](t, client, baseURL+"/auth/info", authHeader)
	if infoResp.Code != uint64(perrors.Success.Code) {
		t.Fatalf("get user info failed: code=%d msg=%s", infoResp.Code, infoResp.Message)
	}

	t.Logf("Updated user info for %s: %+v", username, infoResp.Data)
}

// TestGetUserInfoImmediatelyAfterRegistration tests getting user info right after registration
func TestGetUserInfoImmediatelyAfterRegistration(t *testing.T) {
	baseURL := getTestServer(t)

	client := &http.Client{Timeout: 5 * time.Second}

	suffix := time.Now().UnixNano()
	username, password, token := createAndLoginTestUser(t, client, baseURL, suffix)
	authHeader := map[string]string{"Authorization": fmt.Sprintf("Bearer %s", token)}

	// Get user info immediately
	infoResp := getJSON[map[string]any](t, client, baseURL+"/auth/info", authHeader)
	if infoResp.Code != uint64(perrors.Success.Code) {
		t.Fatalf("get user info failed: code=%d msg=%s", infoResp.Code, infoResp.Message)
	}

	// Verify data integrity
	if infoResp.Data["user"] == nil {
		t.Fatalf("user data is nil")
	}

	user := infoResp.Data["user"].(map[string]any)
	if user["username"] != username {
		t.Errorf("username mismatch: expected %s, got %v", username, user["username"])
	}

	t.Logf("User %s with password %s registered and retrieved successfully", username, password)
}
