package db

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"strings"
	"time"
)

type ImageResponse struct {
	Overlay string `json:"overlay"`
	Data    string `json:"data"`
	Max     string `json:"max"`
	Min     string `json:"min"`
}

var bucket = "silam"

func (db *DB) GetImage(tm string, pol string) (*ImageResponse, error) {
	objectDir := tm[0:10]
	tf := strings.ReplaceAll(tm, ":", "$")
	objectName := fmt.Sprintf("silam_AQ_%s_%s.png", pol, tf)
	dataName := fmt.Sprintf("silam_AQ_%s_%s.jpeg", pol, tf)
	imgUrl, err := db.oss.PresignedGetObject(context.Background(), bucket, objectDir+"/"+objectName, time.Hour, nil)
	if err != nil {
		return nil, err
	}
	dataUrl, err := db.oss.PresignedGetObject(context.Background(), bucket, objectDir+"/"+dataName, time.Hour, nil)
	if err != nil {
		return nil, err
	}
	tagging, err := db.oss.GetObjectTagging(context.Background(), bucket, objectDir+"/"+dataName, minio.GetObjectTaggingOptions{})
	if err != nil {
		return nil, err
	}
	tags := tagging.ToMap()
	resp := &ImageResponse{
		Overlay: imgUrl.String(),
		Data:    dataUrl.String(),
		Max:     tags["max"],
		Min:     tags["min"],
	}
	return resp, nil
}
