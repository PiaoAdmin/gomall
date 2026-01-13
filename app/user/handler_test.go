package main

import (
	"context"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/PiaoAdmin/pmall/app/user/biz/dal/mysql"
	"github.com/PiaoAdmin/pmall/app/user/biz/model"
	"github.com/PiaoAdmin/pmall/app/user/conf"
	user "github.com/PiaoAdmin/pmall/rpc_gen/user"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	// 设置环境为 test，确保加载 test/conf.yaml
	os.Setenv("GO_ENV", "test")
	// 初始化配置
	_ = conf.GetConf()
	// 初始化数据库
	mysql.Init()
	// 运行测试
	code := m.Run()
	os.Exit(code)
}

func getHandler() *UserServiceImpl {
	return &UserServiceImpl{}
}

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func TestUserFlow(t *testing.T) {
	handler := getHandler()
	ctx := context.Background()

	// 生成随机用户名和邮箱，避免冲突
	randStr := generateRandomString(8)
	username := "test_user_" + randStr
	email := "test_" + randStr + "@example.com"
	password := "password123"

	var userId uint64

	t.Run("Register", func(t *testing.T) {
		req := &user.RegisterRequest{
			Username:        username,
			Password:        password,
			PasswordConfirm: password,
			Email:           email,
		}
		resp, err := handler.Register(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotZero(t, resp.User.Id)
		assert.Equal(t, username, resp.User.Username)
		assert.Equal(t, email, resp.User.Email)
		userId = resp.User.Id
	})

	t.Run("RegisterDuplicate", func(t *testing.T) {
		req := &user.RegisterRequest{
			Username:        username, // 重复用户名
			Password:        password,
			PasswordConfirm: password,
			Email:           "new_" + email,
		}
		resp, err := handler.Register(ctx, req)
		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("Login", func(t *testing.T) {
		req := &user.LoginRequest{
			Username: username,
			Password: password,
		}
		resp, err := handler.Login(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, userId, resp.User.Id)
	})

	t.Run("LoginWrongPassword", func(t *testing.T) {
		req := &user.LoginRequest{
			Username: username,
			Password: "wrong_password",
		}
		resp, err := handler.Login(ctx, req)
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "password incorrect")
	})

	t.Run("GetUserInfo", func(t *testing.T) {
		req := &user.GetUserInfoRequest{
			UserId: userId,
		}
		resp, err := handler.GetUserInfo(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, userId, resp.User.Id)
		assert.Equal(t, username, resp.User.Username)
	})

	t.Run("UpdateUser", func(t *testing.T) {
		newEmail := "updated_" + randStr + "@example.com"
		req := &user.UpdateUserRequest{
			UserId: userId,
			Email:  newEmail,
			Phone:  "1234567890",
		}
		resp, err := handler.UpdateUser(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.True(t, resp.Success)
		assert.Equal(t, newEmail, resp.User.Email)
		assert.Equal(t, "1234567890", resp.User.Phone)
	})

	t.Run("UpdatePassword", func(t *testing.T) {
		newPassword := "new_password_123"
		req := &user.UpdatePasswordRequest{
			UserId:      userId,
			OldPassword: password,
			NewPassword: newPassword,
		}
		resp, err := handler.UpdatePassword(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.True(t, resp.Success)

		// 验证旧密码失效
		loginReq := &user.LoginRequest{
			Username: username,
			Password: password,
		}
		_, err = handler.Login(ctx, loginReq)
		assert.Error(t, err)

		// 验证新密码生效
		loginReq.Password = newPassword
		loginResp, err := handler.Login(ctx, loginReq)
		assert.NoError(t, err)
		assert.NotNil(t, loginResp)
	})

	// 清理数据
	t.Cleanup(func() {
		if userId != 0 {
			// 硬删除，彻底清理
			mysql.DB.Unscoped().Delete(&model.User{}, userId)
		}
	})
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
