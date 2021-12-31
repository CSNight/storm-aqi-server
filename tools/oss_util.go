package tools

import (
	"bytes"
	"context"
	"errors"
	"github.com/minio/minio-go/v7"
	"io"
	"net/http"
	"net/url"
	"time"
)

var bucketName = "sys-image"
var httpClient = http.Client{
	Transport: &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse("http://172.16.1.135:3128")
		},
		ForceAttemptHTTP2:   false,
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
		MaxConnsPerHost:     2000,
	},
	Timeout: time.Minute,
}

func GetObject(cli *minio.Client, name string) ([]byte, error) {
	ctx := context.Background()
	object, err := cli.GetObject(ctx, bucketName, "aqi/"+name, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer object.Close()
	data, err := io.ReadAll(object)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func PutObject(cli *minio.Client, object []byte, name string) bool {
	ctx := context.Background()
	_, err := cli.PutObject(ctx, bucketName, name, bytes.NewReader(object), int64(len(object)),
		minio.PutObjectOptions{
			ContentEncoding: "utf-8",
			ContentType:     "image/png",
		})
	if err != nil {
		return false
	}
	return true
}

func ExistObject(cli *minio.Client, objectName string) bool {
	ctx := context.Background()
	info, err := cli.StatObject(ctx, bucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		return false
	}
	if info.Size > 0 {
		return true
	}
	return false
}

func DownloadImage(url string) ([]byte, error) {
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New("request failed")
	}
	if resp.Body != nil {
		defer func() {
			_ = resp.Body.Close()
		}()
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if len(body) == 0 {
		return nil, errors.New("empty body")
	}
	return body, nil
}
