package aws

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/Crampustallin/discord_bot/internal/models"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type Aws struct {
	sess       *session.Session
	svc        *s3.S3
	bucketName string
}

func NewAws(region, bucket string) *Aws {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		panic(err)
	}
	svc := s3.New(sess)
	return &Aws{
		sess:       sess,
		svc:        svc,
		bucketName: bucket,
	}
}

func (a *Aws) Upload(fileName string) error {
	file, err := os.Open(fileName)
	fmt.Println(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = a.svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(a.bucketName),
		Key:    aws.String(fileName),
		Body:   file,
	})
	if err != nil {
		return errors.New("Unable to upload " + fileName + " to " + a.bucketName + ": " + err.Error())
	}
	return nil
}

func (a *Aws) GetObjUrl(key string) (*models.Url, error) {
	req, _ := a.svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(a.bucketName),
		Key:    aws.String(key),
	})
	expires := 24 * time.Hour
	url, err := req.Presign(expires)
	if err != nil {
		return nil, err
	}
	return &models.Url{Link: url, Expires: expires}, nil
}

func (a *Aws) GetObjList() ([]string, error) {
	req, err := a.svc.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(a.bucketName),
	})
	if err != nil {
		return nil, err
	}
	res := make([]string, len(req.Contents))
	for _, obj := range req.Contents {
		res = append(res, *obj.Key)
	}
	return res, nil
}
