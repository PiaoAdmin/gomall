/*
 * @Author: liaosijie
 * @Date: 2025-02-16 17:45:54
 * @Last Modified by: liaosijie
 * @Last Modified time: 2025-02-16 17:53:53
 */

package utils

import (
	"log"

	"github.com/bwmarrin/snowflake"
)

func CreateId(nodeId int64) int64 {
	node, err := snowflake.NewNode(nodeId)
	if err != nil {
		log.Fatalf("failed to create snowflake node: %v", err)
	}
	return node.Generate().Int64()
}
