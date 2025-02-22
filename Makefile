export ROOT_MOD=github.com/PiaoAdmin/gomall
.PHONY:gen-product
gen-product:
	@cd rpc_gen && cwgo client --type RPC --service product --module ${ROOT_MOD}/rpc_gen -I ../idl --idl ../idl/product.proto
	@cd app/product && cwgo server --type RPC --service product --module ${ROOT_MOD}/app/product --pass "-use ${ROOT_MOD}/rpc_gen/kitex_gen" -I ../../idl --idl ../../idl/product.proto
.PHONY:gen-product-page
gen-product-page:
	@cd app/hertz_gateway && cwgo server -I ../../idl --type HTTP --service hertz_gateway --module $(ROOT_MOD)/app/hertz_gateway --idl ../../idl/hertz_gateway/product_page.proto
	@cd app/hertz_gateway && cwgo server -I ../../idl --type HTTP --service hertz_gateway --module $(ROOT_MOD)/app/hertz_gateway --idl ../../idl/hertz_gateway/category_page.proto
.PHONY: gen-payment
gen-payment:
	@cd rpc_gen && cwgo client --I ../idl --type RPC --service payment --module github.com/PiaoAdmin/gomall/rpc_gen --idl ../idl/payment.proto
	@cd app/payment && cwgo server -I ../../idl --type RPC --service payment --module github.com/PiaoAdmin/gomall/app/payment --idl ../../idl/payment.proto --pass "-use github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen"

.PHONY: gen-checkout
gen-checkout:
	@cd rpc_gen && cwgo client --I ../idl --type RPC --service checkout --module github.com/PiaoAdmin/gomall/rpc_gen --idl ../idl/checkout.proto
	@cd app/checkout && cwgo server -I ../../idl --type RPC --service checkout --module github.com/PiaoAdmin/gomall/app/checkout --idl ../../idl/checkout.proto --pass "-use github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen"

.PHONY: gen-hertz_gateway-checkout
gen-hertz_gateway-checkout:
	@cd app/hertz_gateway && cwgo server -I ../../idl --type HTTP --service hertz_gateway --module github.com/PiaoAdmin/gomall/app/hertz_gateway --idl ../../idl/hertz_gateway/checkout_page.proto
