package uniqueid

import (
	"github.com/bwmarrin/snowflake"
)

var node *snowflake.Node

const (
	// NodeID 机器码，在分布式系统中应该通过配置或服务发现获取
	NodeID = 1
)

func init() {
	var err error
	node, err = snowflake.NewNode(NodeID)
	if err != nil {
		panic(err)
	}
}

func GenId() uint64 {
	return uint64(node.Generate().Int64())
}
