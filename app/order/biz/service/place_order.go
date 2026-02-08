package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/PiaoAdmin/pmall/app/order/biz/dal/rabbitmq"
	"github.com/PiaoAdmin/pmall/app/order/biz/rpc"
	"github.com/PiaoAdmin/pmall/common/errs"
	"github.com/PiaoAdmin/pmall/common/uniqueid"
	order "github.com/PiaoAdmin/pmall/rpc_gen/order"
	"github.com/PiaoAdmin/pmall/rpc_gen/product"
	"github.com/cloudwego/kitex/pkg/klog"
)

type PlaceOrderService struct {
	ctx context.Context
}

func NewPlaceOrderService(ctx context.Context) *PlaceOrderService {
	return &PlaceOrderService{ctx: ctx}
}

func (s *PlaceOrderService) Run(req *order.PlaceOrderReq) (*order.PlaceOrderResp, error) {
	if req == nil || req.UserId == 0 {
		return nil, errs.New(errs.ErrParam.Code, "invalid request")
	}
	if len(req.Items) == 0 {
		return nil, errs.New(errs.ErrParam.Code, "items empty")
	}

	newOrderId := fmt.Sprintf("%d", uniqueid.GenId())

	deductItems := make([]*product.SkuDeductItem, 0, len(req.Items))
	for _, it := range req.Items {
		if it == nil || it.GetSkuId() == 0 || it.GetQuantity() <= 0 {
			return nil, errs.New(errs.ErrParam.Code, "invalid sku_id or quantity")
		}
		deductItems = append(deductItems, &product.SkuDeductItem{
			SkuId: it.GetSkuId(),
			Count: it.GetQuantity(),
		})
	}

	// 1. 先扣减库存（同步操作，保证库存准确性）
	if _, err := rpc.ProductClient.DeductStock(s.ctx, &product.DeductStockRequest{
		OrderSn: newOrderId,
		Items:   deductItems,
	}); err != nil {
		return nil, errs.New(errs.ErrInternal.Code, "deduct stock failed: "+err.Error())
	}

	// 2. 构建订单消息
	orderMsg := &rabbitmq.OrderMessage{
		OrderID:   newOrderId,
		UserID:    req.UserId,
		Email:     req.Email,
		CreatedAt: time.Now().Unix(),
		Retry:     0,
	}

	if req.ShippingAddress != nil {
		orderMsg.Address = rabbitmq.OrderAddress{
			Name:          req.ShippingAddress.GetName(),
			StreetAddress: req.ShippingAddress.GetStreetAddress(),
			City:          req.ShippingAddress.GetCity(),
			ZipCode:       req.ShippingAddress.GetZipCode(),
		}
	}

	// 3. 构建订单项
	orderMsg.Items = make([]rabbitmq.OrderMessageItem, 0, len(req.Items))
	for _, it := range req.Items {
		price := 0.0
		if it != nil && it.GetPrice() != "" {
			if p, perr := strconv.ParseFloat(it.GetPrice(), 64); perr == nil {
				price = p
			}
		}
		orderMsg.Items = append(orderMsg.Items, rabbitmq.OrderMessageItem{
			SkuID:    it.GetSkuId(),
			SkuName:  it.GetSkuName(),
			Price:    price,
			Quantity: it.GetQuantity(),
		})
	}

	// 4. 异步发送消息到 RabbitMQ（订单数据库写入将由消费者完成）
	if err := rabbitmq.PublishOrderMessage(s.ctx, orderMsg); err != nil {
		klog.CtxErrorf(s.ctx, "Failed to publish order message, falling back to sync: %v", err)
		// 发送失败时回退库存
		if _, relErr := rpc.ProductClient.ReleaseStock(s.ctx, &product.ReleaseStockRequest{
			OrderSn: newOrderId,
			Items:   deductItems,
		}); relErr != nil {
			klog.CtxErrorf(s.ctx, "Release stock failed: %v", relErr)
		}
		return nil, errs.New(errs.ErrInternal.Code, "place order failed: "+err.Error())
	}

	// 5. 发送延迟取消消息（30分钟后未支付则自动取消）
	if err := rabbitmq.PublishOrderCancelDelay(s.ctx, newOrderId, req.UserId); err != nil {
		// 延迟消息发送失败不影响下单，只记录日志
		klog.CtxWarnf(s.ctx, "Failed to publish cancel delay message: %v", err)
	}

	klog.CtxInfof(s.ctx, "Order placed successfully (async): order_id=%s", newOrderId)

	return &order.PlaceOrderResp{
		Order: &order.OrderResult{OrderId: newOrderId},
	}, nil
}
