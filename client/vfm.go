package client

import (
	"fmt"
	"strconv"

	"github.com/curtisnewbie/miso/miso"
)

type FileType string

const (
	DIR  FileType = "DIR"
	FILE FileType = "FILE"
)

type ListFilesInDirResp struct {
	miso.Resp
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
	miso.Resp
	Data *FileInfoResp `json:"data"`
}

// Get file info from file-service
func GetFileInfo(c miso.Rail, fileKey string) (*GetFileInfoResp, error) {

	var resp GetFileInfoResp
	err := miso.NewDynTClient(c, "/remote/user/file/info", "vfm").
		EnableTracing().
		AddQueryParams("fileKey", fileKey).
		Get().
		Json(&resp)

	if err != nil {
		return nil, err
	}

	if resp.Resp.Error {
		return nil, miso.NewErr(resp.Resp.Msg)
	}
	return &resp, nil
}

// Validate the file key, return true if it's valid else false
func ValidateFileKey(c miso.Rail, fileKey string, userId int) (bool, error) {

	var resp ValidateFileKeyResp
	err := miso.NewDynTClient(c, "/remote/user/file/owner/validation", "vfm").
		EnableTracing().
		AddQueryParams("fileKey", fileKey).
		AddQueryParams("userId", fmt.Sprintf("%v", userId)).
		Get().
		Json(&resp)

	if err != nil {
		return false, err
	}
	if resp.Error {
		return false, miso.NewErr(resp.Resp.Msg)
	}

	return resp.Data, nil
}

type ValidateFileKeyResp struct {
	miso.Resp
	Data bool `json:"data"`
}

// List files in dir from vfm
func ListFilesInDir(c miso.Rail, fileKey string, limit int, page int) (*ListFilesInDirResp, error) {
	var resp ListFilesInDirResp
	err := miso.NewDynTClient(c, "/remote/user/file/indir/list", "vfm").
		EnableTracing().
		AddQueryParams("fileKey", fileKey).
		AddQueryParams("limit", strconv.Itoa(limit)).
		AddQueryParams("page", strconv.Itoa(page)).
		Get().
		Json(&resp)

	if err != nil {
		return nil, err
	}
	if resp.Error {
		return nil, miso.NewErr(resp.Resp.Msg)
	}
	return &resp, nil
}
