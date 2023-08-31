package client

import (
	"fmt"
	"time"

	"github.com/curtisnewbie/gocommon/common"
	"github.com/curtisnewbie/miso/core"
	"github.com/curtisnewbie/miso/rabbitmq"
)

type OperateLog struct {
	OperateName  string       `json:"operateName"`
	OperateDesc  string       `json:"operateDesc"`
	OperateTime  core.ETime `json:"operateTime"`
	OperateParam string       `json:"operateParam"`
	Username     string       `json:"username"`
	UserId       int          `json:"userId"`
}

func DispatchOperateLog(ec core.Rail, ol OperateLog) error {
	return rabbitmq.PublishJson(ec, ol, "auth.operate-log.exg", "auth.operate-log.save")
}

func DispatchUserOpLog(rail core.Rail, opName string, opDesc string, param any, user common.User) {
	if err := DispatchOperateLog(rail, OperateLog{
		OperateName:  opName,
		OperateDesc:  opDesc,
		OperateTime:  core.ETime(time.Now()),
		OperateParam: fmt.Sprintf("%+v", param),
		Username:     user.Username,
		UserId:       user.UserId,
	}); err != nil {
		rail.Errorf("Failed to dispatch operate log, %v", err)
	}
}
