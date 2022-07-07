package helpers

import (
	"log"
	"os"
)

func Logger() (*log.Logger, *log.Logger) {
	logerr, err := os.OpenFile("error.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalln("Error al crear el archivo de logs")
	}
	defer logerr.Close()
	loginfo, err := os.OpenFile("info.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalln("Error al crear el archivo de logs")
	}
	defer loginfo.Close()

	loggerinfo := log.New(loginfo, "INFO ", log.LstdFlags)
	loggererror := log.New(logerr, "ERROR ", log.LstdFlags)

	return loggerinfo, loggererror
}
