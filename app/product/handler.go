package main

import (
	"context"
	"github.com/PiaoAdmin/gomall/app/product/biz/service"
	product "github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen/product"
)

// ProductCatalogServiceImpl implements the last service interface defined in the IDL.
type ProductCatalogServiceImpl struct{}

// ListProducts implements the ProductCatalogServiceImpl interface.
func (s *ProductCatalogServiceImpl) ListProducts(ctx context.Context, req *product.ListProductsReq) (resp *product.ListProductsResp, err error) {
	resp, err = service.NewListProductsService(ctx).Run(req)

	return resp, err
}

// GetProduct implements the ProductCatalogServiceImpl interface.
func (s *ProductCatalogServiceImpl) GetProduct(ctx context.Context, req *product.GetProductReq) (resp *product.GetProductResp, err error) {
	resp, err = service.NewGetProductService(ctx).Run(req)

	return resp, err
}

// SearchProducts implements the ProductCatalogServiceImpl interface.
func (s *ProductCatalogServiceImpl) SearchProducts(ctx context.Context, req *product.SearchProductsReq) (resp *product.SearchProductsResp, err error) {
	resp, err = service.NewSearchProductsService(ctx).Run(req)

	return resp, err
}

// AddProduct implements the ProductCatalogServiceImpl interface.
func (s *ProductCatalogServiceImpl) AddProduct(ctx context.Context, req *product.AddProductReq) (resp *product.AddProductResp, err error) {
	resp, err = service.NewAddProductService(ctx).Run(req)

	return resp, err
}

// AssociateProductWithCategory implements the ProductCatalogServiceImpl interface.
func (s *ProductCatalogServiceImpl) AssociateProductWithCategory(ctx context.Context, req *product.AssociateProductWithCategoryReq) (resp *product.AssociateProductWithCategoryResp, err error) {
	resp, err = service.NewAssociateProductWithCategoryService(ctx).Run(req)

	return resp, err
}

// AddCategory implements the ProductCatalogServiceImpl interface.
func (s *ProductCatalogServiceImpl) AddCategory(ctx context.Context, req *product.AddCategoryReq) (resp *product.AddCategoryResp, err error) {
	resp, err = service.NewAddCategoryService(ctx).Run(req)

	return resp, err
}
