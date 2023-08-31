package client

import (
	"io"
	"net/url"
	"os"

	"github.com/curtisnewbie/miso/client"
	"github.com/curtisnewbie/miso/core"
)

type GenFileTempTokenReq struct {
	Filekeys    []string `json:"fileKeys"`
	ExpireInMin int      `json:"expireInMin"`
}

type GenFileTempTokenResp struct {
	core.Resp
	Data map[string]string `json:"data"`
}

type BatchGenFileKeyReq struct {
	Items []BatchGenFileKeyItem `json:"items"`
}

type BatchGenFileKeyItem struct {
	FileId   string `json:"fileId"`
	Filename string `json:"filename"`
}

type BatchGenFileKeyResp struct {
	FileId  string `json:"fileId"`
	TempKey string `json:"tempKey"`
}

func BatchGetFstoreTmpToken(c core.Rail, req BatchGenFileKeyReq) ([]BatchGenFileKeyResp, error) {
	r := client.NewDynTClient(c, "/file/key/batch", "fstore").
		EnableTracing().
		PostJson(&req)
	if r.Err != nil {
		return nil, r.Err
	}
	defer r.Close()

	var res core.GnResp[[]BatchGenFileKeyResp]
	if e := r.ReadJson(&res); e != nil {
		return nil, e
	}

	if res.Error {
		return nil, res.Err()
	}
	return res.Data, nil
}

func GetFstoreTmpToken(c core.Rail, fileId string, filename string) (string, error) {
	r := client.NewDynTClient(c, "/file/key", "fstore").
		EnableTracing().
		AddQueryParams("fileId", fileId).
		AddQueryParams("filename", url.QueryEscape(filename)).
		Get()
	if r.Err != nil {
		return "", r.Err
	}
	defer r.Close()

	var res core.GnResp[string]
	if e := r.ReadJson(&res); e != nil {
		return "", e
	}

	if res.Error {
		return "", res.Err()
	}
	return res.Data, nil
}

// Download file from mini-fstore
func DownloadFile(c core.Rail, tmpToken string, absPath string) error {
	r := client.NewDynTClient(c, "/file/raw", "fstore").
		EnableTracing().
		AddQueryParams("key", tmpToken).
		Get()
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
