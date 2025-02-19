# \*\*\* Project

## introduce

- Use the [Kitex](https://github.com/cloudwego/kitex/) framework
- Generating the base code for unit tests.
- Provides basic config functions
- Provides the most basic MVC code hierarchy.

## Directory structure

| catalog     | introduce                                       |
| ----------- | ----------------------------------------------- |
| conf        | Configuration files                             |
| main.go     | Startup file                                    |
| handler.go  | Used for request processing return of response. |
| kitex_gen   | kitex generated code                            |
| biz/service | The actual business logic.                      |
| biz/dal     | Logic for operating the storage layer           |

## How to run

```shell
sh build.sh
sh output/bootstrap.sh
```

## 项目文档

概述
本项目文档主要描述了 douyin-gomall 项目中与订单处理相关的功能模块，包括设计目标、主要功能、模块结构和 API 说明等。

设计目标
提供用户订单的查询功能，以便用户可以查看自己在 douyin-gomall 商城中已下订单的详细信息。
保证订单数据的准确性和及时性，为用户提供良好的购物体验。
主要功能
用户可以查询自己的订单列表。
订单查询结果中包含订单 ID、用户 ID、订单总金额、用户货币类型、收货人信息以及订单中包含的商品列表和数量。
模块结构
service: 业务逻辑层，包含订单查询服务的实现。
dal/mysql: 数据访问层，定义了与 MySQL 数据库交互的方法。
model: 数据模型层，定义了订单相关的数据结构。
rpc_gen/kitex_gen/order: 生成的 RPC 代码，包含了与订单相关的服务定义和消息定义。
API 说明
ListOrder: 根据用户 ID 查询订单列表。
请求参数: UserId（int64）
返回结果: 订单列表，每个订单包含订单 ID、用户 ID、订单总金额、用户货币类型、收货人信息以及订单中包含的商品列表和数量。
注意事项
使用前需要修改 app\order\conf 中的配置文件，主要修改的是数据库等服务器的连接信息，包括端口和密码。
