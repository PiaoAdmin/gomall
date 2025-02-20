/*
 * @Author: liaosijie
 * @Date: 2025-02-20 15:14:27
 * @Last Modified by: liaosijie
 * @Last Modified time: 2025-02-20 20:47:58
 */

package service

import (
	"context"
	"time"

	"github.com/PiaoAdmin/gomall/app/hertz_gateway/infra/rpc"
	"github.com/PiaoAdmin/gomall/app/hertz_gateway/types"

	// common "douyin-gomall/gomall/app/hertz_gateway/hertz_gen/hertz_gateway/common"
	common "github.com/PiaoAdmin/gomall/app/hertz_gateway/hertz_gen/hertz_gateway/common"
	hertz_gatewayUtils "github.com/PiaoAdmin/gomall/app/hertz_gateway/utils"
	order "github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/order"
	product "github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/product"

	// "github.com/cloudwego/configmanager/configvalue/items" should be imported if you use configmanager
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
)

type OrderListService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewOrderListService(Context context.Context, RequestContext *app.RequestContext) *OrderListService {
	return &OrderListService{RequestContext: RequestContext, Context: Context}
}

func (h *OrderListService) Run(req *common.Empty) (resp map[string]any, err error) {
	//defer func() {
	// hlog.CtxInfof(h.Context, "req = %+v", req)
	// hlog.CtxInfof(h.Context, "resp = %+v", resp)
	//}()
	// todo edit your code

	userId := hertz_gatewayUtils.GetUserIdFromCtx(h.Context)
	orderResp, err := rpc.OrderClient.ListOrder(h.Context, &order.ListOrderReq{UserId: userId})
	if err != nil {
		return nil, err
	}
	
	var list []types.Order
	for _, v := range orderResp.Orders {
		var (
			total float32
			items []types.OrderItem
		)

		for _, vv := range v.Items {
			total += vv.Cost
			i := vv.Item
			productResp, err := rpc.ProductClient.GetProduct(h.Context, &product.GetProductReq{Id: i.ProductId})
			if err != nil {
				return nil, err
			}
			if productResp == nil || productResp.Product == nil {
				continue
			}
			p := productResp.Product
			items = append(items, types.OrderItem{
				ProductName: p.Name,
				Picture:     p.Picture,
				Cost:        vv.Cost,
				Qty:         i.Quantity,
			})
		}
		created := time.Unix(int64(v.CreatedAt), 0)
		list = append(list, types.Order{
			Cost:        total,
			Items:       items,
			CreatedDate: created.Format("2006-01-02 15:04:05"),
			OrderId:     v.OrderId,
			Consignee:   types.Consignee{Email: v.Email},
		})
	}
	return utils.H{
		"title":  "Order",
		"orders": list,
	},nil
}
