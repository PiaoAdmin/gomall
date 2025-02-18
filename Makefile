export ROOT_MOD=github.com/PiaoAdmin/gomall

.PHONY: gen-order
gen-order:
	@cd rpc_gen && cwgo client --type RPC --service order --module $(ROOT_MOD)/rpc_gen -I ../idl  --idl ../idl/order.proto
	@cd app/order && cwgo server --type RPC --service order --module $(ROOT_MOD)/app/order --pass "-use $(ROOT_MOD)/rpc_gen/kitex_gen" -I ../../idl --idl ../../idl/order.proto