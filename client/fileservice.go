package client

import (
	"io"
	"net/url"
	"os"
	"strconv"

	"github.com/curtisnewbie/gocommon/client"
	"github.com/curtisnewbie/gocommon/common"
	"github.com/curtisnewbie/gocommon/consul"
)

const (
	DIR               FileType = "DIR"
	FILE              FileType = "FILE"
	FILE_SERVICE_NAME string   = "vfm"
	EXP_MIN                    = 15 // expiration time of the token in minutes
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

	LocalPath string `json:"localPath"`

	FstoreFileId string `json:"fstoreFileId"`
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

type GenFileTempTokenReq struct {
	Filekeys    []string `json:"fileKeys"`
	ExpireInMin int      `json:"expireInMin"`
}

type GenFileTempTokenResp struct {
	common.Resp
	Data map[string]string `json:"data"`
}

func GetFstoreTmpToken(c common.ExecContext, fileId string, filename string) (string, error) {
	r := client.NewDynTClient(c, "/file/key", "fstore").
		EnableTracing().
		EnableRequestLog().
		Get(map[string][]string{"fileId": {fileId}, "filename": {url.QueryEscape(filename)}})
	if r.Err != nil {
		return "", r.Err
	}
	defer r.Close()

	var res common.GnResp[string]
	if e := r.ReadJson(&res); e != nil {
		return "", e
	}

	if res.Error {
		return "", res.Err()
	}
	return res.Data, nil
}

// List files in dir from file-service
func ListFilesInDir(c common.ExecContext, fileKey string, limit int, page int) (*ListFilesInDirResp, error) {
	url, e := consul.ResolveRequestUrl(FILE_SERVICE_NAME, "/remote/user/file/indir/list")
	if e != nil {
		return nil, e
	}
	slimit := strconv.Itoa(limit)
	plimit := strconv.Itoa(page)

	r := client.NewDefaultTClient(c, url).
		EnableTracing().
		Get(map[string][]string{
			"fileKey": {fileKey},
			"limit":   {slimit},
			"page":    {plimit},
		})
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

// Get file info from file-service
func GetFileInfo(c common.ExecContext, fileKey string) (*GetFileInfoResp, error) {
	r := client.NewDynTClient(c, "/remote/user/file/info", "vfm").
		EnableTracing().
		Get(map[string][]string{
			"fileKey": {fileKey},
		})
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

// Download file from file-service
func DownloadFile(c common.ExecContext, tmpToken string, absPath string) error {
	r := client.NewDynTClient(c, "/file/raw", "fstore").
		EnableTracing().
		Get(map[string][]string{
			"key": {tmpToken},
		})
	if r.Err != nil {
		return r.Err
	}
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
	return nil
}

// Validate the file key, return true if it's valid else false
func ValidateFileKey(c common.ExecContext, fileKey string, userId string) (bool, error) {
	r := client.NewDynTClient(c, "/remote/user/file/owner/validation", "vfm").
		EnableTracing().
		Get(map[string][]string{
			"fileKey": {fileKey},
			"userId":  {userId},
		})
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
