/*
 * @Author: liaosijie
 * @Date: 2025-02-16 10:15:42
 * @Last Modified by: liaosijie
 * @Last Modified time: 2025-02-16 17:55:16
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
     err := ctx.BindAndValidate(&req)
     if err != nil {
         ctx.JSON(http.StatusBadRequest, utils.H{"msg": "invalid request"})
         return
     }
 
     o := &order.Order{
         UserID:       int64(req.User_id),
         UserCurrency: req.User_currency,
         AddressID:    int64(req.Address.Address_id),
         Email:        req.Email,
         CreatedAt:    time.Now(),
         Status:       "pending",
         OrderItems:   req.Order_items,
     }
 
     err = h.service.PlaceOrder(c, o)
     if err != nil {
         ctx.JSON(http.StatusInternalServerError, utils.H{"msg": "failed to place order"})
         return
     }
 
     ctx.JSON(http.StatusOK, utils.H{"msg": "order placed successfully", "order_id": o.ID})
 }
 
 // ListOrder 处理查询订单请求
 func (h *OrderHandler) ListOrder(c context.Context, ctx *app.RequestContext) {
     userId := ctx.QueryInt64("user_id")
     orders, err := h.service.ListOrder(c, userId)
     if err != nil {
         ctx.JSON(http.StatusInternalServerError, utils.H{"msg": "failed to list orders"})
         return
     }
 
     ctx.JSON(http.StatusOK, orders)
 }
 
 // MarkOrderPaid 处理标记订单为已支付请求
 func (h *OrderHandler) MarkOrderPaid(c context.Context, ctx *app.RequestContext) {
     var req order.MarkOrderPaidReq
     err := ctx.BindAndValidate(&req)
     if err != nil {
         ctx.JSON(http.StatusBadRequest, utils.H{"msg": "invalid request"})
         return
     }
 
     err = h.service.MarkOrderPaid(c, req.Order_id)
     if err != nil {
         ctx.JSON(http.StatusInternalServerError, utils.H{"msg": "failed to mark order as paid"})
         return
     }
 
     ctx.JSON(http.StatusOK, order.MarkOrderPaidResp{})
 }