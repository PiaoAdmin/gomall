.PHONY: all user product cart checkout order payment api

MODULE_PREFIX = github.com/PiaoAdmin/pmall

# 定义生成服务的函数
# 参数 1: 服务名称 (对应 app/ 下的目录名)
# 参数 2: proto 文件名 (对应 idl/ 下的文件)
define gen_service
	@echo "Generating $(1) service..."
	mkdir -p ../app/$(1)
	cd ../app/$(1) && \
	if [ ! -f go.mod ]; then go mod init $(MODULE_PREFIX)/app/$(1); fi && \
	go mod edit -replace $(MODULE_PREFIX)/rpc_gen=../../rpc_gen && \
	kitex -module $(MODULE_PREFIX)/app/$(1) -gen-path ../../rpc_gen -service $(1) -I ../../idl ../../idl/$(2) && \
	go mod tidy
endef

# 定义生成 API 服务的函数
# 参数 1: proto 文件名 (对应 idl/api/ 下的文件)
define gen_api
	@echo "Updating api service with $(1)..."
	cd ../app/api && \
	hz update -module $(MODULE_PREFIX)/app/api -I ../../idl -idl ../../idl/api/$(1) && \
	go mod tidy
endef

all: user product cart checkout order payment api

user:
	$(call gen_service,user,user.proto)

product:
	$(call gen_service,product,product.proto)

cart:
	$(call gen_service,cart,cart.proto)

checkout:
	$(call gen_service,checkout,checkout.proto)

order:
	$(call gen_service,order,order.proto)

payment:
	$(call gen_service,payment,payment.proto)

api_home:
	$(call gen_api,home.proto)
api_auth:
	$(call gen_api,auth.proto) 

api_product:
	$(call gen_api,product_api.proto) 

api_cart:
	$(call gen_api,cart_api.proto) 

api_order:
	$(call gen_api,order_api.proto) 

