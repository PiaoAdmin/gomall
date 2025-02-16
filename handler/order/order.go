/*
 * @Author: liaosijie
 * @Date: 2025-02-16 10:15:42
 * @Last Modified by:   liaosijie
 * @Last Modified time: 2025-02-16 13:15:42
 */

package order

import (
	"context"
	"gomall/biz/model/order"
	"gomall/biz/service/order"
	"net/http"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
)

type OrderHandler struct {
    service *order.OrderService
}

func NewOrderHandler(service *order.OrderService) *OrderHandler {
    return &OrderHandler{service: service}
}

// PlaceOrder 处理创建订单请求
func (h *OrderHandler) PlaceOrder(c context.Context, ctx *app.RequestContext) {
    var req order.PlaceOrderReq
    if err := ctx.BindAndValidate(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, utils.H{"msg": "invalid request"})
        return
    }

    o := &order.Order{
        UserID:      int(req.UserID),
        UserCurrency: req.UserCurrency,
        AddressID:   req.Address.AddressID,
        Email:       req.Email,
        CreatedAt:   time.Now(),
        Status:      "pending",
        OrderItems:  req.OrderItems,
    }

    if err := h.service.PlaceOrder(c, o); err != nil {
        ctx.JSON(http.StatusInternalServerError, utils.H{"msg": "failed to place order"})
        return
    }

    resp := order.PlaceOrderResp{
        OrderResult: order.OrderResult{
            OrderID: o.OrderID,
        },
    }

    ctx.JSON(http.StatusOK, resp)
}

// ListOrder 处理查询订单请求
func (h *OrderHandler) ListOrder(c context.Context, ctx *app.RequestContext) {
    userId := ctx.QueryInt("user_id")
    orders, err := h.service.ListOrder(c, userId)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, utils.H{"msg": "failed to list orders"})
        return
    }

    resp := order.ListOrderResp{
        Orders: orders,
    }

    ctx.JSON(http.StatusOK, resp)
}

// MarkOrderPaid 处理标记订单为已支付请求
func (h *OrderHandler) MarkOrderPaid(c context.Context, ctx *app.RequestContext) {
    var req order.MarkOrderPaidReq
    if err := ctx.BindAndValidate(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, utils.H{"msg": "invalid request"})
        return
    }

    if err := h.service.MarkOrderPaid(c, req.OrderID); err != nil {
        ctx.JSON(http.StatusInternalServerError, utils.H{"msg": "failed to mark order as paid"})
        return
    }

    ctx.JSON(http.StatusOK, order.MarkOrderPaidResp{})
}
