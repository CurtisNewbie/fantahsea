package client

import (
	"fmt"
	"testing"

	"github.com/curtisnewbie/gocommon/config"
	log "github.com/sirupsen/logrus"
)

func TestDownloadFile(t *testing.T) {
	conf, _ := config.ParseJsonConfig(fmt.Sprintf("../app-conf-%v.json", "dev"))
	config.SetGlobalConfig(conf)

	fileKey := "4b10b75b-501e-4010-b6e7-e4724835e210"
	err := DownloadFile(fileKey, fmt.Sprintf("/tmp/%s.png", fileKey))
	if err != nil {
		t.Error(err)
	}
}

func TestValidateFileKey(t *testing.T) {
	conf, _ := config.ParseJsonConfig(fmt.Sprintf("../app-conf-%v.json", "dev"))
	config.SetGlobalConfig(conf)

	fileKey := "4b10b75b-501e-4010-b6e7-e4724835e210"
	userId := "30"
	hasAccess, err := ValidateFileKey(fileKey, userId)
	if err != nil {
		t.Error(err)
	}
	if !hasAccess {
		t.Errorf("User %s should have access to file %s", userId, fileKey)
	}

}

func TestGetFileInfo(t *testing.T) {
	conf, _ := config.ParseJsonConfig(fmt.Sprintf("../app-conf-%v.json", "dev"))
	config.SetGlobalConfig(conf)

	fileKey := "4b10b75b-501e-4010-b6e7-e4724835e210"
	resp, err := GetFileInfo(fileKey)
	if err != nil {
		t.Error(err)
	}
	if resp.Data == nil {
		t.Error("Resp doesn't contain data")
	}
	log.Infof("Normal Resp.Data: %+v", resp.Data)

	fileKey = "non-existing-file-key"
	resp, err = GetFileInfo(fileKey)
	if err == nil {
		t.Error("Should have received error")
	}
	log.Infof("Missing Resp: %+v", resp)
}

func TestListFilesInDir(t *testing.T) {
	conf, _ := config.ParseJsonConfig(fmt.Sprintf("../app-conf-%v.json", "dev"))
	config.SetGlobalConfig(conf)

	fileKey := "5ddf49ca-dec9-4ecf-962d-47b0f3eab90c"
	resp, err := ListFilesInDir(fileKey, 100, 1)
	if err != nil {
		t.Error(err)
	}
	if resp.Data == nil {
		t.Error("Resp doesn't contain data")
	}
	log.Infof("Resp.Data: %+v", resp.Data)
}
