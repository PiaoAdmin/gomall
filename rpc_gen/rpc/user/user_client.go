package user

import (
	"context"
	user "github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/user"

	"github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/user/userservice"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/callopt"
)

type RPCClient interface {
	KitexClient() userservice.Client
	Service() string
	CreateUser(ctx context.Context, Req *user.CreateUserRequest, callOptions ...callopt.Option) (r *user.CreateUserResponse, err error)
	DeleteUser(ctx context.Context, Req *user.DeleteUserRequest, callOptions ...callopt.Option) (r *user.DeleteUserResponse, err error)
	UpdateBaseUser(ctx context.Context, Req *user.UpdateBaseUserRequest, callOptions ...callopt.Option) (r *user.UpdateBaseUserResponse, err error)
	UpdateUserPassword(ctx context.Context, Req *user.UpdateUserPasswordRequest, callOptions ...callopt.Option) (r *user.UpdateUserPasswordResponse, err error)
	UpdateUserBalance(ctx context.Context, Req *user.UpdateUserBalanceRequest, callOptions ...callopt.Option) (r *user.UpdateUserBalanceResponse, err error)
	UpdateUserStatus(ctx context.Context, Req *user.UpdateUserStatusRequest, callOptions ...callopt.Option) (r *user.UpdateUserStatusResponse, err error)
	GetUserInfo(ctx context.Context, Req *user.GetUserInfoRequest, callOptions ...callopt.Option) (r *user.GetUserInfoResponse, err error)
}

func NewRPCClient(dstService string, opts ...client.Option) (RPCClient, error) {
	kitexClient, err := userservice.NewClient(dstService, opts...)
	if err != nil {
		return nil, err
	}
	cli := &clientImpl{
		service:     dstService,
		kitexClient: kitexClient,
	}

	return cli, nil
}

type clientImpl struct {
	service     string
	kitexClient userservice.Client
}

func (c *clientImpl) Service() string {
	return c.service
}

func (c *clientImpl) KitexClient() userservice.Client {
	return c.kitexClient
}

func (c *clientImpl) CreateUser(ctx context.Context, Req *user.CreateUserRequest, callOptions ...callopt.Option) (r *user.CreateUserResponse, err error) {
	return c.kitexClient.CreateUser(ctx, Req, callOptions...)
}

func (c *clientImpl) DeleteUser(ctx context.Context, Req *user.DeleteUserRequest, callOptions ...callopt.Option) (r *user.DeleteUserResponse, err error) {
	return c.kitexClient.DeleteUser(ctx, Req, callOptions...)
}

func (c *clientImpl) UpdateBaseUser(ctx context.Context, Req *user.UpdateBaseUserRequest, callOptions ...callopt.Option) (r *user.UpdateBaseUserResponse, err error) {
	return c.kitexClient.UpdateBaseUser(ctx, Req, callOptions...)
}

func (c *clientImpl) UpdateUserPassword(ctx context.Context, Req *user.UpdateUserPasswordRequest, callOptions ...callopt.Option) (r *user.UpdateUserPasswordResponse, err error) {
	return c.kitexClient.UpdateUserPassword(ctx, Req, callOptions...)
}

func (c *clientImpl) UpdateUserBalance(ctx context.Context, Req *user.UpdateUserBalanceRequest, callOptions ...callopt.Option) (r *user.UpdateUserBalanceResponse, err error) {
	return c.kitexClient.UpdateUserBalance(ctx, Req, callOptions...)
}

func (c *clientImpl) UpdateUserStatus(ctx context.Context, Req *user.UpdateUserStatusRequest, callOptions ...callopt.Option) (r *user.UpdateUserStatusResponse, err error) {
	return c.kitexClient.UpdateUserStatus(ctx, Req, callOptions...)
}

func (c *clientImpl) GetUserInfo(ctx context.Context, Req *user.GetUserInfoRequest, callOptions ...callopt.Option) (r *user.GetUserInfoResponse, err error) {
	return c.kitexClient.GetUserInfo(ctx, Req, callOptions...)
}
