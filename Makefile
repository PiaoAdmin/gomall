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
.PHONY: gen-cart
gen-cart:
	@cd rpc_gen && cwgo client --I ../idl --type RPC --service cart --module github.com/PiaoAdmin/gomall/rpc_gen --idl ../idl/cart.proto
	@cd app/cart && cwgo server -I ../../idl --type RPC --service cart --module github.com/PiaoAdmin/gomall/app/cart --idl ../../idl/cart.proto --pass "-use github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen"
.PHONY:gen-cart-page
gen-cart-page:
	@cd app/hertz_gateway && cwgo server -I ../../idl --type HTTP --service hertz_gateway --module $(ROOT_MOD)/app/hertz_gateway --idl ../../idl/hertz_gateway/cart_page.proto
.PHONY:gen-common-page
gen-common-page:
	@cd app/hertz_gateway && cwgo server -I ../../idl --type HTTP --service hertz_gateway --module $(ROOT_MOD)/app/hertz_gateway --idl ../../idl/hertz_gateway/common.proto

.PHONY: gen-order
gen-order:
	@cd rpc_gen && cwgo client --type RPC --service order --module ${ROOT_MOD}/rpc_gen -I ../idl --idl ../idl/order.proto
	@cd app/order && cwgo server --type RPC --service order --module ${ROOT_MOD}/app/order --pass "-use ${ROOT_MOD}/rpc_gen/kitex_gen" -I ../../idl --idl ../../idl/order.proto


.PHONY: gen-hertz_gateway-checkout
gen-hertz_gateway-checkout:
	@cd app/hertz_gateway && cwgo server -I ../../idl --type HTTP --service hertz_gateway --module github.com/PiaoAdmin/gomall/app/hertz_gateway --idl ../../idl/hertz_gateway/checkout_page.proto
