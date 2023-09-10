package data

import (
	"testing"

	"github.com/curtisnewbie/gocommon/common"
	"github.com/curtisnewbie/miso/miso"
)

func TestListGalleryImages(t *testing.T) {
	c := miso.EmptyRail()
	miso.LoadConfigFromFile("../app-conf-dev.yml", c)
	miso.InitRedisFromProp(c)
	if _, e := miso.GetConsulClient(); e != nil {
		t.Fatal(e)
	}
	if e := miso.InitMySQLFromProp(); e != nil {
		t.Fatal(e)
	}

	cmd := ListGalleryImagesCmd{GalleryNo: "GALZRQG0RP8KPMUU0HQ4P7N7LACG", Paging: miso.Paging{Limit: 5, Page: 1}}
	user := common.User{UserId: 1, UserNo: "UE202205142310076187414", Username: "zhuangyongj"}
	r, err := ListGalleryImages(miso.EmptyRail(), cmd, user)
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Images) < 1 {
		t.Fatalf("Images is empty")
	}
	t.Logf("%+v", r.Images)
}
