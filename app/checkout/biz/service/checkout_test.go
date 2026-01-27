package service

import (
	"context"
	"errors"
	"testing"

	"github.com/PiaoAdmin/pmall/app/checkout/biz/rpc"
	checkout "github.com/PiaoAdmin/pmall/rpc_gen/checkout"
	"github.com/PiaoAdmin/pmall/rpc_gen/order"
	"github.com/PiaoAdmin/pmall/rpc_gen/order/orderservice"
	"github.com/PiaoAdmin/pmall/rpc_gen/payment"
	"github.com/PiaoAdmin/pmall/rpc_gen/payment/paymentservice"
	"github.com/PiaoAdmin/pmall/rpc_gen/product"
	"github.com/PiaoAdmin/pmall/rpc_gen/product/productservice"
	"github.com/PiaoAdmin/pmall/rpc_gen/user"
	"github.com/PiaoAdmin/pmall/rpc_gen/user/userservice"
	"github.com/cloudwego/kitex/client/callopt"
)

var errNotImplemented = errors.New("not implemented")

type userClientMock struct {
	getResp *user.GetUserInfoResponse
	getErr  error
}

func (m *userClientMock) Register(ctx context.Context, Req *user.RegisterRequest, callOptions ...callopt.Option) (*user.RegisterResponse, error) {
	return nil, errNotImplemented
}

func (m *userClientMock) Login(ctx context.Context, Req *user.LoginRequest, callOptions ...callopt.Option) (*user.LoginResponse, error) {
	return nil, errNotImplemented
}

func (m *userClientMock) GetUserInfo(ctx context.Context, Req *user.GetUserInfoRequest, callOptions ...callopt.Option) (*user.GetUserInfoResponse, error) {
	return m.getResp, m.getErr
}

func (m *userClientMock) UpdateUser(ctx context.Context, Req *user.UpdateUserRequest, callOptions ...callopt.Option) (*user.UpdateUserResponse, error) {
	return nil, errNotImplemented
}

func (m *userClientMock) UpdatePassword(ctx context.Context, Req *user.UpdatePasswordRequest, callOptions ...callopt.Option) (*user.UpdatePasswordResponse, error) {
	return nil, errNotImplemented
}

type productClientMock struct {
	getResp     *product.GetSkusByIdsResponse
	getErr      error
	deductResp  *product.DeductStockResponse
	deductErr   error
	releaseResp *product.ReleaseStockResponse
	releaseErr  error

	deductCalled  bool
	releaseCalled bool
	releaseReq    *product.ReleaseStockRequest
}

func (m *productClientMock) CreateProduct(ctx context.Context, Req *product.CreateProductRequest, callOptions ...callopt.Option) (*product.CreateProductResponse, error) {
	return nil, errNotImplemented
}

func (m *productClientMock) UpdateProduct(ctx context.Context, Req *product.UpdateProductRequest, callOptions ...callopt.Option) (*product.UpdateProductResponse, error) {
	return nil, errNotImplemented
}

func (m *productClientMock) UpdateProductStatus(ctx context.Context, Req *product.UpdateProductStatusRequest, callOptions ...callopt.Option) (*product.UpdateProductStatusResponse, error) {
	return nil, errNotImplemented
}

func (m *productClientMock) GetProductDetail(ctx context.Context, Req *product.GetProductDetailRequest, callOptions ...callopt.Option) (*product.GetProductDetailResponse, error) {
	return nil, errNotImplemented
}

func (m *productClientMock) ListProducts(ctx context.Context, Req *product.ListProductsRequest, callOptions ...callopt.Option) (*product.ListProductsResponse, error) {
	return nil, errNotImplemented
}

func (m *productClientMock) GetProductsByIds(ctx context.Context, Req *product.GetProductsByIdsRequest, callOptions ...callopt.Option) (*product.GetProductsByIdsResponse, error) {
	return nil, errNotImplemented
}

func (m *productClientMock) GetSkusByIds(ctx context.Context, Req *product.GetSkusByIdsRequest, callOptions ...callopt.Option) (*product.GetSkusByIdsResponse, error) {
	return m.getResp, m.getErr
}

func (m *productClientMock) BatchUpdateSku(ctx context.Context, Req *product.BatchUpdateSkuRequest, callOptions ...callopt.Option) (*product.BatchUpdateSkuResponse, error) {
	return nil, errNotImplemented
}

func (m *productClientMock) DeductStock(ctx context.Context, Req *product.DeductStockRequest, callOptions ...callopt.Option) (*product.DeductStockResponse, error) {
	m.deductCalled = true
	return m.deductResp, m.deductErr
}

func (m *productClientMock) ReleaseStock(ctx context.Context, Req *product.ReleaseStockRequest, callOptions ...callopt.Option) (*product.ReleaseStockResponse, error) {
	m.releaseCalled = true
	m.releaseReq = Req
	return m.releaseResp, m.releaseErr
}

func (m *productClientMock) ListCategories(ctx context.Context, Req *product.ListCategoriesRequest, callOptions ...callopt.Option) (*product.ListCategoriesResponse, error) {
	return nil, errNotImplemented
}

func (m *productClientMock) ListBrands(ctx context.Context, Req *product.ListBrandsRequest, callOptions ...callopt.Option) (*product.ListBrandsResponse, error) {
	return nil, errNotImplemented
}

func (m *productClientMock) SearchProducts(ctx context.Context, Req *product.SearchProductsRequest, callOptions ...callopt.Option) (*product.SearchProductsResponse, error) {
	return nil, errNotImplemented
}

type orderClientMock struct {
	placeResp *order.PlaceOrderResp
	placeErr  error

	cancelCalled bool
	cancelReq    *order.CancelOrderReq

	markPaidCalled bool
}

func (m *orderClientMock) ListOrder(ctx context.Context, Req *order.ListOrderReq, callOptions ...callopt.Option) (*order.ListOrderResp, error) {
	return nil, errNotImplemented
}

func (m *orderClientMock) CancelOrder(ctx context.Context, Req *order.CancelOrderReq, callOptions ...callopt.Option) (*order.CancelOrderResp, error) {
	m.cancelCalled = true
	m.cancelReq = Req
	return &order.CancelOrderResp{Success: true}, nil
}

func (m *orderClientMock) PlaceOrder(ctx context.Context, Req *order.PlaceOrderReq, callOptions ...callopt.Option) (*order.PlaceOrderResp, error) {
	return m.placeResp, m.placeErr
}

func (m *orderClientMock) MarkOrderPaid(ctx context.Context, Req *order.MarkOrderPaidReq, callOptions ...callopt.Option) (*order.MarkOrderPaidResp, error) {
	m.markPaidCalled = true
	return &order.MarkOrderPaidResp{Success: true}, nil
}

type paymentClientMock struct {
	payResp *payment.PayResponse
	payErr  error
}

func (m *paymentClientMock) Pay(ctx context.Context, Req *payment.PayRequest, callOptions ...callopt.Option) (*payment.PayResponse, error) {
	return m.payResp, m.payErr
}

func TestCheckoutServiceSuccess(t *testing.T) {
	oldUser := rpc.UserClient
	oldProduct := rpc.ProductClient
	oldOrder := rpc.OrderClient
	oldPayment := rpc.PaymentClient
	defer func() {
		rpc.UserClient = oldUser
		rpc.ProductClient = oldProduct
		rpc.OrderClient = oldOrder
		rpc.PaymentClient = oldPayment
	}()

	rpc.UserClient = &userClientMock{getResp: &user.GetUserInfoResponse{User: &user.User{Id: 1, Email: "a@b.com"}}}
	rpc.ProductClient = &productClientMock{getResp: &product.GetSkusByIdsResponse{Skus: map[uint64]*product.ProductSKU{
		1001: {Id: 1001, Name: "sku-1", Price: "10.50", Stock: 10, MainImage: "img", MarketPrice: "12.00", SpuId: 2001},
	}}, deductResp: &product.DeductStockResponse{Success: true}}
	rpc.OrderClient = &orderClientMock{placeResp: &order.PlaceOrderResp{Order: &order.OrderResult{OrderId: "oid-1"}}}
	rpc.PaymentClient = &paymentClientMock{payResp: &payment.PayResponse{Success: true, TradeNo: "t1"}}

	svc := NewCheckoutService(context.Background())
	resp, err := svc.Run(&checkout.CheckoutRequest{
		UserId: 1,
		Items:  []*checkout.CheckoutItem{{SkuId: 1001, Quantity: 2}},
		ShippingAddress: &checkout.Address{
			StreetAddress: "street",
			City:          "city",
			Name:          "name",
			ZipCode:       123,
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.OrderId != "oid-1" {
		t.Fatalf("unexpected order id: %s", resp.OrderId)
	}
	if resp.TotalAmount != "21.00" {
		t.Fatalf("unexpected total amount: %s", resp.TotalAmount)
	}

	orderMock := rpc.OrderClient.(*orderClientMock)
	if !orderMock.markPaidCalled {
		t.Fatalf("mark order paid not called")
	}
}

func TestCheckoutServicePaymentFail(t *testing.T) {
	oldUser := rpc.UserClient
	oldProduct := rpc.ProductClient
	oldOrder := rpc.OrderClient
	oldPayment := rpc.PaymentClient
	defer func() {
		rpc.UserClient = oldUser
		rpc.ProductClient = oldProduct
		rpc.OrderClient = oldOrder
		rpc.PaymentClient = oldPayment
	}()

	prodMock := &productClientMock{
		getResp:     &product.GetSkusByIdsResponse{Skus: map[uint64]*product.ProductSKU{1001: {Id: 1001, Name: "sku-1", Price: "10.00", Stock: 10}}},
		deductResp:  &product.DeductStockResponse{Success: true},
		releaseResp: &product.ReleaseStockResponse{Success: true},
	}
	orderMock := &orderClientMock{placeResp: &order.PlaceOrderResp{Order: &order.OrderResult{OrderId: "oid-2"}}}

	rpc.UserClient = &userClientMock{getResp: &user.GetUserInfoResponse{User: &user.User{Id: 2, Email: "b@c.com"}}}
	rpc.ProductClient = prodMock
	rpc.OrderClient = orderMock
	rpc.PaymentClient = &paymentClientMock{payResp: &payment.PayResponse{Success: false}}

	svc := NewCheckoutService(context.Background())
	_, err := svc.Run(&checkout.CheckoutRequest{
		UserId:          2,
		Items:           []*checkout.CheckoutItem{{SkuId: 1001, Quantity: 1}},
		ShippingAddress: &checkout.Address{StreetAddress: "street", City: "city", Name: "name", ZipCode: 123},
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !prodMock.releaseCalled {
		t.Fatalf("release stock not called")
	}
	if prodMock.releaseReq == nil || prodMock.releaseReq.OrderSn != "oid-2" {
		t.Fatalf("release stock order_sn not set")
	}
	if !orderMock.cancelCalled {
		t.Fatalf("cancel order not called")
	}
}

var _ userservice.Client = (*userClientMock)(nil)
var _ productservice.Client = (*productClientMock)(nil)
var _ orderservice.Client = (*orderClientMock)(nil)
var _ paymentservice.Client = (*paymentClientMock)(nil)
