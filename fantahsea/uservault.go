package fantahsea

import (
	"fmt"

	"github.com/curtisnewbie/miso/miso"
)

type FindUserReq struct {
	UserId   *int    `json:"userId"`
	UserNo   *string `json:"userNo"`
	Username *string `json:"username"`
}

type UserInfo struct {
	Id       int
	Username string
	UserNo   string
}

func FindUser(rail miso.Rail, req FindUserReq) (UserInfo, error) {
	var r miso.GnResp[UserInfo]
	err := miso.NewDynTClient(rail, "/remote/user/info", "user-vault").
		Require2xx().
		PostJson(req).
		Json(&r)
	if err != nil {
		return UserInfo{}, fmt.Errorf("failed to find user (user-vault), %v", err)
	}
	return r.Res()
}
