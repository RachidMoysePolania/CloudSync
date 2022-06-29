package cmd

import (
	"TerritoriumSync/selective"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"time"
)

var blobtos3 = &cobra.Command{
	Use:   "aztos3",
	Short: "command to transmit data from blobstorage to s3 bucket",
	Long:  "si",
	Run: func(cmd *cobra.Command, args []string) {
		pathfile := "/Users/r4st4m4n/Downloads/2022 Abril - Evidencias TEC.csv"
		models, err := selective.ReadCSV(pathfile)
		if err != nil {
			log.Fatalln(err)
		}
		for _, data := range models {
			start := time.Now()
			parsedurl, err := selective.ParsingUrl(data.Destino)
			if err != nil {
				log.Fatalln(err)
			}
			result := selective.BlobtoS3(parsedurl[0], data.Url)
			log.Println(fmt.Sprintf("Item Uploaded %v Time Elapsed: %v", result.Location, time.Since(start)))
		}

	},
}

func init() {
	rootCmd.AddCommand(blobtos3)
}
