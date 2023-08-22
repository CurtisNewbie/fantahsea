package data

import (
	"testing"

	"github.com/curtisnewbie/gocommon/common"
	"github.com/curtisnewbie/gocommon/consul"
	"github.com/curtisnewbie/gocommon/mysql"
	"github.com/curtisnewbie/gocommon/redis"
)

func TestListGalleryImages(t *testing.T) {
	c := common.EmptyRail()
	common.LoadConfigFromFile("../app-conf-dev.yml", c)
	redis.InitRedisFromProp()
	if _, e := consul.GetConsulClient(); e != nil {
		t.Fatal(e)
	}
	if e := mysql.InitMySqlFromProp(); e != nil {
		t.Fatal(e)
	}

	cmd := ListGalleryImagesCmd{GalleryNo: "GALZRQG0RP8KPMUU0HQ4P7N7LACG", Paging: common.Paging{Limit: 5, Page: 1}}
	user := common.User{UserId: 1, UserNo: "UE202205142310076187414", Username: "zhuangyongj"}
	r, err := ListGalleryImages(common.EmptyRail(), cmd, user)
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Images) < 1 {
		t.Fatalf("Images is empty")
	}
	t.Logf("%+v", r.Images)
}
