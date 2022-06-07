package db

import (
	"context"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/minio/minio-go/v7"
)

type ImageResponse struct {
	Data string `json:"data"`
	Max  string `json:"max"`
	Min  string `json:"min"`
}

var bucket = "silam"

func (db *DB) GetImage(tm string, pol string) (*ImageResponse, error) {
	objectDir := tm[0:10]
	tf := strings.ReplaceAll(tm, ":", "$")
	objectName := fmt.Sprintf("silam_AQ_%s_%s.png", pol, tf)
	tagging, err := db.oss.GetObjectTagging(context.Background(), bucket, objectDir+"/"+objectName, minio.GetObjectTaggingOptions{})
	if err != nil {
		return nil, err
	}
	tags := tagging.ToMap()
	resp := &ImageResponse{
		Data: objectDir + "/" + objectName,
		Max:  tags["max"],
		Min:  tags["min"],
	}
	return resp, nil
}

func (db *DB) DownloadImage(dir string, file string) ([]byte, error) {
	fmt.Println(file)
	imgObj, err := db.oss.GetObject(context.Background(), bucket, dir+"/"+file, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer imgObj.Close()
	body, err := ioutil.ReadAll(imgObj)
	if err != nil {
		return nil, err
	}
	return body, nil
}
