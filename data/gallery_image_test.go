package data

import (
	"context"
	"testing"

	"github.com/curtisnewbie/gocommon/common"
	"github.com/curtisnewbie/gocommon/consul"
	"github.com/curtisnewbie/gocommon/mysql"
	"github.com/curtisnewbie/gocommon/redis"
)

func TestGuessIsImage(t *testing.T) {
	if isImage := guessIsImage("someFile.pdf", 1000); isImage {
		t.Errorf("someFile.pdf shouldn't be considered as an image")
	}

	if isImage := guessIsImage("someFile.jpg", 1000); !isImage {
		t.Errorf("someFile.jpg should be considered as an image")
	}

	if isImage := guessIsImage("oversized.jpg", 31*1048576); isImage {
		t.Errorf("oversized.jpg shouldn't be considered as an image")
	}
}

func TestListGalleryImages(t *testing.T) {
	common.LoadConfigFromFile("../app-conf-dev.yml")
	redis.InitRedisFromProp()
	consul.MustInitConsulClient()
	mysql.MustInitMySqlFromProp()

	cmd := ListGalleryImagesCmd{GalleryNo: "GALZRQG0RP8KPMUU0HQ4P7N7LACG", Paging: common.Paging{Limit: 5, Page: 1}}
	user := common.User{UserId: "1", UserNo: "UE202205142310076187414", Username: "zhuangyongj", Role: "admin"}
	r, err := ListGalleryImages(cmd, common.NewExecContext(context.Background(), &user))
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Images) < 1 {
		t.Fatalf("Images is empty")
	}
	t.Logf("%+v", r.Images)
}
