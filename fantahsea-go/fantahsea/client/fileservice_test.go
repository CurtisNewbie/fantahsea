package client

import (
	"fmt"
	"testing"

	"github.com/curtisnewbie/github.com/curtisnewbie/fantahsea/config"
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
