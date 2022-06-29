package selective

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gocarina/gocsv"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Model struct {
	Id      string `csv:"Id"`
	Url     string `csv:"Url"`
	Destino string `csv:"Destino"`
}

func BlobtoS3(filename string, evidencia string) *manager.UploadOutput {
	resp, err := http.Get(evidencia)
	if err != nil {
		log.Fatalln(err)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	//UploadToS3
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalln(err)
	}
	//Creating a new S3 client
	client := s3.NewFromConfig(cfg)
	uploader := manager.NewUploader(client)
	result, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String("pruebas-devops-2022"),
		Key:    aws.String(filename),
		Body:   strings.NewReader(string(data)),
	})
	if err != nil {
		log.Fatalln(err)
	}
	return result
}

func ReadCSV(pathfile string) ([]Model, error) {
	csvfile, err := os.OpenFile(pathfile, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	defer csvfile.Close()
	var models []Model = []Model{}
	if err := gocsv.UnmarshalFile(csvfile, &models); err != nil {
		log.Fatalln(err)
		return nil, err
	}
	return models, nil
}

func ParsingUrl(urls ...string) ([]string, error) {
	var decodedurls []string
	for _, u := range urls {
		decoded, err := url.QueryUnescape(u)
		if err != nil {
			log.Println("error al decodear la url")
			return nil, err
		}
		decodedurls = append(decodedurls, decoded)
	}
	return decodedurls, nil
}
