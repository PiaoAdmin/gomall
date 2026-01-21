package model

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/PiaoAdmin/pmall/app/cart/biz/dal/redis"
	"github.com/cloudwego/kitex/pkg/klog"
)

const (
	CartKeyPrefix  = "cart:user:"
	CartExpiration = 7 * 24 * time.Hour
)

func GetCartKey(userID uint64) string {
	return fmt.Sprintf("%s%d", CartKeyPrefix, userID)
}

func AddToCart(ctx context.Context, userID, skuID uint64, quantity int32) error {
	key := GetCartKey(userID)
	field := strconv.FormatUint(skuID, 10)

	err := redis.RedisClient.HIncrBy(ctx, key, field, int64(quantity)).Err()
	if err != nil {
		return err
	}

	redis.RedisClient.Expire(ctx, key, CartExpiration)

	return nil
}

func RemoveFromCart(ctx context.Context, userID uint64, skuIDs []uint64) error {
	removeCount := make(map[uint64]int64)
	for _, skuID := range skuIDs {
		removeCount[skuID]++
	}

	klog.CtxInfof(ctx, "RemoveFromCart: userID=%d, removeCount=%+v", userID, removeCount)
	key := GetCartKey(userID)
	pipe := redis.RedisClient.Pipeline()
	for skuID, count := range removeCount {
		field := strconv.FormatUint(skuID, 10)
		pipe.HIncrBy(ctx, key, field, -count)
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		return err
	}
	return cleanUpZeroQuantityItems(ctx, userID)
}

func cleanUpZeroQuantityItems(ctx context.Context, userID uint64) error {
	key := GetCartKey(userID)

	items, err := redis.RedisClient.HGetAll(ctx, key).Result()
	if err != nil {
		return err
	}

	var toDelete []string
	for field, qtyStr := range items {
		qty, _ := strconv.ParseInt(qtyStr, 10, 64)
		if qty <= 0 {
			toDelete = append(toDelete, field)
		}
	}

	if len(toDelete) > 0 {
		return redis.RedisClient.HDel(ctx, key, toDelete...).Err()
	}

	return nil
}

func GetCartItems(ctx context.Context, userID uint64) (map[uint64]int32, error) {
	key := GetCartKey(userID)

	result, err := redis.RedisClient.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	items := make(map[uint64]int32)
	for skuIDStr, quantityStr := range result {
		skuID, _ := strconv.ParseUint(skuIDStr, 10, 64)
		quantity, _ := strconv.ParseInt(quantityStr, 10, 32)
		items[skuID] = int32(quantity)
	}

	return items, nil
}

func GetCartItemQuantity(ctx context.Context, userID, skuID uint64) (int32, error) {
	key := GetCartKey(userID)
	field := strconv.FormatUint(skuID, 10)

	quantityStr, err := redis.RedisClient.HGet(ctx, key, field).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			return 0, nil
		}
		return 0, err
	}

	quantity, _ := strconv.ParseInt(quantityStr, 10, 32)
	return int32(quantity), nil
}

func ClearCart(ctx context.Context, userID uint64) error {
	key := GetCartKey(userID)
	return redis.RedisClient.Del(ctx, key).Err()
}

func UpdateCartItemQuantity(ctx context.Context, userID, skuID uint64, quantity int32) error {
	key := GetCartKey(userID)
	field := strconv.FormatUint(skuID, 10)

	if quantity <= 0 {
		return redis.RedisClient.HDel(ctx, key, field).Err()
	}

	err := redis.RedisClient.HSet(ctx, key, field, quantity).Err()
	if err != nil {
		return err
	}

	redis.RedisClient.Expire(ctx, key, CartExpiration)

	return nil
}
