package auth

import (
	"context"
	auth "github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/auth"

	"github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/auth/authservice"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/callopt"
)

type RPCClient interface {
	KitexClient() authservice.Client
	Service() string
	IssueToken(ctx context.Context, Req *auth.IssueTokenRequest, callOptions ...callopt.Option) (r *auth.IssueTokenResponse, err error)
	VerifyToken(ctx context.Context, Req *auth.VerifyTokenRequest, callOptions ...callopt.Option) (r *auth.VerifyTokenResponse, err error)
	RefreshToken(ctx context.Context, Req *auth.RefreshTokenRequest, callOptions ...callopt.Option) (r *auth.RefreshTokenResponse, err error)
	Logout(ctx context.Context, Req *auth.LogoutRequest, callOptions ...callopt.Option) (r *auth.LogoutResponse, err error)
	CheckPermission(ctx context.Context, Req *auth.PermissionCheckRequest, callOptions ...callopt.Option) (r *auth.PermissionCheckResponse, err error)
	AddUserRole(ctx context.Context, Req *auth.AddUserRoleRequest, callOptions ...callopt.Option) (r *auth.AddUserRoleResponse, err error)
	RemoveUserRole(ctx context.Context, Req *auth.RemoveUserRoleRequest, callOptions ...callopt.Option) (r *auth.RemoveUserRoleResponse, err error)
	GetUserRole(ctx context.Context, Req *auth.GetUserRoleRequest, callOptions ...callopt.Option) (r *auth.GetUserRoleResponse, err error)
}

func NewRPCClient(dstService string, opts ...client.Option) (RPCClient, error) {
	kitexClient, err := authservice.NewClient(dstService, opts...)
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
	kitexClient authservice.Client
}

func (c *clientImpl) Service() string {
	return c.service
}

func (c *clientImpl) KitexClient() authservice.Client {
	return c.kitexClient
}

func (c *clientImpl) IssueToken(ctx context.Context, Req *auth.IssueTokenRequest, callOptions ...callopt.Option) (r *auth.IssueTokenResponse, err error) {
	return c.kitexClient.IssueToken(ctx, Req, callOptions...)
}

func (c *clientImpl) VerifyToken(ctx context.Context, Req *auth.VerifyTokenRequest, callOptions ...callopt.Option) (r *auth.VerifyTokenResponse, err error) {
	return c.kitexClient.VerifyToken(ctx, Req, callOptions...)
}

func (c *clientImpl) RefreshToken(ctx context.Context, Req *auth.RefreshTokenRequest, callOptions ...callopt.Option) (r *auth.RefreshTokenResponse, err error) {
	return c.kitexClient.RefreshToken(ctx, Req, callOptions...)
}

func (c *clientImpl) Logout(ctx context.Context, Req *auth.LogoutRequest, callOptions ...callopt.Option) (r *auth.LogoutResponse, err error) {
	return c.kitexClient.Logout(ctx, Req, callOptions...)
}

func (c *clientImpl) CheckPermission(ctx context.Context, Req *auth.PermissionCheckRequest, callOptions ...callopt.Option) (r *auth.PermissionCheckResponse, err error) {
	return c.kitexClient.CheckPermission(ctx, Req, callOptions...)
}

func (c *clientImpl) AddUserRole(ctx context.Context, Req *auth.AddUserRoleRequest, callOptions ...callopt.Option) (r *auth.AddUserRoleResponse, err error) {
	return c.kitexClient.AddUserRole(ctx, Req, callOptions...)
}

func (c *clientImpl) RemoveUserRole(ctx context.Context, Req *auth.RemoveUserRoleRequest, callOptions ...callopt.Option) (r *auth.RemoveUserRoleResponse, err error) {
	return c.kitexClient.RemoveUserRole(ctx, Req, callOptions...)
}

func (c *clientImpl) GetUserRole(ctx context.Context, Req *auth.GetUserRoleRequest, callOptions ...callopt.Option) (r *auth.GetUserRoleResponse, err error) {
	return c.kitexClient.GetUserRole(ctx, Req, callOptions...)
}
