package client

import (
	"context"
	"testing"
	"time"

	"github.com/curtisnewbie/gocommon/common"
	"github.com/curtisnewbie/gocommon/rabbitmq"
)

func TestDispatchOperateLog(t *testing.T) {
	common.LoadConfigFromFile("../app-conf-dev.yml")			
	rabbitmq.StartRabbitMqClient(context.Background())

	ol := OperateLog{
		OperateName:  "Fantahsea test operate log",
		OperateDesc:  "just a unit test",
		OperateTime:  common.ETime(time.Now()),
		OperateParam: "{  }",
		Username:     "yongj.zhuang",
		UserId:       1,
	}
	err := DispatchOperateLog(common.EmptyExecContext(), ol)
	if err != nil {
		t.Fatal(err)
	}

}
