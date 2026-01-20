package mongo

import (
	"context"
	"time"

	"github.com/PiaoAdmin/pmall/app/product/conf"
	"github.com/cloudwego/kitex/pkg/klog"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	DB *mongo.Database
)

func Init() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoConf := conf.GetConf().MongoDB
	clientOptions := options.Client().ApplyURI(mongoConf.URI)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		klog.Fatalf("Failed to connect to MongoDB: %v", err)
		panic(err)
	}

	// 测试连接
	err = client.Ping(ctx, nil)
	if err != nil {
		klog.Fatalf("Failed to ping MongoDB: %v", err)
		panic(err)
	}

	klog.Info("Successfully connected to MongoDB")

	DB = client.Database(mongoConf.Database)
}
