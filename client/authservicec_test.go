package client

import (
	"context"
	"testing"
	"time"

	"github.com/curtisnewbie/miso/core"
	"github.com/curtisnewbie/miso/rabbitmq"
)

func TestDispatchOperateLog(t *testing.T) {
	c := core.EmptyRail()
	core.LoadConfigFromFile("../app-conf-dev.yml", c)
	rabbitmq.StartRabbitMqClient(context.Background())

	ol := OperateLog{
		OperateName:  "Fantahsea test operate log",
		OperateDesc:  "just a unit test",
		OperateTime:  core.ETime(time.Now()),
		OperateParam: "{  }",
		Username:     "yongj.zhuang",
		UserId:       1,
	}
	err := DispatchOperateLog(c, ol)
	if err != nil {
		t.Fatal(err)
	}

}
