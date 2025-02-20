.PHONY: gen-payment
gen-payment:
	@cd rpc_gen && cwgo client --I ../idl --type RPC --service payment --module github.com/PiaoAdmin/gomall/rpc_gen --idl ../idl/payment.proto
	@cd app/payment && cwgo server -I ../../idl --type RPC --service payment --module github.com/PiaoAdmin/gomall/app/payment --idl ../../idl/payment.proto --pass "-use github.com/PiaoAdmin/gomall/rpc_gen/kitex_gen"