package client

import (
	"io"
	"net/url"
	"os"

	"github.com/curtisnewbie/gocommon/client"
	"github.com/curtisnewbie/gocommon/common"
)

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

// Download file from mini-fstore
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
