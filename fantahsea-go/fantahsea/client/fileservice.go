package client

import (
	"encoding/json"
	"fantahsea/config"
	"fantahsea/web/dto"
	"fmt"
	"io"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
)

type ValidateFileKeyResp struct {
	dto.Resp
	Data bool `json:"data"`
}

// Download file from file-service
func DownloadFile(fileKey string, absPath string) error {
	base := config.GlobalConfig.ClientConf.FileServiceUrl
	url := base + fmt.Sprintf("/remote/user/file/download?fileKey=%s", fileKey)
	log.Infof("Download file, url: %s, absPath: %s", url, absPath)

	out, err := os.Create(absPath)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	log.Infof("Finished downloading file, url: %s", url)
	return nil
}

// Validate the file key, return true if it's valid else false
func ValidateFileKey(fileKey string, userId string) (bool, error) {
	base := config.GlobalConfig.ClientConf.FileServiceUrl
	url := base + fmt.Sprintf("/remote/user/file/owner/validation?fileKey=%s&userId=%s", fileKey, userId)
	log.Infof("Validate file key, url: %s", url)

	r, e := http.Get(url)
	defer r.Body.Close()
	if e != nil {
		return false, e
	}

	body, e := io.ReadAll(r.Body)
	if e != nil {
		return false, e
	}

	var resp ValidateFileKeyResp
	if e := json.Unmarshal(body, &resp); e != nil {
		return false, e
	}

	return resp.Data, nil
}
