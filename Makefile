export ROOT_MOD=github.com/PiaoAdmin/gomall
.PHONY:gen-product
gen-product:
	@cd rpc_gen && cwgo client --type RPC --service product --module ${ROOT_MOD}/rpc_gen -I ../idl --idl ../idl/product.proto
	@cd app/product && cwgo server --type RPC --service product --module ${ROOT_MOD}/app/product --pass "-use ${ROOT_MOD}/rpc_gen/kitex_gen" -I ../../idl --idl ../../idl/product.proto
.PHONY:gen-product-page
gen-product-page:
	@cd app/hertz_gateway && cwgo server -I ../../idl --type HTTP --service hertz_gateway --module $(ROOT_MOD)/app/hertz_gateway --idl ../../idl/hertz_gateway/product_page.proto
	@cd app/hertz_gateway && cwgo server -I ../../idl --type HTTP --service hertz_gateway --module $(ROOT_MOD)/app/hertz_gateway --idl ../../idl/hertz_gateway/category_page.proto


.PHONY: gen-order
gen-order:
	@cd rpc_gen && cwgo client --type RPC --service order --module $(ROOT_MOD)/rpc_gen -I ../idl  --idl ../idl/order.proto
	@cd app/order && cwgo server --type RPC --service order --module $(ROOT_MOD)/app/order --pass "-use $(ROOT_MOD)/rpc_gen/kitex_gen" -I ../../idl --idl ../../idl/order.proto

.PHONY: gen-hertz_gateway
gen-hertz_gateway:
	@cd app/hertz_gateway && cwgo server -I ../../idl --type HTTP --service hertz_gateway --module $(ROOT_MOD)/gomall/app/hertz_gateway --idl ../../idl/hertz_gateway/order_page.proto