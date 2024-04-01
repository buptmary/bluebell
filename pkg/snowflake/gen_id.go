package snowflake

import (
	"bluebell/settings"
	"github.com/bwmarrin/snowflake"
	"time"
)

var node *snowflake.Node

// Init 传入机器ID
func Init(machineID int64) (err error) {
	var st time.Time
	st, err = time.Parse("2006-01-02", settings.Conf.StartTime)
	if err != nil {
		return
	}

	snowflake.Epoch = st.UnixNano() / 1000000
	node, err = snowflake.NewNode(machineID)
	return err
}

// GenID 返回生成的id
func GenID() int64 {
	return node.Generate().Int64()
}
