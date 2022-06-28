package cmd

import (
	"TerritoriumSync/recursive"
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strings"
	"time"
)

var bucketname string
var prefixfolder []string
var azdestination bool
var awsdestination bool
var localdestination bool
var workingdirectory string
var downloadCmd = &cobra.Command{
	Use:   "awsdownload",
	Short: "download all files from aws s3",
	Long:  "Use this module for list and download recursive folders from aws s3",
	Run: func(cmd *cobra.Command, args []string) {
		//Logger
		logfile, err := os.OpenFile("info.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalln("Error al crear el logger")
		}
		defer logfile.Close()

		logger := log.New(logfile, "", log.LstdFlags)
		files := recursive.GetObjects(bucketname, prefixfolder...)
		if localdestination {
			os.Chdir(workingdirectory)
			for _, archivo := range files {
				filename := strings.Split(archivo, "/")
				if filename[len(filename)-1] != "" {
					log.Println(fmt.Sprintf("[!] Descargando el archivo %v en la ruta %v", filename[len(filename)-1], strings.Join(filename[:len(filename)-1], "/")))
					logger.Println(fmt.Sprintf("[!] Descargando el archivo %v en la ruta %v", filename[len(filename)-1], strings.Join(filename[:len(filename)-1], "/")))
					err := os.MkdirAll(strings.Join(filename[:len(filename)-1], "/"), 0755)
					if err != nil {
						logger.Fatalln("Error al crear las carpetas de destino", err)
					}
					file, err := os.Create(strings.Join(filename[:len(filename)-1], "/") + "/" + filename[len(filename)-1])
					if err != nil {
						logger.Fatalln(err)
					}
					defer file.Close()
					//Download file
					data, err := recursive.CopyFiles(bucketname, archivo)
					if err != nil {
						log.Println(fmt.Sprintf("[x] Error al descargar el archivo %v, fallo con el siguiente error: %v", filename[len(filename)-1], err))
						logger.Fatalln(fmt.Sprintf("[x] Error al descargar el archivo %v, fallo con el siguiente error: %v", filename[len(filename)-1], err))
					}

					file.Write(data)
					log.Println(fmt.Sprintf("[+] Finalizada correctamente la descarga del archivo %v", filename[len(filename)-1]))
					logger.Println(fmt.Sprintf("[+] Finalizada correctamente la descarga del archivo %v", filename[len(filename)-1]))
					time.Sleep(time.Millisecond * 100)
				}
			}
		}

		if azdestination {
			var StorageAccountName string
			fmt.Print("Enter your Storage Account name where the data will be stored -> ")
			fmt.Scan(&StorageAccountName)

			var containername string
			fmt.Print("Enter your container name where the data will be stored -> ")
			fmt.Scan(&containername)

			url := fmt.Sprintf("https://%v.blob.core.windows.net/", StorageAccountName)
			credential, err := azidentity.NewDefaultAzureCredential(nil)
			if err != nil {
				log.Fatal("Invalid credentials with error: ", err)
			}
			for _, archivo := range files {
				filename := strings.Split(archivo, "/")
				if filename[len(filename)-1] != "" {
					blobclient, err := azblob.NewBlockBlobClient(url+containername+"/"+filename[len(filename)-1], credential, nil)
					if err != nil {
						log.Println("Error al crear el cliente blobstorage", err)
						logger.Fatalln("Error al crear el cliente blobstorage", err)
					}
					data, err := recursive.CopyFiles(bucketname, archivo)
					if err != nil {
						log.Println(fmt.Sprintf("[x] Error al descargar el archivo %v, fallo con el siguiente error: %v", filename[len(filename)-1], err))
						logger.Fatalln(fmt.Sprintf("[x] Error al descargar el archivo %v, fallo con el siguiente error: %v", filename[len(filename)-1], err))
					}
					_, err = blobclient.UploadBuffer(context.Background(), data, azblob.UploadOption{
						BlockSize: 1024,
					})
					if err != nil {
						log.Println("Error al subir el archivo al blobstorage", err)
						logger.Fatalln("Error al subir el archivo al blobstorage", err)
					}
					time.Sleep(time.Millisecond * 100)
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)
	downloadCmd.Flags().StringVarP(&bucketname, "bucketname", "b", "", "The bucketname for copy recursive file")
	downloadCmd.Flags().StringSliceVarP(&prefixfolder, "foldername", "f", nil, "The foler or folders for copy recursive file, comma separated for multiple folders")
	downloadCmd.Flags().BoolVarP(&azdestination, "Azure", "a", true, "By default transmit all data to azure")
	downloadCmd.Flags().BoolVarP(&awsdestination, "AWS", "w", false, "transmit data to AWS")
	downloadCmd.Flags().BoolVarP(&localdestination, "Local", "l", false, "Download data to local drive")
	downloadCmd.Flags().StringVarP(&workingdirectory, "savepath", "s", ".", "Set the directory where you need to save the data")
}
