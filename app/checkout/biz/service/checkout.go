package service

import (
	"context"
	"strconv"

	"github.com/PiaoAdmin/pmall/app/checkout/biz/rpc"
	"github.com/PiaoAdmin/pmall/common/errs"
	checkout "github.com/PiaoAdmin/pmall/rpc_gen/checkout"
	"github.com/PiaoAdmin/pmall/rpc_gen/order"
	"github.com/PiaoAdmin/pmall/rpc_gen/payment"
	"github.com/PiaoAdmin/pmall/rpc_gen/product"
	"github.com/PiaoAdmin/pmall/rpc_gen/user"
	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/cloudwego/kitex/pkg/klog"
)

type CheckoutService struct {
	ctx context.Context
}

func NewCheckoutService(ctx context.Context) *CheckoutService {
	return &CheckoutService{ctx: ctx}
}

func (s *CheckoutService) Run(req *checkout.CheckoutRequest) (*checkout.CheckoutResponse, error) {
	if req == nil || req.UserId == 0 {
		return nil, errs.New(errs.ErrParam.Code, "invalid request")
	}
	if len(req.Items) == 0 {
		return nil, errs.New(errs.ErrParam.Code, "items empty")
	}
	if req.ShippingAddress == nil {
		return nil, errs.New(errs.ErrParam.Code, "shipping_address empty")
	}

	skuOrder, qtyMap, err := s.normalizeItems(req.Items)
	if err != nil {
		return nil, err
	}

	userResp, err := rpc.UserClient.GetUserInfo(s.ctx, &user.GetUserInfoRequest{UserId: req.UserId})
	if err != nil {
		return nil, wrapRPC(err, "get user info failed")
	}
	if userResp == nil || userResp.User == nil {
		return nil, errs.New(errs.ErrRecordNotFound.Code, "user not found")
	}

	skuResp, err := rpc.ProductClient.GetSkusByIds(s.ctx, &product.GetSkusByIdsRequest{SkuIds: skuOrder})
	if err != nil {
		return nil, wrapRPC(err, "get skus failed")
	}
	if skuResp == nil || len(skuResp.Skus) == 0 {
		return nil, errs.New(errs.ErrRecordNotFound.Code, "sku not found")
	}

	orderItems := make([]*order.CartItem, 0, len(skuOrder))
	resultItems := make([]*checkout.CheckoutItemResult, 0, len(skuOrder))
	var totalAmount float64

	for _, skuID := range skuOrder {
		qty := qtyMap[skuID]
		sku, ok := skuResp.Skus[skuID]
		if !ok || sku == nil {
			return nil, errs.New(errs.ErrRecordNotFound.Code, "sku not found")
		}
		if sku.Stock < qty {
			return nil, errs.New(errs.ErrParam.Code, "stock not enough")
		}

		orderItems = append(orderItems, &order.CartItem{
			SkuId:    skuID,
			Quantity: qty,
			SkuName:  sku.Name,
			Price:    sku.Price,
		})
		resultItems = append(resultItems, &checkout.CheckoutItemResult{
			SkuId:       skuID,
			Quantity:    qty,
			SkuName:     sku.Name,
			SkuImage:    sku.MainImage,
			Price:       sku.Price,
			MarketPrice: sku.MarketPrice,
			SpuId:       sku.SpuId,
			SkuSpecData: sku.SkuSpecData,
		})
		if price, perr := strconv.ParseFloat(sku.Price, 64); perr == nil {
			totalAmount += price * float64(qty)
		}
	}

	placeResp, err := rpc.OrderClient.PlaceOrder(s.ctx, &order.PlaceOrderReq{
		UserId: req.UserId,
		Email:  userResp.User.Email,
		Items:  orderItems,
		ShippingAddress: &order.Address{
			StreetAddress: req.ShippingAddress.GetStreetAddress(),
			City:          req.ShippingAddress.GetCity(),
			Name:          req.ShippingAddress.GetName(),
			ZipCode:       req.ShippingAddress.GetZipCode(),
		},
	})
	if err != nil {
		return nil, wrapRPC(err, "place order failed")
	}
	if placeResp == nil || placeResp.Order == nil || placeResp.Order.OrderId == "" {
		return nil, errs.New(errs.ErrInternal.Code, "place order failed: empty order_id")
	}

	payResp, err := rpc.PaymentClient.Pay(s.ctx, &payment.PayRequest{
		OrderId:    placeResp.Order.OrderId,
		UserId:     req.UserId,
		Amount:     strconv.FormatFloat(totalAmount, 'f', 2, 64),
		CreditCard: req.CreditCard,
	})
	if err != nil || payResp == nil || !payResp.Success {
		if err != nil {
			klog.CtxErrorf(s.ctx, "Pay failed: %v", err)
		} else {
			klog.CtxErrorf(s.ctx, "Pay failed: empty response")
		}
		if _, cancelErr := rpc.OrderClient.CancelOrder(s.ctx, &order.CancelOrderReq{OrderId: placeResp.Order.OrderId}); cancelErr != nil {
			klog.CtxWarnf(s.ctx, "CancelOrder failed: %v", cancelErr)
		}
		if err != nil {
			return nil, wrapRPC(err, "payment failed")
		}
		return nil, errs.New(errs.ErrInternal.Code, "payment failed")
	}

	_, err = rpc.OrderClient.MarkOrderPaid(s.ctx, &order.MarkOrderPaidReq{
		UserId:  req.UserId,
		OrderId: placeResp.Order.OrderId,
	})
	if err != nil {
		return nil, wrapRPC(err, "mark order paid failed")
	}

	return &checkout.CheckoutResponse{
		OrderId:     placeResp.Order.OrderId,
		TotalAmount: strconv.FormatFloat(totalAmount, 'f', 2, 64),
		Items:       resultItems,
	}, nil
}

func (s *CheckoutService) normalizeItems(items []*checkout.CheckoutItem) ([]uint64, map[uint64]int32, error) {
	qtyMap := make(map[uint64]int32)
	order := make([]uint64, 0, len(items))
	seen := make(map[uint64]struct{})

	for _, item := range items {
		if item == nil || item.SkuId == 0 || item.Quantity <= 0 {
			return nil, nil, errs.New(errs.ErrParam.Code, "invalid sku_id or quantity")
		}
		qtyMap[item.SkuId] += item.Quantity
		if _, ok := seen[item.SkuId]; !ok {
			seen[item.SkuId] = struct{}{}
			order = append(order, item.SkuId)
		}
	}

	return order, qtyMap, nil
}

func wrapRPC(err error, msg string) error {
	if err == nil {
		return nil
	}
	if bizErr, ok := kerrors.FromBizStatusError(err); ok {
		code := errs.ErrorType(bizErr.BizStatusCode())
		return errs.New(code, bizErr.BizMessage())
	}
	if e, ok := err.(*errs.Error); ok {
		return e
	}
	return errs.New(errs.ErrInternal.Code, msg+": "+err.Error())
}
