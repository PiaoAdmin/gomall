/*
 * @Author: liaosijie
 * @Date: 2025-02-16 10:14:49
 * @Last Modified by:  liaosijie
 * @Last Modified time: 2025-02-16 13:14:49
 */

package order

import (
	"time"
)

type Order struct {
    OrderID     string `json:"order_id"`
    UserID      int    `json:"user_id"`
    UserCurrency string `json:"user_currency"`
    AddressID   int    `json:"address_id"`
    Email       string `json:"email"`
    CreatedAt   time.Time `json:"created_at"`
    Status      string `json:"status"`
}

type OrderItem struct {
    OrderItemID int     `json:"order_item_id"`
    OrderID     string  `json:"order_id"`
    ProductID   int     `json:"product_id"`
    Quantity    int     `json:"quantity"`
    Price       float64 `json:"price"`
}

type Address struct {
    AddressID    int    `json:"address_id"`
    UserID       int    `json:"user_id"`
    StreetAddress string `json:"street_address"`
    City         string `json:"city"`
    State        string `json:"state"`
    Country      string `json:"country"`
    ZipCode      int    `json:"zip_code"`
}
