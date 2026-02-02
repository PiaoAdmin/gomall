package service

import (
	"context"

	"github.com/PiaoAdmin/pmall/app/product/biz/dal/mysql"
	"github.com/PiaoAdmin/pmall/app/product/biz/dal/redis"
	"github.com/PiaoAdmin/pmall/app/product/biz/model"
	"github.com/PiaoAdmin/pmall/app/product/biz/utils"
	"github.com/PiaoAdmin/pmall/common/errs"
	product "github.com/PiaoAdmin/pmall/rpc_gen/product"
	"github.com/cloudwego/kitex/pkg/klog"
)

type UpdateProductStatusService struct {
	ctx context.Context
}

func NewUpdateProductStatusService(ctx context.Context) *UpdateProductStatusService {
	return &UpdateProductStatusService{ctx: ctx}
}

func (s *UpdateProductStatusService) Run(req *product.UpdateProductStatusRequest) (*product.UpdateProductStatusResponse, error) {
	if len(req.Ids) == 0 {
		return nil, errs.New(errs.ErrParam.Code, "ids is empty")
	}

	updates := make(map[string]interface{})

	if req.PublishStatus > 0 {
		if !utils.IsValidPublishStatus(int8(req.PublishStatus)) {
			return nil, errs.New(errs.ErrParam.Code, "invalid publish_status")
		}
		updates["publish_status"] = req.PublishStatus
	}

	if req.VerifyStatus > 0 {
		if !utils.IsValidVerifyStatus(int8(req.VerifyStatus)) {
			return nil, errs.New(errs.ErrParam.Code, "invalid verify_status")
		}
		updates["verify_status"] = req.VerifyStatus
	}

	if len(updates) == 0 {
		return &product.UpdateProductStatusResponse{
			Success: true,
		}, nil
	}

	err := model.BatchUpdateStatus(s.ctx, mysql.DB, req.Ids, updates)
	if err != nil {
		return nil, errs.New(errs.ErrInternal.Code, "update product status failed: "+err.Error())
	}

	// 异步清除相关商品缓存
	go func() {
		if err := redis.BatchDeleteProductDetailCache(s.ctx, req.Ids); err != nil {
			klog.Warnf("Failed to batch delete product detail cache: %v", err)
		}
		// 清除商品列表缓存
		if err := redis.InvalidateProductListCache(s.ctx); err != nil {
			klog.Warnf("Failed to invalidate product list cache: %v", err)
		}
		// 清除热门商品缓存（状态变更可能影响热门商品列表）
		if err := redis.DeleteHotProductsCache(s.ctx); err != nil {
			klog.Warnf("Failed to delete hot products cache: %v", err)
		}
	}()

	return &product.UpdateProductStatusResponse{
		Success: true,
	}, nil
}
