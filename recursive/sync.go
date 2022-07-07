package recursive

import (
	"TerritoriumSync/helpers"
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"log"
)

var _, errlog *log.Logger = helpers.Logger()

func GetObjects(origin string, carpeta ...string) []string {
	var allobjects []string
	//load the shared AWS configuration in ~/.aws/config
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		errlog.Fatalln(err)
	}

	//Creating a new S3 client
	client := s3.NewFromConfig(cfg)

	//Get firstPage of results
	for _, f := range carpeta {
		output, err := client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
			Bucket: aws.String(fmt.Sprintf("%v", origin)),
			Prefix: aws.String(fmt.Sprintf("%v", f)),
		})
		if err != nil {
			errlog.Fatalln(err)
		}
		for _, object := range output.Contents {
			allobjects = append(allobjects, aws.ToString(object.Key))
		}
	}

	return allobjects
}

func CopyFiles(bucket string, files ...string) ([]byte, error) {
	//load the shared AWS configuration in ~/.aws/config
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		errlog.Fatalln(err)
	}
	//Creating a new S3 client
	client := s3.NewFromConfig(cfg)

	buffer := manager.NewWriteAtBuffer([]byte{})
	downloader := manager.NewDownloader(client)

	for _, file := range files {
		numbytes, err := downloader.Download(context.TODO(), buffer, &s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(file),
		})
		if err != nil {
			return nil, err
		}
		if numbytes < 1 {
			return nil, errors.New("Zero bytes written to memory")
		}
	}
	return buffer.Bytes(), nil
}
