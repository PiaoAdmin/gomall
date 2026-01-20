DROP DATABASE IF EXISTS `p_product`;
CREATE DATABASE `p_product` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
USE `p_product`;

DROP TABLE IF EXISTS `product_spu`;

CREATE TABLE `product_spu` (
  `id` bigint NOT NULL COMMENT '商品ID',
  `brand_id` bigint NOT NULL COMMENT '品牌ID',
  `category_id` bigint NOT NULL COMMENT '分类ID',
  `name` varchar(255) NOT NULL COMMENT '商品名称',
  `sub_title` varchar(500) DEFAULT NULL COMMENT '商品副标题',
  `main_image` varchar(1000) DEFAULT NULL COMMENT '商品主图',

  `publish_status` tinyint DEFAULT '0' COMMENT '商品发布状态:0-未发布,1-已发布',
  `verify_status` tinyint DEFAULT '0' COMMENT '商品审核状态:0-未审核,1-审核通过,2-审核不通过',

  `low_price` decimal(10,2) DEFAULT NULL COMMENT '最低售价',
  `high_price` decimal(10,2) DEFAULT NULL COMMENT '最高售价',

  `sale_count` int DEFAULT '0' COMMENT '销量',
  `sort` int DEFAULT '0' COMMENT '排序权重',

  `service_bits` BIGINT(20) DEFAULT '0' COMMENT '商品服务:用二进制位存储,每一位代表一种服务',
  `version` int DEFAULT '1' COMMENT '乐观锁版本号',

  `create_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `delete_at` datetime DEFAULT NULL COMMENT '删除时间',
  `is_deleted` tinyint DEFAULT '0' COMMENT '逻辑删除标记:0-未删除,1-已删除',
  PRIMARY KEY (`id`),
  KEY `idx_cat_brand` (`category_id`, `brand_id`), -- 联合索引优化筛选
  KEY `idx_status` (`is_deleted`, `publish_status`, `verify_status`),
  KEY `idx_sale_count` (`sale_count` DESC),
  KEY `idx_sort` (`sort` DESC),
  KEY `idx_brand_id` (`brand_id`),
  KEY `idx_is_deleted` (`is_deleted`),
  KEY `idx_price_range` (`low_price`, `high_price`)
) ENGINE=InnoDB COMMENT='商品SPU表';

DROP TABLE IF EXISTS `product_sku`;

CREATE TABLE `product_sku` (
  `id` bigint NOT NULL COMMENT '商品SKU ID',
  `spu_id` bigint NOT NULL COMMENT '商品SPU ID',
  `sku_code` varchar(64) NOT NULL COMMENT '商家内部SKU编码',
  `name` varchar(255) NOT NULL COMMENT 'SKU名称',
  `sub_title` varchar(500) DEFAULT NULL COMMENT 'SKU副标题',
  `main_image` varchar(1000) DEFAULT NULL COMMENT 'SKU商品主图',

  `publish_status` tinyint DEFAULT '0' COMMENT '商品发布状态:0-未发布,1-已发布',
  `verify_status` tinyint DEFAULT '0' COMMENT '商品审核状态:0-未审核,1-审核通过,2-审核不通过',

  `price` decimal(10,2) NOT NULL COMMENT '销售价',
  `market_price` decimal(10,2) DEFAULT NULL COMMENT '市场价',
  `stock` int DEFAULT '0' COMMENT '库存数量',
  `lock_stock` int DEFAULT '0' COMMENT '锁定库存(下单未付)',

  `sku_spec_data` json DEFAULT NULL COMMENT '规格键值对',

  `version` int DEFAULT '1' COMMENT '乐观锁版本号',

  `create_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `delete_at` datetime DEFAULT NULL COMMENT '删除时间',
  `is_deleted` tinyint DEFAULT '0' COMMENT '逻辑删除标记:0-未删除,1-已删除',
  PRIMARY KEY (`id`),
  KEY `idx_spu_id` (`spu_id`, `is_deleted`, `publish_status`),
  KEY `idx_status` (`is_deleted`, `publish_status`, `verify_status`),
  UNIQUE KEY `uk_sku_code` (`sku_code`),
  KEY `idx_stock` (`stock`),
  KEY `idx_price` (`price`),
  KEY `idx_is_deleted` (`is_deleted`)
) ENGINE=InnoDB COMMENT='商品SKU表';

DROP TABLE IF EXISTS `product_category`;

CREATE TABLE `product_category` (
  `id` bigint NOT NULL COMMENT '分类ID',
  `parent_id` bigint DEFAULT '0' COMMENT '父分类ID',
  `name` varchar(64) NOT NULL COMMENT '分类名称',
  `level` tinyint DEFAULT '1' COMMENT '层级: 1/2/3',
  `icon` varchar(500) DEFAULT NULL COMMENT '图标URL',
  `unit` varchar(32) DEFAULT NULL COMMENT '数量单位(件/台)',
  `sort` int(11) DEFAULT '0',
  `show_status` tinyint DEFAULT '1' COMMENT '是否显示',

  `create_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `delete_at` datetime DEFAULT NULL COMMENT '删除时间',
  `is_deleted` tinyint DEFAULT '0' COMMENT '逻辑删除标记:0-未删除,1-已删除',
  PRIMARY KEY (`id`),
  KEY `idx_parent_id` (`parent_id`, `sort`),
  KEY `idx_show_status` (`show_status`, `sort`),
  KEY `idx_level` (`level`)
) ENGINE=InnoDB COMMENT='商品分类表';

DROP TABLE IF EXISTS `product_brand`;

CREATE TABLE `product_brand` (
  `id` bigint NOT NULL COMMENT '品牌ID',
  `name` varchar(64) NOT NULL UNIQUE COMMENT '品牌名称',
  `first_letter` char(3) DEFAULT NULL COMMENT '品牌首字母',
  `logo` varchar(500) DEFAULT NULL COMMENT '品牌LOGO URL',
  `sort` int(11) DEFAULT '0',
  `show_status` tinyint DEFAULT '1' COMMENT '是否显示',

  `create_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `delete_at` datetime DEFAULT NULL COMMENT '删除时间',
  `is_deleted` tinyint DEFAULT '0' COMMENT '逻辑删除标记:0-未删除,1-已删除',
  PRIMARY KEY (`id`),
  KEY `idx_first_letter` (`first_letter`, `sort`),
  KEY `idx_show_status` (`show_status`, `sort`)
) ENGINE=InnoDB COMMENT='商品品牌表';