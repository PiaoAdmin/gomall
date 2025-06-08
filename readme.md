**你需求写的不队抖音商城项目文档**

**一、项目介绍**

| <font style="color:#646a73;">概况：基于</font><font style="color:#646a73;">go</font><font style="color:#646a73;">语言实现的抖音商城微服务项目，具有完善的用户身份认证、商品管理、购物车、订单管理和支付等核心功能。</font><br/><font style="color:#646a73;">项目地址：</font><font style="color:#646a73;">https://github.com/PiaoAdmin/gomall/tree/dev</font> |
| --- |


**二、项目分工**

| **团队成员** | **主要贡献** |
| --- | --- |
| 张浩辰 | 负责开发**商品模块**，处理git代码合并冲突 |
| 汤兵兵（队长） | 负责开发**用户模块**和**鉴权模块**，**整合网关和各个微服务之间的接口调用以及后续优化** |
| 杨国栋 | 负责开发**购物车服务** |
| 卢嘉钦 | 负责开发**支付服务和结算** |
| 廖思杰 | 负责开发**订单服务** |


**三、项目实现**

**3.1 ****技术选型与相关开发文档**

| <font style="color:#646a73;">可以补充场景分析环节，明确要解决的问题和前提假设，比如按当前选型和架构总体预计需要</font><font style="color:#646a73;">xxx</font><font style="color:#646a73;">存储空间，</font><font style="color:#646a73;">xxx</font><font style="color:#646a73;">台服务器</font><font style="color:#646a73;">......</font><font style="color:#646a73;">。</font> |
| --- |


**场景分析**

本项目为**抖音商城**，需要支持大规模用户访问和高并发交易处理，核心挑战包括：

<font style="color:#3370ff;">•</font><font style="color:#3370ff;">          </font>**高并发与高可用**：秒杀、促销等场景下，流量会瞬时激增，系统需具备良好的负载均衡和弹性扩展能力。

<font style="color:#3370ff;">•</font><font style="color:#3370ff;">          </font>**数据一致性**：订单、库存、支付等关键业务需要保证事务一致性，避免超卖等问题。

<font style="color:#3370ff;">•</font><font style="color:#3370ff;">          </font>**低延迟要求**：用户期望快速响应，API 接口需具备毫秒级响应能力，减少用户等待时间。

<font style="color:#3370ff;">•</font><font style="color:#3370ff;">          </font>**存储与计算需求**：因为用户和订单数量大需要分布式存储，所以要保证存储时不出现主键冲突。

**技术选型**

基于上述需求，技术栈的选择遵循**高并发、低延迟、可扩展性**原则，选定如下方案：

<font style="color:#3370ff;">•</font><font style="color:#3370ff;">          </font>**开发语言：****Go**

<font style="color:#3370ff;">•</font><font style="color:#3370ff;">          </font>**微服务架构：****Hertz + Kitex**

<font style="color:#3370ff;">￮</font><font style="color:#3370ff;">        </font>**Hertz****（****HTTP ****框架）**：

<font style="color:#3370ff;">▪</font><font style="color:#3370ff;">         </font>针对高并发 HTTP 请求进行了优化，使用其做为统一的网关

<font style="color:#3370ff;">￮</font><font style="color:#3370ff;">        </font>**Kitex****（****RPC ****框架）**：

<font style="color:#3370ff;">▪</font><font style="color:#3370ff;">         </font>支持高效的 RPC 通信，减少跨服务调用开销，使用其做微服务

<font style="color:#3370ff;">￮</font><font style="color:#3370ff;">        </font>**Consul:**

<font style="color:#3370ff;">▪</font><font style="color:#3370ff;">         </font>服务发现和服务注册

<font style="color:#3370ff;">•</font><font style="color:#3370ff;">          </font>**消息队列：****RocketMQ**（用于削峰填谷，支持异步任务，如订单超时取消）

<font style="color:#3370ff;">•</font><font style="color:#3370ff;">          </font>**数据库：****MySQL + Redis**

<font style="color:#3370ff;">￮</font><font style="color:#3370ff;">        </font>**MySQL**：存储核心业务数据，如用户、订单、支付信息

<font style="color:#3370ff;">￮</font><font style="color:#3370ff;">        </font>**Redis**：用于缓存、热点数据存储、订单流水号生成等，减少数据库压力

<font style="color:#3370ff;">•</font><font style="color:#3370ff;">          </font>**容器化与运维**

<font style="color:#3370ff;">￮</font><font style="color:#3370ff;">        </font>**Kubernetes****（****K8s****）**：管理微服务，提供弹性伸缩

<font style="color:#3370ff;">￮</font><font style="color:#3370ff;">        </font>**Prometheus + Grafana**：监控系统性能，优化服务稳定性

**3.2 ****架构设计**

总体架构采用**微服务架构**，主要包括以下几个核心模块：

<font style="color:#3370ff;">1.</font><font style="color:#3370ff;">        </font>**认证服务（****Auth Service****）**：负责用户身份认证、权限管理及 Token 生成，确保系统安全性。

<font style="color:#3370ff;">2.</font><font style="color:#3370ff;">        </font>**用户服务（****User Service****）**：管理用户的基本信息、账号管理等操作。

<font style="color:#3370ff;">3.</font><font style="color:#3370ff;">        </font>**商品服务（****Product Service****）**：提供商品的增删改查、库存管理及分类查询等功能。

<font style="color:#3370ff;">4.</font><font style="color:#3370ff;">        </font>**购物车服务（****Cart Service****）**：处理用户的购物车操作，包括添加、删除、修改商品，以及购物车结算等功能。

<font style="color:#3370ff;">5.</font><font style="color:#3370ff;">        </font>**订单服务（****Order Service****）**：管理订单的创建、支付、状态更新、超时取消等流程，确保订单数据一致性。

<font style="color:#3370ff;">6.</font><font style="color:#3370ff;">        </font>**支付服务（****Payment Service****）**：对接第三方支付渠道（如支付宝、微信支付），处理支付请求及支付状态回调。

<font style="color:#3370ff;">7.</font><font style="color:#3370ff;">        </font>**结算服务（****Settlement Service****）**：处理订单结算等相关操作。

<font style="color:#3370ff;">8.</font><font style="color:#3370ff;">        </font>**网关服务（****Gateway Service****）**：作为 API 入口，负责服务发现、请求路由、流量控制及安全防护，处理来自前端的 HTTP 请求。

架构图如下:

![](https://cdn.nlark.com/yuque/0/2025/png/43778648/1749373078084-318d6088-e735-41f0-b653-fd988ca0001c.png)

**3.3 ****项目代码介绍**

| <font style="color:#646a73;">go   </font><font style="color:black;">.   </font><font style="color:black;">├── 1.txt   </font><font style="color:black;">├── Makefile   </font><font style="color:black;">├── README.md   </font><font style="color:black;">├── app   </font><font style="color:black;">│   ├── auth //</font><font style="color:black;">剩余微服务类似</font><font style="color:black;">   </font><font style="color:black;">│   │   ├── biz   </font><font style="color:black;">│   │   │   ├── dal   </font><font style="color:black;">│   │   │   ├── model   </font><font style="color:black;">│   │   │   ├── service   </font><font style="color:black;">│   │   │   └── utils   </font><font style="color:black;">│   │   ├── build.sh   </font><font style="color:black;">│   │   ├── conf   </font><font style="color:black;">│   │   ├── go.mod   </font><font style="color:black;">│   │   ├── go.sum   </font><font style="color:black;">│   │   ├── handler.go   </font><font style="color:black;">│   │   ├── kitex_info.yaml   </font><font style="color:black;">│   │   ├── log   </font><font style="color:black;">│   │   ├── main.go   </font><font style="color:black;">│   ├── cart   </font><font style="color:black;">│   ├── checkout   </font><font style="color:black;">│   ├── hertz_gateway //</font><font style="color:black;">网关</font><font style="color:black;">   </font><font style="color:black;">│   │   ├── biz   </font><font style="color:black;">│   │   │   ├── dal   </font><font style="color:black;">│   │   │   │   ├── init.go   </font><font style="color:black;">│   │   │   │   ├── mysql   </font><font style="color:black;">│   │   │   │   │   └── init.go   </font><font style="color:black;">│   │   │   │   └── redis   </font><font style="color:black;">│   │   │   │  </font><font style="color:black;">     </font><font style="color:black;">└── init.go   </font><font style="color:black;">│   │   │   ├── handler   </font><font style="color:black;">│   │   │   │   ├── auth   </font><font style="color:black;">│   │   │   │   │   ├── auth_service.go   </font><font style="color:black;">│   │   │   │   │   └── auth_service_test.go   </font><font style="color:black;">│   │   │   │   ├── cart   </font><font style="color:black;">│   │   │   │   │   ├── cart_service.go   </font><font style="color:black;">│   │   │   │   │   └── cart_service_test.go   </font><font style="color:black;">│   │   │   │   ├── category   </font><font style="color:black;">│   │   │   │   │   ├── category_service.go   </font><font style="color:black;">│   │   │   │   │   └── category_service_test.go   </font><font style="color:black;">│   │   │   │   ├── checkout   </font><font style="color:black;">│   │   │   │   │   ├── checkout_service.go   </font><font style="color:black;">│   │   │   │   │   └── checkout_service_test.go   </font><font style="color:black;">│   │   │   │   ├── home   </font><font style="color:black;">│   │   │   │   │   ├── home_service.go   </font><font style="color:black;">│   │   │   │   │   └── home_service_test.go   </font><font style="color:black;">│   │   │   │   ├── order   </font><font style="color:black;">│   │   │   │   │   ├── order_service.go   </font><font style="color:black;">│   │   │   │   │   └── order_service_test.go   </font><font style="color:black;">│   │   │   │   ├── product   </font><font style="color:black;">│   │   │   │   │   ├── product_service.go   </font><font style="color:black;">│   │   │   │   │   └── product_service_test.go   </font><font style="color:black;">│   │   │   │   └── user   </font><font style="color:black;">│   │   │   │  </font><font style="color:black;">     </font><font style="color:black;">├── user_service.go   </font><font style="color:black;">│   │   │   │  </font><font style="color:black;">     </font><font style="color:black;">└── user_service_test.go   </font><font style="color:black;">│   │   │   ├── router   </font><font style="color:black;">│   │   │   │   ├── auth   </font><font style="color:black;">│   │   │   │   │   ├── auth_api.go   </font><font style="color:black;">│   │   │   │   │   └── middleware.go   </font><font style="color:black;">│   │   │   │   ├── cart   </font><font style="color:black;">│   │   │   │   │   ├── cart_page.go   </font><font style="color:black;">│   │   │   │   │   └── middleware.go   </font><font style="color:black;">│   │   │   │   ├── category   </font><font style="color:black;">│   │   │   │   │   ├── category_page.go   </font><font style="color:black;">│   │   │   │   │   └── middleware.go   </font><font style="color:black;">│   │   │   │   ├── checkout   </font><font style="color:black;">│   │   │   │   │   ├── checkout_page.go   </font><font style="color:black;">│   │   │   │   │   └── middleware.go   </font><font style="color:black;">│   │   │   │   ├── home   </font><font style="color:black;">│   │   │   │   │   ├── home.go   </font><font style="color:black;">│   │   │   │   │   └── middleware.go   </font><font style="color:black;">│   │   │   │   ├── order   </font><font style="color:black;">│   │   │   │   │   ├── middleware.go   </font><font style="color:black;">│   │   │   │   │   └── order_page.go   </font><font style="color:black;">│   │   │   │   ├── product   </font><font style="color:black;">│   │   │   │   │   ├── middleware.go   </font><font style="color:black;">│   │   │   │   │   └── product_page.go   </font><font style="color:black;">│   │   │   │   ├── register.go   </font><font style="color:black;">│   │   │   │   └── user   </font><font style="color:black;">│   │   │   │  </font><font style="color:black;">     </font><font style="color:black;">├── middleware.go   </font><font style="color:black;">│   │   │   │  </font><font style="color:black;">     </font><font style="color:black;">└── user_api.go   </font><font style="color:black;">│   │   │   ├── service   </font><font style="color:black;">│   │   │   └── utils   </font><font style="color:black;">│   │   ├── build.sh   </font><font style="color:black;">│   │   ├── conf   </font><font style="color:black;">│   │   ├── docker-compose.yaml   </font><font style="color:black;">│   │   ├── go.mod   </font><font style="color:black;">│   │   ├── go.sum   </font><font style="color:black;">│   │   ├── hertz_gen   </font><font style="color:black;">│   │   ├── infra   </font><font style="color:black;">│   │   │   └── rpc   </font><font style="color:black;">│   │   │  </font><font style="color:black;">     </font><font style="color:black;">├── client.go   </font><font style="color:black;">│   │   │  </font><font style="color:black;">     </font><font style="color:black;">└── client_test.go   </font><font style="color:black;">│   │   ├── log   </font><font style="color:black;">│   │   ├── main.go   </font><font style="color:black;">│   │   ├── middleware   </font><font style="color:black;">│   │   ├── readme.md   </font><font style="color:black;">│   │   ├── script   </font><font style="color:black;">│   │   ├── types   </font><font style="color:black;">│   │   └── utils   </font><font style="color:black;">│   ├── order   </font><font style="color:black;">│   ├── payment   </font><font style="color:black;">│   ├── product   </font><font style="color:black;">│   └── user   </font><font style="color:black;">├── common //</font><font style="color:black;">公共组件</font><font style="color:black;">   </font><font style="color:black;">│   ├── constant   </font><font style="color:black;">│   │   ├── error.go   </font><font style="color:black;">│   │   └── role.go   </font><font style="color:black;">│   └── go.mod   </font><font style="color:black;">├── db   </font><font style="color:black;">├── docker-compose.yaml   </font><font style="color:black;">├── go.work   </font><font style="color:black;">├── go.work.sum   </font><font style="color:black;">├── idl //protobuf</font><font style="color:black;">文件</font><font style="color:black;">   </font><font style="color:black;">└── rpc_gen   </font><font style="color:black;">    </font><font style="color:black;">├── go.mod   </font><font style="color:black;">    </font><font style="color:black;">├── go.sum   </font><font style="color:black;">    </font><font style="color:black;">├── kitex_gen   </font><font style="color:black;">    </font><font style="color:black;">│   ├── auth //</font><font style="color:black;">剩余微服务类似</font><font style="color:black;">   </font><font style="color:black;">    </font><font style="color:black;">│   │   ├── auth.pb.fast.go   </font><font style="color:black;">    </font><font style="color:black;">│   │   ├── auth.pb.go   </font><font style="color:black;">    </font><font style="color:black;">│   │   └── authservice   </font><font style="color:black;">    </font><font style="color:black;">│   │  </font><font style="color:black;">     </font><font style="color:black;">├── authservice.go   </font><font style="color:black;">    </font><font style="color:black;">│   │  </font><font style="color:black;">     </font><font style="color:black;">├── client.go   </font><font style="color:black;">    </font><font style="color:black;">│   │  </font><font style="color:black;">     </font><font style="color:black;">├── invoker.go   </font><font style="color:black;">    </font><font style="color:black;">│   │  </font><font style="color:black;">     </font><font style="color:black;">└── server.go   </font><font style="color:black;">    </font><font style="color:black;">│   ├── cart   </font><font style="color:black;">    </font><font style="color:black;">│   ├── checkout   </font><font style="color:black;">    </font><font style="color:black;">│   ├── order   </font><font style="color:black;">    </font><font style="color:black;">│   ├── payment   </font><font style="color:black;">    </font><font style="color:black;">│   ├── product   </font><font style="color:black;">    </font><font style="color:black;">│   └── user   </font><font style="color:black;">    </font><font style="color:black;">└── rpc   </font><font style="color:black;">        </font><font style="color:black;">├── auth //</font><font style="color:black;">剩余微服务类似</font><font style="color:black;">   </font><font style="color:black;">        </font><font style="color:black;">│   ├── auth_client.go   </font><font style="color:black;">        </font><font style="color:black;">│   ├── auth_default.go   </font><font style="color:black;">        </font><font style="color:black;">│   └── auth_init.go   </font><font style="color:black;">        </font><font style="color:black;">├── cart   </font><font style="color:black;">        </font><font style="color:black;">├── checkout   </font><font style="color:black;">        </font><font style="color:black;">├── order   </font><font style="color:black;">        </font><font style="color:black;">├── payment   </font><font style="color:black;">        </font><font style="color:black;">├── product   </font><font style="color:black;">        </font><font style="color:black;">└── user</font> |
| --- |


利用cwgo生成rpc的客户端和服务端代码：

| <font style="color:#646a73;">bash   </font><font style="color:black;">.PHONY: gen-auth   </font><font style="color:black;">gen-auth:   </font><font style="color:black;">    </font><font style="color:black;">@cd rpc_gen && cwgo client --type RPC --service auth --module ${ROOT_MOD}/rpc_gen -I ../idl --idl ../idl/auth.proto   </font><font style="color:black;">    </font><font style="color:black;">@cd app/auth && cwgo server --type RPC --service auth --module ${ROOT_MOD}/app/auth --pass "-use ${ROOT_MOD}/rpc_gen/kitex_gen" -I ../../idl --idl ../../idl/auth.proto</font> |
| --- |


使用jwt进行认证服务：

| <font style="color:#646a73;">go   </font><font style="color:black;">var (   </font><font style="color:black;">    </font><font style="color:black;">arj</font><font style="color:black;">  </font><font style="color:black;">*ARJWT   </font><font style="color:black;">    </font><font style="color:black;">once sync.Once   </font><font style="color:black;">)      </font><font style="color:black;">type ARJWT struct {   </font><font style="color:black;">    </font><font style="color:black;">// </font><font style="color:black;">密钥，用以加密</font><font style="color:black;"> JWT   </font><font style="color:black;">    </font><font style="color:black;">Key []byte      </font><font style="color:black;">    </font><font style="color:black;">// </font><font style="color:black;">定义</font><font style="color:black;"> access token </font><font style="color:black;">过期时间（单位：分钟）即当颁发</font><font style="color:black;"> access token </font><font style="color:black;">后，多少分钟后</font><font style="color:black;"> access token </font><font style="color:black;">过期</font><font style="color:black;">   </font><font style="color:black;">    </font><font style="color:black;">AccessExpireTime int64      </font><font style="color:black;">    </font><font style="color:black;">// </font><font style="color:black;">定义</font><font style="color:black;"> refresh token </font><font style="color:black;">过期时间（单位：分钟）即当颁发</font><font style="color:black;"> refresh token </font><font style="color:black;">后，多少分钟后</font><font style="color:black;"> refresh token </font><font style="color:black;">过期</font><font style="color:black;">   </font><font style="color:black;">    </font><font style="color:black;">RefreshExpireTime int64      </font><font style="color:black;">    </font><font style="color:black;">// token </font><font style="color:black;">的签发者</font><font style="color:black;">   </font><font style="color:black;">    </font><font style="color:black;">Issuer string   </font><font style="color:black;">}      </font><font style="color:black;">type JWTCustomClaims struct {   </font><font style="color:black;">    </font><font style="color:black;">UserID</font><font style="color:black;">   </font><font style="color:black;">int64</font><font style="color:black;">  </font><font style="color:black;">`json:"user_id"`   </font><font style="color:black;">    </font><font style="color:black;">Username string `json:"username"`   </font><font style="color:black;">    </font><font style="color:black;">jwt.RegisteredClaims   </font><font style="color:black;">}      </font><font style="color:black;">func NewARJWT() *ARJWT {   </font><font style="color:black;">   </font><font style="color:black;"> </font><font style="color:black;">once.Do(func() {   </font><font style="color:black;">        </font><font style="color:black;">arj = &ARJWT{   </font><font style="color:black;">            </font><font style="color:black;">Key:</font><font style="color:black;">               </font><font style="color:black;">[]byte(conf.GetConf().JWT.Secret),   </font><font style="color:black;">            </font><font style="color:black;">AccessExpireTime:</font><font style="color:black;">  </font><font style="color:black;">conf.GetConf().JWT.AccessExpireTime,   </font><font style="color:black;">            </font><font style="color:black;">RefreshExpireTime: conf.GetConf().JWT.RefreshExpireTime,   </font><font style="color:black;">            </font><font style="color:black;">Issuer:</font><font style="color:black;">            </font><font style="color:black;">conf.GetConf().JWT.Issuer,   </font><font style="color:black;">        </font><font style="color:black;">}   </font><font style="color:black;">    </font><font style="color:black;">})   </font><font style="color:black;">    </font><font style="color:black;">return arj   </font><font style="color:black;">}      </font><font style="color:black;">func (arj *ARJWT) GenerateToken(userId int64, username string) (accessToken, refreshToken string, err error) {   </font><font style="color:black;">    </font><font style="color:black;">// </font><font style="color:black;">生成</font><font style="color:black;"> access token   </font><font style="color:black;">    </font><font style="color:black;">mc := JWTCustomClaims{   </font><font style="color:black;">        </font><font style="color:black;">UserID:</font><font style="color:black;">   </font><font style="color:black;">userId,   </font><font style="color:black;">        </font><font style="color:black;">Username: username,   </font><font style="color:black;">        </font><font style="color:black;">RegisteredClaims: jwt.RegisteredClaims{   </font><font style="color:black;">            </font><font style="color:black;">ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(arj.AccessExpireTime) * time.Minute)),   </font><font style="color:black;">            </font><font style="color:black;">IssuedAt:</font><font style="color:black;">  </font><font style="color:black;">jwt.NewNumericDate(time.Now()),   </font><font style="color:black;">            </font><font style="color:black;">Issuer:</font><font style="color:black;">    </font><font style="color:black;">arj.Issuer,   </font><font style="color:black;">        </font><font style="color:black;">},   </font><font style="color:black;">    </font><font style="color:black;">}      </font><font style="color:black;">    </font><font style="color:black;">accessToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, mc).SignedString(arj.Key)   </font><font style="color:black;">    </font><font style="color:black;">if err != nil {   </font><font style="color:black;">        </font><font style="color:black;">log.Printf("generate access token failed: %v \n", err)   </font><font style="color:black;">        </font><font style="color:black;">return "", "", err   </font><font style="color:black;">    </font><font style="color:black;">}      </font><font style="color:black;">    </font><font style="color:black;">// </font><font style="color:black;">生成</font><font style="color:black;"> refresh token   </font><font style="color:black;">    </font><font style="color:black;">refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{   </font><font style="color:black;">        </font><font style="color:black;">ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(arj.RefreshExpireTime) * time.Minute)),   </font><font style="color:black;">        </font><font style="color:black;">Issuer:</font><font style="color:black;">    </font><font style="color:black;">arj.Issuer,   </font><font style="color:black;">    </font><font style="color:black;">}).SignedString(arj.Key)   </font><font style="color:black;">    </font><font style="color:black;">if err != nil {   </font><font style="color:black;">        </font><font style="color:black;">log.Printf("generate refresh token failed: %v \n", err)   </font><font style="color:black;">        </font><font style="color:black;">return "", "", err   </font><font style="color:black;">    </font><font style="color:black;">}   </font><font style="color:black;">    </font><font style="color:black;">return   </font><font style="color:black;">}      </font><font style="color:black;">func (arj *ARJWT) ParseAccessToken(tokenString string) (*JWTCustomClaims, error) {   </font><font style="color:black;">    </font><font style="color:black;">claims := JWTCustomClaims{}      </font><font style="color:black;">    </font><font style="color:black;">token, err := jwt.ParseWithClaims(tokenString, &claims,   </font><font style="color:black;">        </font><font style="color:black;">func(token *jwt.Token) (interface{}, error) {   </font><font style="color:black;">            </font><font style="color:black;">return arj.Key, nil   </font><font style="color:black;">        </font><font style="color:black;">},   </font><font style="color:black;">    </font><font style="color:black;">)      </font><font style="color:black;">    </font><font style="color:black;">if err != nil {   </font><font style="color:black;">        </font><font style="color:black;">validationErr, ok := err.(*jwt.ValidationError)   </font><font style="color:black;">        </font><font style="color:black;">if ok {   </font><font style="color:black;">            </font><font style="color:black;">switch validationErr.Errors {   </font><font style="color:black;">            </font><font style="color:black;">case jwt.ValidationErrorMalformed:   </font><font style="color:black;">                </font><font style="color:black;">return nil, errors.New("</font><font style="color:black;">请求令牌格式有误</font><font style="color:black;">")   </font><font style="color:black;">            </font><font style="color:black;">case jwt.ValidationErrorExpired:   </font><font style="color:black;">                </font><font style="color:black;">return nil, errors.New("</font><font style="color:black;">令牌已过期</font><font style="color:black;">")   </font><font style="color:black;">        </font><font style="color:black;">    </font><font style="color:black;">}   </font><font style="color:black;">        </font><font style="color:black;">}   </font><font style="color:black;">        </font><font style="color:black;">return nil, errors.New("</font><font style="color:black;">请求令牌无效</font><font style="color:black;">")   </font><font style="color:black;">    </font><font style="color:black;">}      </font><font style="color:black;">    </font><font style="color:black;">if _, ok := token.Claims.(*JWTCustomClaims); ok && token.Valid {   </font><font style="color:black;">        </font><font style="color:black;">return &claims, nil   </font><font style="color:black;">    </font><font style="color:black;">}      </font><font style="color:black;">    </font><font style="color:black;">return nil, errors.New("</font><font style="color:black;">请求令牌无效</font><font style="color:black;">")   </font><font style="color:black;">}      </font><font style="color:black;">func (arj *ARJWT) RefreshToken(accessToken, refreshToken string) (newAccessToken, newRefreshToken string, err error) {   </font><font style="color:black;">    </font><font style="color:black;">// </font><font style="color:black;">先判断</font><font style="color:black;"> refresh token </font><font style="color:black;">是否有效</font><font style="color:black;">   </font><font style="color:black;">    </font><font style="color:black;">if _, err = jwt.Parse(refreshToken,   </font><font style="color:black;">        </font><font style="color:black;">func(token *jwt.Token) (interface{}, error) {   </font><font style="color:black;">            </font><font style="color:black;">return arj.Key, nil   </font><font style="color:black;">        </font><font style="color:black;">},   </font><font style="color:black;">    </font><font style="color:black;">); err != nil {   </font><font style="color:black;">        </font><font style="color:black;">return   </font><font style="color:black;">    </font><font style="color:black;">}      </font><font style="color:black;">    </font><font style="color:black;">// </font><font style="color:black;">从旧的</font><font style="color:black;"> access token </font><font style="color:black;">中解析出</font><font style="color:black;"> JWTCustomClaims </font><font style="color:black;">数据出来</font><font style="color:black;">   </font><font style="color:black;">    </font><font style="color:black;">claims := JWTCustomClaims{}   </font><font style="color:black;">    </font><font style="color:black;">_, err = jwt.ParseWithClaims(accessToken, &claims,   </font><font style="color:black;">        </font><font style="color:black;">func(token *jwt.Token) (interface{}, error) {   </font><font style="color:black;">            </font><font style="color:black;">return arj.Key, nil   </font><font style="color:black;">        </font><font style="color:black;">},   </font><font style="color:black;">    </font><font style="color:black;">)   </font><font style="color:black;">    </font><font style="color:black;">if err != nil {   </font><font style="color:black;">        </font><font style="color:black;">validationErr, ok := err.(*jwt.ValidationError)   </font><font style="color:black;">        </font><font style="color:black;">// </font><font style="color:black;">当</font><font style="color:black;"> access token </font><font style="color:black;">是过期错误，并且</font><font style="color:black;"> refresh token </font><font style="color:black;">没有过期时就创建一个新的</font><font style="color:black;"> access token </font><font style="color:black;">和</font><font style="color:black;"> refresh token   </font><font style="color:black;">        </font><font style="color:black;">if ok && validationErr.Errors == jwt.ValidationErrorExpired {   </font><font style="color:black;">           </font><font style="color:black;"> </font><font style="color:black;">// </font><font style="color:black;">重新生成新的</font><font style="color:black;"> access token </font><font style="color:black;">和</font><font style="color:black;"> refresh token   </font><font style="color:black;">            </font><font style="color:black;">return arj.GenerateToken(claims.UserID, claims.Username)   </font><font style="color:black;">        </font><font style="color:black;">}   </font><font style="color:black;">    </font><font style="color:black;">}      </font><font style="color:black;">    </font><font style="color:black;">return accessToken, refreshToken, errors.New("access token still valid")   </font><font style="color:black;">}</font> |
| --- |


**四、测试结果**

| <font style="color:#646a73;">建议从功能测试和性能测试两部分分析，其中功能测试补充测试用例，性能测试补充性能分析报告、可优化点等内容。</font> |
| --- |


*******功能测试为必填**

**功能测试：**

**用户服务：**

创建用户：

![](https://cdn.nlark.com/yuque/0/2025/png/43778648/1749373078155-5085f32a-933f-448d-83df-9edfbb98b8c8.png)

获取用户信息：

![](https://cdn.nlark.com/yuque/0/2025/png/43778648/1749373078170-1ae768a6-e561-4c42-b0c4-6c71ea61b1fa.png)

更新用户：

![](https://cdn.nlark.com/yuque/0/2025/png/43778648/1749373078174-f517ebf2-1fdb-4ecb-9073-ead81daeaabe.png)

删除用户：

![](https://cdn.nlark.com/yuque/0/2025/png/43778648/1749373078308-2093d213-b03f-4c4b-bf5b-ba61ce70b1c3.png)

登录：

![](https://cdn.nlark.com/yuque/0/2025/png/43778648/1749373078476-2a9a50bf-a66f-48c2-a9a8-79b3e23059a4.png)

更新token：

![](https://cdn.nlark.com/yuque/0/2025/png/43778648/1749373078680-f6d65503-ef9d-45a2-9b48-a36065f19ada.png)

**角色服务****:**

获取用户角色：

![](https://cdn.nlark.com/yuque/0/2025/png/43778648/1749373078751-4fe8a3c0-d033-4ec0-87c7-814fd2594f67.png)

添加用户角色:

![](https://cdn.nlark.com/yuque/0/2025/png/43778648/1749373078622-d4f30b1c-10da-4f75-be99-c8ed94a66ec6.png)

删除用户角色：

![](https://cdn.nlark.com/yuque/0/2025/png/43778648/1749373078883-b919400d-70ba-4757-a197-7120e1819817.png)

**商品服务：**

创建商品：

![](https://cdn.nlark.com/yuque/0/2025/png/43778648/1749373079054-fe31719b-e621-437e-96e2-3728125ee6bc.png)

查找商品：

![](https://cdn.nlark.com/yuque/0/2025/png/43778648/1749373079058-8a220789-742a-4fae-a1ad-95a0a1e6efe9.png)

分类商品：

![](https://cdn.nlark.com/yuque/0/2025/png/43778648/1749373079401-d3ffe1b8-9b92-4643-9ac2-10929fe783db.png)

查找某类商品：

![](https://cdn.nlark.com/yuque/0/2025/png/43778648/1749373079690-adc24d4a-fd77-4bc3-b479-4157d15868a9.png)

创建分类：

![](https://cdn.nlark.com/yuque/0/2025/png/43778648/1749373079631-ceb48179-d3f8-4b50-9257-82065baff42b.png)

**购物车服务：**

获取购物车：

![](https://cdn.nlark.com/yuque/0/2025/png/43778648/1749373080290-52710807-b4a5-486c-ae3e-612adc0d513d.png)

添加产品到购物车：

![](https://cdn.nlark.com/yuque/0/2025/png/43778648/1749373080240-8660c8ab-15e5-4d0f-a1d6-9dec7bd182ba.png)

**结算服务：**

结算服务：

![](https://cdn.nlark.com/yuque/0/2025/png/43778648/1749373080341-d2e795fc-73ed-418e-b47a-841f34107f19.png)

支付服务：

![](https://cdn.nlark.com/yuque/0/2025/png/43778648/1749373080238-2ad2b602-5565-4ef0-a991-c7102700a1d2.png)

**订单服务：**

查看订单：

![](https://cdn.nlark.com/yuque/0/2025/png/43778648/1749373080638-f8f00a64-1b33-4575-b4b1-3d889963868d.png)



