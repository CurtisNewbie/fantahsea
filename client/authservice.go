package client

import (
	"fmt"
	"strconv"
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

func DispatchOperateLog(ec common.ExecContext, ol OperateLog) error {
	return rabbitmq.PublishJson(ol, "auth.operate-log.exg", "auth.operate-log.save")
}

func DispatchUserOpLog(ec common.ExecContext, opName string, opDesc string, param any) {
	id, _ := strconv.Atoi(ec.User.UserId)
	if err := DispatchOperateLog(ec, OperateLog{
		OperateName:  opName,
		OperateDesc:  opDesc,
		OperateTime:  common.ETime(time.Now()),
		OperateParam: fmt.Sprintf("%+v", param),
		Username:     ec.User.Username,
		UserId:       id,
	}); err != nil {
		ec.Log.Errorf("Failed to dispatch operate log, %v", err)
	}
}
