package storage

import (
	"mime/multipart"
	"time"

	"github.com/aws/aws-sdk-go/private/protocol/rest"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type StorageService interface {
	Upload(key string, file *multipart.FileHeader) error
	GetUrl(key string) string
}

// impl `StorageService` start

func (storage *StorageModule) Upload(key string, file *multipart.FileHeader) error {
	reader, err := file.Open()

	if err != nil {
		return err
	}
	defer reader.Close()

	_, err = storage.uploader.Upload(&s3manager.UploadInput{
		Bucket: &storage.bucket,
		Key:    &key,
		Body:   reader,
	})

	return err
}

func (storage *StorageModule) GetUrl(key string) string {
	req, _ := storage.s3.GetObjectRequest(&s3.GetObjectInput{
		Bucket: &storage.bucket,
		Key:    &key,
	})

	if storage.enablePresign {
		urlStr, err := req.Presign(15 * time.Minute)

		if err == nil {
			return urlStr
		}
	}

	rest.Build(req)

	return req.HTTPRequest.URL.String()
}

// impl `StorageService` end
