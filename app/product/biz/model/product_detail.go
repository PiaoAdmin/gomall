package model

import (
	"context"
	"time"

	"github.com/PiaoAdmin/pmall/app/product/biz/dal/mongo"
	"github.com/cloudwego/kitex/pkg/klog"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	CollectionProductDetails = "product_details"
)

// ProductDetail 商品详情 - 存储在MongoDB中
type ProductDetail struct {
	SpuID         uint64   `bson:"spu_id" json:"spu_id"`
	Description   string   `bson:"description" json:"description"`                   // 富文本描述
	Images        []string `bson:"images" json:"images"`                             // 商品图
	Videos        []string `bson:"videos,omitempty" json:"videos"`                   // 商品视频
	MarketTagJSON string   `bson:"market_tag_json,omitempty" json:"market_tag_json"` // 营销标签 JSON 字符串
	TechTagJSON   string   `bson:"tech_tag_json,omitempty" json:"tech_tag_json"`     // 技术参数 JSON 字符串
	FaqJSON       string   `bson:"faq_json,omitempty" json:"faq_json"`               // 常见问题 JSON 字符串
	CreatedAt     int64    `bson:"created_at" json:"created_at"`                     // Unix 时间戳
	UpdatedAt     int64    `bson:"updated_at" json:"updated_at"`                     // Unix 时间戳
}

func getProductDetailCollection() *mongodriver.Collection {
	return mongo.DB.Collection(CollectionProductDetails)
}

// 创建 spu_id 唯一索引
func InitProductDetailIndexes() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := getProductDetailCollection()

	indexModel := mongodriver.IndexModel{
		Keys:    bson.D{{Key: "spu_id", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	_, err := collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		klog.Warnf("Failed to create spu_id index for product_details: %v", err)
	} else {
		klog.Info("Created spu_id index for product_details")
	}
}

func CreateProductDetail(ctx context.Context, detail *ProductDetail) error {
	detail.CreatedAt = time.Now().Unix()
	detail.UpdatedAt = time.Now().Unix()
	collection := getProductDetailCollection()
	_, err := collection.InsertOne(ctx, detail)
	return err
}

// UpdateProductDetail 增量更新商品详情，只更新 updates map 中指定的字段
func UpdateProductDetail(ctx context.Context, spuID uint64, updates map[string]interface{}) (int64, error) {
	if len(updates) == 0 {
		return 0, nil
	}

	updates["updated_at"] = time.Now().Unix()

	collection := getProductDetailCollection()
	filter := bson.M{"spu_id": spuID}
	update := bson.M{"$set": updates}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return 0, err
	}
	return result.ModifiedCount, nil
}

func GetProductDetailBySpuID(ctx context.Context, spuID uint64) (*ProductDetail, error) {
	collection := getProductDetailCollection()
	filter := bson.M{"spu_id": spuID}

	var detail ProductDetail
	err := collection.FindOne(ctx, filter).Decode(&detail)
	if err != nil {
		if err == mongodriver.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &detail, nil
}
