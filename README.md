# P-Mall

一个基于 Go 微服务架构的电商系统。

## 技术栈

- Go 1.25+
- [CloudWeGo Hertz](https://github.com/cloudwego/hertz) - HTTP 框架
- [CloudWeGo Kitex](https://github.com/cloudwego/kitex) - RPC 框架
- gRPC + Protobuf
- Consul - 服务注册与发现
- MySQL + Redis + MongoDB
- Python + LangGraph - AI Agent

## 项目结构

```
├── app/                    # 微服务
│   ├── api/               # API Gateway (Hertz)
│   ├── user/              # 用户服务 (Kitex)
│   ├── product/           # 商品服务 (Kitex)
│   ├── cart/              # 购物车服务 (Kitex)
│   ├── order/             # 订单服务 (Kitex)
│   ├── checkout/          # 结算服务 (Kitex)
│   └── payment/           # 支付服务 (Kitex)
├── agent/                 # AI Agent
│   ├── product_listing_agent/  # 商品上架 Agent
│   └── auto_order_agent/       # 自动下单 Agent
├── idl/                   # Protobuf IDL 定义
├── rpc_gen/               # RPC 代码生成
├── common/                # 公共库
├── sql/                   # 数据库脚本
└── scripts/               # 构建脚本
```

## 快速开始

### 前置要求

- Go 1.25+
- MySQL
- Redis
- Consul
- Python 3.9+ (可选，用于 AI Agent)

### 启动服务

1. 启动基础设施（MySQL、Redis、Consul）

2. 初始化数据库
   ```bash
   mysql -u root -p < sql/user.sql
   mysql -u root -p < sql/product.sql
   mysql -u root -p < sql/order.sql
   ```

3. 启动微服务
   ```bash
   # 启动所有服务
   bash scripts/start_services.sh
   
   # 或单独启动每个服务
   cd app/user && sh build.sh && ./output/bootstrap.sh
   cd app/product && sh build.sh && ./output/bootstrap.sh
   # ...以此类推
   ```

4. 启动 API Gateway
   ```bash
   cd app/api && sh build.sh && ./output/bootstrap.sh
   ```

5. 访问 API 文档
   ```
   http://localhost:8888/swagger/index.html
   ```

### AI Agent (可选)

```bash
cd agent

# 商品上架 Agent
bash run_product_agent.sh

# 自动下单 Agent
cd auto_order_agent && python main.py
```

## 服务端口

| 服务 | 端口 |
|------|------|
| API Gateway | 8888 |
| User Service | 9900 |
| Product Service | 9901 |
| Order Service | 9902 |
| Checkout Service | 9903 |
| Cart Service | 9904 |
| Payment Service | 9905 |

## 开发

### 代码生成

```bash
# 生成 RPC 代码
cd scripts
make gen-rpc

# 生成 API 代码
cd app/api
hz update -idl ../../idl/api/*.proto
```

### 构建

```bash
# 构建所有服务
cd scripts
make build

# 构建单个服务
cd app/user
sh build.sh
```

## License

MIT
