package main

import (
	"flag"
	"log"
)

func main() {
	path := flag.String("logpath", "logdata.json", "path to logfile")
	flag.Parse()
	logStorage := newStorage()
	logService := newLogService(logStorage)
	err := logService.getData(*path)
	if err != nil {
		log.Println(err)
		return
	}
	webService := newWebService(logService)
	webService.run()
}
