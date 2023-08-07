package client

import (
	"strconv"

	"github.com/curtisnewbie/gocommon/client"
	"github.com/curtisnewbie/gocommon/common"
)

type FileType string

const (
	DIR  FileType = "DIR"
	FILE FileType = "FILE"
)

type ListFilesInDirResp struct {
	common.Resp
	// list of file key
	Data []string `json:"data"`
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

	LocalPath string `json:"localPath"`

	FstoreFileId string `json:"fstoreFileId"`

	Thumbnail string `json:"thumbnail"` // also a mini-fstore file_id
}

type GetFileInfoResp struct {
	common.Resp
	Data *FileInfoResp `json:"data"`
}

// Get file info from file-service
func GetFileInfo(c common.ExecContext, fileKey string) (*GetFileInfoResp, error) {
	r := client.NewDynTClient(c, "/remote/user/file/info", "vfm").
		EnableTracing().
		AddQueryParams("fileKey", fileKey).
		Get()
	defer r.Close()

	if r.Err != nil {
		return nil, r.Err
	}

	var resp GetFileInfoResp
	if e := r.ReadJson(&resp); e != nil {
		return nil, e
	}

	if resp.Resp.Error {
		return nil, common.NewWebErr(resp.Resp.Msg)
	}
	return &resp, nil
}

// Validate the file key, return true if it's valid else false
func ValidateFileKey(c common.ExecContext, fileKey string, userId string) (bool, error) {
	r := client.NewDynTClient(c, "/remote/user/file/owner/validation", "vfm").
		EnableTracing().
		AddQueryParams("fileKey", fileKey).
		AddQueryParams("userId", userId).
		Get()
	defer r.Close()

	if r.Err != nil {
		return false, r.Err
	}

	var resp ValidateFileKeyResp
	if e := r.ReadJson(&resp); e != nil {
		return false, e
	}

	if resp.Error {
		return false, common.NewWebErr(resp.Resp.Msg)
	}

	return resp.Data, nil
}

type ValidateFileKeyResp struct {
	common.Resp
	Data bool `json:"data"`
}

// List files in dir from vfm
func ListFilesInDir(c common.ExecContext, fileKey string, limit int, page int) (*ListFilesInDirResp, error) {

	r := client.NewDynTClient(c, "/remote/user/file/indir/list", "vfm").
		EnableTracing().
		AddQueryParams("fileKey", fileKey).
		AddQueryParams("limit", strconv.Itoa(limit)).
		AddQueryParams("page", strconv.Itoa(page)).
		Get()
	defer r.Close()

	if r.Err != nil {
		return nil, r.Err
	}

	var resp ListFilesInDirResp
	if e := r.ReadJson(&resp); e != nil {
		return nil, e
	}

	if resp.Error {
		return nil, common.NewWebErr(resp.Resp.Msg)
	}
	return &resp, nil
}
