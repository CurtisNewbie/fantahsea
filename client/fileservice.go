package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/curtisnewbie/gocommon/config"
	"github.com/curtisnewbie/gocommon/web/dto"
	"github.com/curtisnewbie/gocommon/weberr"
	log "github.com/sirupsen/logrus"
)

const (
	DIR  FileType = "DIR"
	FILE FileType = "FILE"
)

type FileType string

type ValidateFileKeyResp struct {
	dto.Resp
	Data bool `json:"data"`
}

type FileInfoResp struct {

	/** name of the file */
	Name string `json:"name"`

	/** file's uuid */
	Uuid string `json:"uuid"`

	/** size of file in bytes */
	SizeInBytes int64 `json:"sizeInBytes"`

	/** uploader id, i.e., user.id */
	UploaderId int `json:"uploaderId"`

	/** uploader name */
	UploaderName string `json:"uploaderName"`

	/** when the file is deleted */
	IsDeleted bool `json:"isDeleted"`

	/** file type: FILE, DIR */
	FileType FileType `json:"fileType"`

	/** parent file's uuid */
	ParentFile string `json:"parentFile"`
}

type GetFileInfoResp struct {
	dto.Resp
	Data *FileInfoResp `json:"data"`
}

type ListFilesInDirResp struct {
	dto.Resp
	// list of file key
	Data []string `json:"data"`
}

// List files in dir from file-service
func ListFilesInDir(fileKey string, limit int, page int) (*ListFilesInDirResp, error) {
	url := BuildFileServiceUrl(fmt.Sprintf("/remote/user/file/indir/list?fileKey=%s&limit=%d&page=%d", fileKey, limit, page))
	log.Infof("List files in dir, url: %s", url)
	r, e := http.Get(url)
	if e != nil {
		return nil, e
	}
	defer r.Body.Close()

	body, e := io.ReadAll(r.Body)
	if e != nil {
		return nil, e
	}
	log.Infof("List files in dir, resp: %v", string(body))

	var resp ListFilesInDirResp
	if e = json.Unmarshal(body, &resp); e != nil {
		return nil, e
	}

	if resp.Resp.Error {
		return nil, weberr.NewWebErr(resp.Resp.Msg)
	}
	return &resp, nil
}

// Get file info from file-service
func GetFileInfo(fileKey string) (*GetFileInfoResp, error) {
	url := BuildFileServiceUrl(fmt.Sprintf("/remote/user/file/info?fileKey=%s", fileKey))
	log.Infof("Get file info, url: %s", url)
	r, e := http.Get(url)
	if e != nil {
		return nil, e
	}
	defer r.Body.Close()

	body, e := io.ReadAll(r.Body)
	if e != nil {
		return nil, e
	}
	log.Infof("Get file info, resp: %v", string(body))

	var resp GetFileInfoResp
	if e = json.Unmarshal(body, &resp); e != nil {
		return nil, e
	}

	if resp.Resp.Error {
		return nil, weberr.NewWebErr(resp.Resp.Msg)
	}
	return &resp, nil
}

// Download file from file-service
func DownloadFile(fileKey string, absPath string) error {
	url := BuildFileServiceUrl(fmt.Sprintf("/remote/user/file/download?fileKey=%s", fileKey))
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

	url := BuildFileServiceUrl(fmt.Sprintf("/remote/user/file/owner/validation?fileKey=%s&userId=%s", fileKey, userId))
	log.Infof("Validate file key, url: %s", url)

	r, e := http.Get(url)
	if e != nil {
		return false, e
	}
	defer r.Body.Close()

	body, e := io.ReadAll(r.Body)
	if e != nil {
		return false, e
	}
	log.Infof("Validate file key, key: %v resp: %v", fileKey, string(body))

	var resp ValidateFileKeyResp
	if e := json.Unmarshal(body, &resp); e != nil {
		return false, e
	}

	if resp.Resp.Error {
		return false, weberr.NewWebErr(resp.Resp.Msg)
	}

	return resp.Data, nil
}

// Concatenate given relative url to base url, the relUrl may or may not start with "/"
func BuildFileServiceUrl(relUrl string) string {
	if !strings.HasPrefix(relUrl, "/") {
		relUrl = "/" + relUrl
	}
	return config.GlobalConfig.ClientConf.FileServiceUrl + relUrl
}
