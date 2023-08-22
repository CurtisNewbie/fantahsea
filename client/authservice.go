package client

import (
	"fmt"
	"time"

	"github.com/curtisnewbie/gocommon/common"
	"github.com/curtisnewbie/gocommon/rabbitmq"
)

type OperateLog struct {
	OperateName  string       `json:"operateName"`
	OperateDesc  string       `json:"operateDesc"`
	OperateTime  common.ETime `json:"operateTime"`
	OperateParam string       `json:"operateParam"`
	Username     string       `json:"username"`
	UserId       int          `json:"userId"`
}

func DispatchOperateLog(ec common.Rail, ol OperateLog) error {
	return rabbitmq.PublishJson(ec, ol, "auth.operate-log.exg", "auth.operate-log.save")
}

func DispatchUserOpLog(rail common.Rail, opName string, opDesc string, param any, user common.User) {
	if err := DispatchOperateLog(rail, OperateLog{
		OperateName:  opName,
		OperateDesc:  opDesc,
		OperateTime:  common.ETime(time.Now()),
		OperateParam: fmt.Sprintf("%+v", param),
		Username:     user.Username,
		UserId:       user.UserId,
	}); err != nil {
		rail.Errorf("Failed to dispatch operate log, %v", err)
	}
}
