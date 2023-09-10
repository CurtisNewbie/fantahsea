package client

import (
	"testing"
	"time"

	"github.com/curtisnewbie/miso/miso"
)

func TestDispatchOperateLog(t *testing.T) {
	c := miso.EmptyRail()
	miso.LoadConfigFromFile("../app-conf-dev.yml", c)
	miso.StartRabbitMqClient(c)

	ol := OperateLog{
		OperateName:  "Fantahsea test operate log",
		OperateDesc:  "just a unit test",
		OperateTime:  miso.ETime(time.Now()),
		OperateParam: "{  }",
		Username:     "yongj.zhuang",
		UserId:       1,
	}
	err := DispatchOperateLog(c, ol)
	if err != nil {
		t.Fatal(err)
	}

}
