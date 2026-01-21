CREATE DATABASE IF NOT EXISTS `p_order`
    DEFAULT CHARACTER SET = 'utf8mb4';
USE `p_order`;

-- ----------------------------
-- 1. 订单主表 (orders)
-- ----------------------------
DROP TABLE IF EXISTS `orders`;
CREATE TABLE `orders` (
  `id` bigint NOT NULL COMMENT '订单ID',
  `order_id` varchar(64) NOT NULL COMMENT '订单号，对应 proto order_id',
  `user_id` bigint unsigned NOT NULL COMMENT '用户ID，对应 proto user_id',
  `email` varchar(255) NOT NULL DEFAULT '' COMMENT '邮箱，对应 proto email',
  `shipping_name` varchar(64) NOT NULL DEFAULT '' COMMENT '对应 proto Address.name',
  `shipping_street_address` varchar(255) NOT NULL DEFAULT '' COMMENT '对应 proto Address.street_address',
  `shipping_city` varchar(64) NOT NULL DEFAULT '' COMMENT '对应 proto Address.city',
  `shipping_zip_code` int NOT NULL DEFAULT 0 COMMENT '对应 proto Address.zip_code',
  `status` varchar(32) NOT NULL DEFAULT '' COMMENT '状态，对应 proto status',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  `is_deleted` tinyint DEFAULT '0' COMMENT '逻辑删除标记:0-未删除,1-已删除',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_order_id` (`order_id`),
  KEY `idx_orders_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- 2. 订单项表 (order_items)
-- ----------------------------
DROP TABLE IF EXISTS `order_items`;
CREATE TABLE `order_items` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `order_id` varchar(64) NOT NULL COMMENT '外键，关联 orders 表',
  `sku_id` bigint unsigned NOT NULL COMMENT '对应 proto CartItem.sku_id',
  `sku_name` varchar(255) NOT NULL DEFAULT '' COMMENT '对应 proto CartItem.sku_name',
  `price` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '对应 proto CartItem.price',
  `quantity` int NOT NULL DEFAULT 1 COMMENT '对应 proto CartItem.quantity',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  `is_deleted` tinyint DEFAULT '0' COMMENT '逻辑删除标记:0-未删除,1-已删除',
  PRIMARY KEY (`id`),
  KEY `idx_order_items_order_id` (`order_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;