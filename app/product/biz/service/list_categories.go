package service

import (
	"context"

	"github.com/PiaoAdmin/pmall/app/product/biz/dal/mysql"
	"github.com/PiaoAdmin/pmall/app/product/biz/model"
	product "github.com/PiaoAdmin/pmall/rpc_gen/product"
)

type ListCategoriesService struct {
	ctx context.Context
}

func NewListCategoriesService(ctx context.Context) *ListCategoriesService {
	return &ListCategoriesService{ctx: ctx}
}

func (s *ListCategoriesService) Run(req *product.ListCategoriesRequest) (*product.ListCategoriesResponse, error) {
	var categories []*model.ProductCategory
	var err error

	if req.ParentId > 0 {
		categories, err = model.GetCategoriesByParentID(s.ctx, mysql.DB, req.ParentId)
	} else {
		categories, err = model.GetAllCategories(s.ctx, mysql.DB)
	}

	if err != nil {
		return nil, err
	}

	respCategories := make([]*product.Category, 0, len(categories))

	if req.ParentId == 0 {
		categoryMap := make(map[uint64]*product.Category)
		var rootCategories []*product.Category

		for _, cat := range categories {
			protoCategory := &product.Category{
				Id:       cat.ID,
				ParentId: cat.ParentID,
				Name:     cat.Name,
				Level:    int32(cat.Level),
				Icon:     cat.Icon,
				Unit:     cat.Unit,
				Sort:     int32(cat.Sort),
				Children: []*product.Category{},
			}
			categoryMap[cat.ID] = protoCategory

			if cat.ParentID == 0 {
				rootCategories = append(rootCategories, protoCategory)
			}
		}

		for _, cat := range categories {
			if cat.ParentID != 0 {
				if parent, ok := categoryMap[cat.ParentID]; ok {
					parent.Children = append(parent.Children, categoryMap[cat.ID])
				}
			}
		}

		respCategories = rootCategories
	} else {
		for _, cat := range categories {
			respCategories = append(respCategories, &product.Category{
				Id:       cat.ID,
				ParentId: cat.ParentID,
				Name:     cat.Name,
				Level:    int32(cat.Level),
				Icon:     cat.Icon,
				Unit:     cat.Unit,
				Sort:     int32(cat.Sort),
			})
		}
	}

	return &product.ListCategoriesResponse{
		Categories: respCategories,
	}, nil
}
