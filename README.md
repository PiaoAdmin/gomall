你需求写的不队抖音商城项目文档
一、项目介绍
概况：基于go语言实现的抖音商城微服务项目，具有完善的用户身份认证、商品管理、购物车、订单管理和支付等核心功能。
项目地址：https://github.com/PiaoAdmin/gomall/tree/dev
二、项目分工
团队成员	主要贡献
张浩辰	负责开发商品模块，处理git代码合并冲突
汤兵兵（队长）	负责开发用户模块和鉴权模块，整合网关和各个微服务之间的接口调用以及后续优化
杨国栋	负责开发购物车服务
卢嘉钦	负责开发支付服务和结算
廖思杰	负责开发订单服务
三、项目实现
3.1 技术选型与相关开发文档
可以补充场景分析环节，明确要解决的问题和前提假设，比如按当前选型和架构总体预计需要xxx存储空间，xxx台服务器......。
场景分析
本项目为抖音商城，需要支持大规模用户访问和高并发交易处理，核心挑战包括：
•          高并发与高可用：秒杀、促销等场景下，流量会瞬时激增，系统需具备良好的负载均衡和弹性扩展能力。
•          数据一致性：订单、库存、支付等关键业务需要保证事务一致性，避免超卖等问题。
•          低延迟要求：用户期望快速响应，API 接口需具备毫秒级响应能力，减少用户等待时间。
•          存储与计算需求：因为用户和订单数量大需要分布式存储，所以要保证存储时不出现主键冲突。
技术选型
基于上述需求，技术栈的选择遵循高并发、低延迟、可扩展性原则，选定如下方案：
•          开发语言：Go
•          微服务架构：Hertz + Kitex
￮        Hertz（HTTP 框架）：
▪         针对高并发 HTTP 请求进行了优化，使用其做为统一的网关
￮        Kitex（RPC 框架）：
▪         支持高效的 RPC 通信，减少跨服务调用开销，使用其做微服务
￮        Consul:
▪         服务发现和服务注册
•          消息队列：RocketMQ（用于削峰填谷，支持异步任务，如订单超时取消）
•          数据库：MySQL + Redis
￮        MySQL：存储核心业务数据，如用户、订单、支付信息
￮        Redis：用于缓存、热点数据存储、订单流水号生成等，减少数据库压力
•          容器化与运维
￮        Kubernetes（K8s）：管理微服务，提供弹性伸缩
￮        Prometheus + Grafana：监控系统性能，优化服务稳定性
3.2 架构设计
总体架构采用微服务架构，主要包括以下几个核心模块：
1.        认证服务（Auth Service）：负责用户身份认证、权限管理及 Token 生成，确保系统安全性。
2.        用户服务（User Service）：管理用户的基本信息、账号管理等操作。
3.        商品服务（Product Service）：提供商品的增删改查、库存管理及分类查询等功能。
4.        购物车服务（Cart Service）：处理用户的购物车操作，包括添加、删除、修改商品，以及购物车结算等功能。
5.        订单服务（Order Service）：管理订单的创建、支付、状态更新、超时取消等流程，确保订单数据一致性。
6.        支付服务（Payment Service）：对接第三方支付渠道（如支付宝、微信支付），处理支付请求及支付状态回调。
7.        结算服务（Settlement Service）：处理订单结算等相关操作。
8.        网关服务（Gateway Service）：作为 API 入口，负责服务发现、请求路由、流量控制及安全防护，处理来自前端的 HTTP 请求。
架构图如下:

3.3 项目代码介绍
go
.
├── 1.txt
├── Makefile
├── README.md
├── app
│   ├── auth //剩余微服务类似
│   │   ├── biz
│   │   │   ├── dal
│   │   │   ├── model
│   │   │   ├── service
│   │   │   └── utils
│   │   ├── build.sh
│   │   ├── conf
│   │   ├── go.mod
│   │   ├── go.sum
│   │   ├── handler.go
│   │   ├── kitex_info.yaml
│   │   ├── log
│   │   ├── main.go
│   ├── cart
│   ├── checkout
│   ├── hertz_gateway //网关
│   │   ├── biz
│   │   │   ├── dal
│   │   │   │   ├── init.go
│   │   │   │   ├── mysql
│   │   │   │   │   └── init.go
│   │   │   │   └── redis
│   │   │   │       └── init.go
│   │   │   ├── handler
│   │   │   │   ├── auth
│   │   │   │   │   ├── auth_service.go
│   │   │   │   │   └── auth_service_test.go
│   │   │   │   ├── cart
│   │   │   │   │   ├── cart_service.go
│   │   │   │   │   └── cart_service_test.go
│   │   │   │   ├── category
│   │   │   │   │   ├── category_service.go
│   │   │   │   │   └── category_service_test.go
│   │   │   │   ├── checkout
│   │   │   │   │   ├── checkout_service.go
│   │   │   │   │   └── checkout_service_test.go
│   │   │   │   ├── home
│   │   │   │   │   ├── home_service.go
│   │   │   │   │   └── home_service_test.go
│   │   │   │   ├── order
│   │   │   │   │   ├── order_service.go
│   │   │   │   │   └── order_service_test.go
│   │   │   │   ├── product
│   │   │   │   │   ├── product_service.go
│   │   │   │   │   └── product_service_test.go
│   │   │   │   └── user
│   │   │   │       ├── user_service.go
│   │   │   │       └── user_service_test.go
│   │   │   ├── router
│   │   │   │   ├── auth
│   │   │   │   │   ├── auth_api.go
│   │   │   │   │   └── middleware.go
│   │   │   │   ├── cart
│   │   │   │   │   ├── cart_page.go
│   │   │   │   │   └── middleware.go
│   │   │   │   ├── category
│   │   │   │   │   ├── category_page.go
│   │   │   │   │   └── middleware.go
│   │   │   │   ├── checkout
│   │   │   │   │   ├── checkout_page.go
│   │   │   │   │   └── middleware.go
│   │   │   │   ├── home
│   │   │   │   │   ├── home.go
│   │   │   │   │   └── middleware.go
│   │   │   │   ├── order
│   │   │   │   │   ├── middleware.go
│   │   │   │   │   └── order_page.go
│   │   │   │   ├── product
│   │   │   │   │   ├── middleware.go
│   │   │   │   │   └── product_page.go
│   │   │   │   ├── register.go
│   │   │   │   └── user
│   │   │   │       ├── middleware.go
│   │   │   │       └── user_api.go
│   │   │   ├── service
│   │   │   └── utils
│   │   ├── build.sh
│   │   ├── conf
│   │   ├── docker-compose.yaml
│   │   ├── go.mod
│   │   ├── go.sum
│   │   ├── hertz_gen
│   │   ├── infra
│   │   │   └── rpc
│   │   │       ├── client.go
│   │   │       └── client_test.go
│   │   ├── log
│   │   ├── main.go
│   │   ├── middleware
│   │   ├── readme.md
│   │   ├── script
│   │   ├── types
│   │   └── utils
│   ├── order
│   ├── payment
│   ├── product
│   └── user
├── common //公共组件
│   ├── constant
│   │   ├── error.go
│   │   └── role.go
│   └── go.mod
├── db
├── docker-compose.yaml
├── go.work
├── go.work.sum
├── idl //protobuf文件
└── rpc_gen
    ├── go.mod
    ├── go.sum
    ├── kitex_gen
    │   ├── auth //剩余微服务类似
    │   │   ├── auth.pb.fast.go
    │   │   ├── auth.pb.go
    │   │   └── authservice
    │   │       ├── authservice.go
    │   │       ├── client.go
    │   │       ├── invoker.go
    │   │       └── server.go
    │   ├── cart
    │   ├── checkout
    │   ├── order
    │   ├── payment
    │   ├── product
    │   └── user
    └── rpc
        ├── auth //剩余微服务类似
        │   ├── auth_client.go
        │   ├── auth_default.go
        │   └── auth_init.go
        ├── cart
        ├── checkout
        ├── order
        ├── payment
        ├── product
        └── user
利用cwgo生成rpc的客户端和服务端代码：
bash
.PHONY: gen-auth
gen-auth:
    @cd rpc_gen && cwgo client --type RPC --service auth --module ${ROOT_MOD}/rpc_gen -I ../idl --idl ../idl/auth.proto
    @cd app/auth && cwgo server --type RPC --service auth --module ${ROOT_MOD}/app/auth --pass "-use ${ROOT_MOD}/rpc_gen/kitex_gen" -I ../../idl --idl ../../idl/auth.proto
使用jwt进行认证服务：
go
var (
    arj  *ARJWT
    once sync.Once
)

type ARJWT struct {
    // 密钥，用以加密 JWT
    Key []byte

    // 定义 access token 过期时间（单位：分钟）即当颁发 access token 后，多少分钟后 access token 过期
    AccessExpireTime int64

    // 定义 refresh token 过期时间（单位：分钟）即当颁发 refresh token 后，多少分钟后 refresh token 过期
    RefreshExpireTime int64

    // token 的签发者
    Issuer string
}

type JWTCustomClaims struct {
    UserID   int64  `json:"user_id"`
    Username string `json:"username"`
    jwt.RegisteredClaims
}

func NewARJWT() *ARJWT {
    once.Do(func() {
        arj = &ARJWT{
            Key:               []byte(conf.GetConf().JWT.Secret),
            AccessExpireTime:  conf.GetConf().JWT.AccessExpireTime,
            RefreshExpireTime: conf.GetConf().JWT.RefreshExpireTime,
            Issuer:            conf.GetConf().JWT.Issuer,
        }
    })
    return arj
}

func (arj *ARJWT) GenerateToken(userId int64, username string) (accessToken, refreshToken string, err error) {
    // 生成 access token
    mc := JWTCustomClaims{
        UserID:   userId,
        Username: username,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(arj.AccessExpireTime) * time.Minute)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Issuer:    arj.Issuer,
        },
    }

    accessToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, mc).SignedString(arj.Key)
    if err != nil {
        log.Printf("generate access token failed: %v \n", err)
        return "", "", err
    }

    // 生成 refresh token
    refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
        ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(arj.RefreshExpireTime) * time.Minute)),
        Issuer:    arj.Issuer,
    }).SignedString(arj.Key)
    if err != nil {
        log.Printf("generate refresh token failed: %v \n", err)
        return "", "", err
    }
    return
}

func (arj *ARJWT) ParseAccessToken(tokenString string) (*JWTCustomClaims, error) {
    claims := JWTCustomClaims{}

    token, err := jwt.ParseWithClaims(tokenString, &claims,
        func(token *jwt.Token) (interface{}, error) {
            return arj.Key, nil
        },
    )

    if err != nil {
        validationErr, ok := err.(*jwt.ValidationError)
        if ok {
            switch validationErr.Errors {
            case jwt.ValidationErrorMalformed:
                return nil, errors.New("请求令牌格式有误")
            case jwt.ValidationErrorExpired:
                return nil, errors.New("令牌已过期")
            }
        }
        return nil, errors.New("请求令牌无效")
    }

    if _, ok := token.Claims.(*JWTCustomClaims); ok && token.Valid {
        return &claims, nil
    }

    return nil, errors.New("请求令牌无效")
}

func (arj *ARJWT) RefreshToken(accessToken, refreshToken string) (newAccessToken, newRefreshToken string, err error) {
    // 先判断 refresh token 是否有效
    if _, err = jwt.Parse(refreshToken,
        func(token *jwt.Token) (interface{}, error) {
            return arj.Key, nil
        },
    ); err != nil {
        return
    }

    // 从旧的 access token 中解析出 JWTCustomClaims 数据出来
    claims := JWTCustomClaims{}
    _, err = jwt.ParseWithClaims(accessToken, &claims,
        func(token *jwt.Token) (interface{}, error) {
            return arj.Key, nil
        },
    )
    if err != nil {
        validationErr, ok := err.(*jwt.ValidationError)
        // 当 access token 是过期错误，并且 refresh token 没有过期时就创建一个新的 access token 和 refresh token
        if ok && validationErr.Errors == jwt.ValidationErrorExpired {
            // 重新生成新的 access token 和 refresh token
            return arj.GenerateToken(claims.UserID, claims.Username)
        }
    }

    return accessToken, refreshToken, errors.New("access token still valid")
}
四、测试结果
建议从功能测试和性能测试两部分分析，其中功能测试补充测试用例，性能测试补充性能分析报告、可优化点等内容。
*功能测试为必填
功能测试：
用户服务：
创建用户：

获取用户信息：

更新用户：

删除用户：

登录：

更新token：

角色服务:
获取用户角色：

添加用户角色:

删除用户角色：

商品服务：
创建商品：

查找商品：

分类商品：

查找某类商品：

创建分类：

购物车服务：
获取购物车：

添加产品到购物车：

结算服务：
结算服务：

支付服务：

订单服务：
查看订单：

