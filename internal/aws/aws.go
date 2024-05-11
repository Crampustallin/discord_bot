package aws

import (
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type Aws struct {
	sess       *session.Session
	uploader   *s3manager.Uploader
	bucketName string
}

func NewAws(region, bucket string) *Aws {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		panic(err)
	}
	uploader := s3manager.NewUploader(sess)
	return &Aws{
		sess:       sess,
		uploader:   uploader,
		bucketName: bucket,
	}
}

func (a *Aws) Upload(fileName string) error {
	file, err := os.Open(fileName + ".ogg")
	fmt.Println(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = a.uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(a.bucketName),
		Key:    aws.String(fileName),
		Body:   file,
	})
	if err != nil {
		return errors.New("Unable to upload " + fileName + " to " + a.bucketName + ": " + err.Error())
	}
	return nil
}
