package fclient

import (
	"context"
	"io"
	"os"

	"github.com/curtisnewbie/gocommon/client"
	"github.com/curtisnewbie/gocommon/common"
	"github.com/curtisnewbie/gocommon/consul"
	"github.com/sirupsen/logrus"
)

const (
	DIR               FileType = "DIR"
	FILE              FileType = "FILE"
	FILE_SERVICE_NAME string   = "file-service"
)

type FileType string

type ValidateFileKeyResp struct {
	common.Resp
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
	common.Resp
	Data *FileInfoResp `json:"data"`
}

type ListFilesInDirResp struct {
	common.Resp
	// list of file key
	Data []string `json:"data"`
}

// List files in dir from file-service
func ListFilesInDir(ctx context.Context, fileKey string, limit int, page int) (*ListFilesInDirResp, error) {
	url := consul.ResolveRequestUrl(FILE_SERVICE_NAME, "/remote/user/file/indir/list")
	r := client.NewDefaultTClient(ctx, url).
		EnableTracing().
		Get(map[string][]string{
			"fileKey": {fileKey},
			"limit":   {string(limit)},
			"page":    {string(page)},
		})
	defer r.Close()

	var resp ListFilesInDirResp
	if e := r.ReadJson(&resp); e != nil {
		return nil, e
	}

	if resp.Resp.Error {
		return nil, common.NewWebErr(resp.Resp.Msg)
	}
	return &resp, nil
}

// Get file info from file-service
func GetFileInfo(ctx context.Context, fileKey string) (*GetFileInfoResp, error) {
	url := consul.ResolveRequestUrl(FILE_SERVICE_NAME, "/remote/user/file/info")
	r := client.NewDefaultTClient(ctx, url).
		EnableTracing().
		Get(map[string][]string{
			"fileKey": {fileKey},
		})
	defer r.Close()

	var resp GetFileInfoResp
	if e := r.ReadJson(&resp); e != nil {
		return nil, e
	}

	if resp.Resp.Error {
		return nil, common.NewWebErr(resp.Resp.Msg)
	}
	return &resp, nil
}

// Download file from file-service
func DownloadFile(ctx context.Context, fileKey string, absPath string) error {
	url := consul.ResolveRequestUrl(FILE_SERVICE_NAME, "/remote/user/file/download")
	r := client.NewDefaultTClient(ctx, url).
		EnableTracing().
		Get(map[string][]string{
			"fileKey": {fileKey},
		})
	defer r.Close()

	out, err := os.Create(absPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, r.Resp.Body)
	if err != nil {
		return err
	}

	logrus.Infof("Finished downloading file, url: %s", url)
	return nil
}

// Validate the file key, return true if it's valid else false
func ValidateFileKey(ctx context.Context, fileKey string, userId string) (bool, error) {
	url := consul.ResolveRequestUrl(FILE_SERVICE_NAME, "/remote/user/file/owner/validation")
	r := client.NewDefaultTClient(ctx, url).
		EnableTracing().
		Get(map[string][]string{
			"fileKey": {fileKey},
			"userId":  {userId},
		})
	defer r.Close()

	var resp ValidateFileKeyResp
	if e := r.ReadJson(&resp); e != nil {
		return false, e
	}

	if resp.Resp.Error {
		return false, common.NewWebErr(resp.Resp.Msg)
	}

	return resp.Data, nil
}
