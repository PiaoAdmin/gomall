package main

import (
	"context"

	"github.com/PiaoAdmin/pmall/app/product/biz/service"
	product "github.com/PiaoAdmin/pmall/rpc_gen/product"
)

// ProductServiceImpl implements the last service interface defined in the IDL.
type ProductServiceImpl struct{}

// CreateProduct implements the ProductServiceImpl interface.
func (s *ProductServiceImpl) CreateProduct(ctx context.Context, req *product.CreateProductRequest) (resp *product.CreateProductResponse, err error) {
	resp, err = service.NewCreateProductService(ctx).Run(req)
	return
}

// UpdateProduct implements the ProductServiceImpl interface.
func (s *ProductServiceImpl) UpdateProduct(ctx context.Context, req *product.UpdateProductRequest) (resp *product.UpdateProductResponse, err error) {
	resp, err = service.NewUpdateProductService(ctx).Run(req)
	return
}

// UpdateProductStatus implements the ProductServiceImpl interface.
func (s *ProductServiceImpl) UpdateProductStatus(ctx context.Context, req *product.UpdateProductStatusRequest) (resp *product.UpdateProductStatusResponse, err error) {
	resp, err = service.NewUpdateProductStatusService(ctx).Run(req)
	return
}

// GetProductDetail implements the ProductServiceImpl interface.
func (s *ProductServiceImpl) GetProductDetail(ctx context.Context, req *product.GetProductDetailRequest) (resp *product.GetProductDetailResponse, err error) {
	resp, err = service.NewGetProductDetailService(ctx).Run(req)
	return
}

// ListProducts implements the ProductServiceImpl interface.
func (s *ProductServiceImpl) ListProducts(ctx context.Context, req *product.ListProductsRequest) (resp *product.ListProductsResponse, err error) {
	resp, err = service.NewListProductsService(ctx).Run(req)
	return
}

// GetProductsByIds implements the ProductServiceImpl interface.
func (s *ProductServiceImpl) GetProductsByIds(ctx context.Context, req *product.GetProductsByIdsRequest) (resp *product.GetProductsByIdsResponse, err error) {
	resp, err = service.NewGetProductsByIdsService(ctx).Run(req)
	return
}

// BatchUpdateSku implements the ProductServiceImpl interface.
func (s *ProductServiceImpl) BatchUpdateSku(ctx context.Context, req *product.BatchUpdateSkuRequest) (resp *product.BatchUpdateSkuResponse, err error) {
	resp, err = service.NewBatchUpdateSkuService(ctx).Run(req)
	return
}

// DeductStock implements the ProductServiceImpl interface.
func (s *ProductServiceImpl) DeductStock(ctx context.Context, req *product.DeductStockRequest) (resp *product.DeductStockResponse, err error) {
	resp, err = service.NewDeductStockService(ctx).Run(req)
	return
}

// ReleaseStock implements the ProductServiceImpl interface.
func (s *ProductServiceImpl) ReleaseStock(ctx context.Context, req *product.ReleaseStockRequest) (resp *product.ReleaseStockResponse, err error) {
	resp, err = service.NewReleaseStockService(ctx).Run(req)
	return
}

// ListCategories implements the ProductServiceImpl interface.
func (s *ProductServiceImpl) ListCategories(ctx context.Context, req *product.ListCategoriesRequest) (resp *product.ListCategoriesResponse, err error) {
	resp, err = service.NewListCategoriesService(ctx).Run(req)
	return
}

// ListBrands implements the ProductServiceImpl interface.
func (s *ProductServiceImpl) ListBrands(ctx context.Context, req *product.ListBrandsRequest) (resp *product.ListBrandsResponse, err error) {
	resp, err = service.NewListBrandsService(ctx).Run(req)
	return
}

// SearchProducts implements the ProductServiceImpl interface.
func (s *ProductServiceImpl) SearchProducts(ctx context.Context, req *product.SearchProductsRequest) (resp *product.SearchProductsResponse, err error) {
	resp, err = service.NewSearchProductsService(ctx).Run(req)
	return
}

// GetSkusByIds implements the ProductServiceImpl interface.
func (s *ProductServiceImpl) GetSkusByIds(ctx context.Context, req *product.GetSkusByIdsRequest) (resp *product.GetSkusByIdsResponse, err error) {
	resp, err = service.NewGetSkusByIdsService(ctx).Run(req)
	return
}
