package client

import (
	"fmt"
	"time"

	"github.com/curtisnewbie/gocommon/common"
	"github.com/curtisnewbie/miso/miso"
)

type OperateLog struct {
	OperateName  string     `json:"operateName"`
	OperateDesc  string     `json:"operateDesc"`
	OperateTime  miso.ETime `json:"operateTime"`
	OperateParam string     `json:"operateParam"`
	Username     string     `json:"username"`
	UserId       int        `json:"userId"`
}

func DispatchOperateLog(ec miso.Rail, ol OperateLog) error {
	return miso.PublishJson(ec, ol, "auth.operate-log.exg", "auth.operate-log.save")
}

func DispatchUserOpLog(rail miso.Rail, opName string, opDesc string, param any, user common.User) {
	if err := DispatchOperateLog(rail, OperateLog{
		OperateName:  opName,
		OperateDesc:  opDesc,
		OperateTime:  miso.ETime(time.Now()),
		OperateParam: fmt.Sprintf("%+v", param),
		Username:     user.Username,
		UserId:       user.UserId,
	}); err != nil {
		rail.Errorf("Failed to dispatch operate log, %v", err)
	}
}
