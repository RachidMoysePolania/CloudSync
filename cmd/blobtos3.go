package cmd

import (
	"TerritoriumSync/selective"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strings"
	"time"
)

var pathfile string
var buckname string
var local bool
var blobtos3 = &cobra.Command{
	Use:   "aztos3",
	Short: "command to transmit data from blobstorage to s3 bucket",
	Long:  "upload recursively data from blobstorage in azure to s3 bucket",
	Run: func(cmd *cobra.Command, args []string) {
		globalstart := time.Now()
		models, err := selective.ReadCSV(pathfile)
		if err != nil {
			log.Fatalln(err)
		}
		if local {
			for _, data := range models {
				destinourl, err := selective.ParsingUrl(data.Destino)
				if err != nil {
					log.Fatalln(err)
				}
				filename := strings.Split(destinourl[0], "/")
				err = os.MkdirAll(strings.Join(filename[:len(filename)-1], "/"), 0755)
				if err != nil {
					log.Fatalln(err)
				}
				data := selective.LocalStore(data.Url)
				f, err := os.Create(strings.Join(filename[:len(filename)-1], "/") + "/" + filename[len(filename)-1])
				if err != nil {
					log.Fatalln(err)
				}
				defer f.Close()
				f.Write(data)
			}
			log.Println("Tarea completa")
			time.Sleep(time.Minute * 2)
			os.Exit(1)
		}
		for _, data := range models {
			start := time.Now()
			parsedurl, err := selective.ParsingUrl(data.Destino)
			if err != nil {
				log.Fatalln(err)
			}
			result := selective.BlobtoS3(parsedurl[0], data.Url, buckname)
			log.Println(fmt.Sprintf("Item Uploaded %v Time Elapsed: %v", result.Location, time.Since(start)))
		}
		log.Println(fmt.Sprintf("Task Completed %v", time.Since(globalstart)))
	},
}

func init() {
	rootCmd.AddCommand(blobtos3)
	blobtos3.Flags().StringVarP(&pathfile, "csvfile", "p", ".", "Use this parameter to set the the path to CSV file to upload")
	blobtos3.Flags().StringVarP(&buckname, "bucketname", "b", "pruebas-devops-2022", "Use this parameter to set the bucket where upload the data")
	blobtos3.Flags().BoolVarP(&local, "local", "l", false, "Enable this feature to store from blobstorage in local")
}
