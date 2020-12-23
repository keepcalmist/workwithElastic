package main

import (
	"fmt"
	"github.com/keepcalmist/workwithElastic/pkg/Server"
	"github.com/keepcalmist/workwithElastic/pkg/config"
	"github.com/keepcalmist/workwithElastic/pkg/storage"
	log "github.com/sirupsen/logrus"

	"os"
)

var (
	LOG_FILE = "logs.log"
)

func init(){
	_ = config.New()
}

func main(){
	//creating new file to write logs
	file, err := os.OpenFile(LOG_FILE, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		//TODO add handle this error
		panic("Cant open or create logFile")
		return
	}
	//init logger
	logs := initLogger(file)

	status := make(chan int,1 )
	//server with chan to graceful shutdown
	go Server.Run(status,logs)


	cl, err := storage.InitElastic()
	if err != nil {
		log.Println(err)
	}


	res, err := cl.Info()
	if err != nil {
		log.Println("Error getting info")
		fmt.Println("Error getting info")
	}
	fmt.Println(res)
	fmt.Println("Exit from main with status: ",<-status)
}


func initLogger(file *os.File) *log.Logger {
	logger := log.New()
//Set output for logs
	logger.SetOutput(file)
	logger.SetFormatter(&log.TextFormatter{})
	log.RegisterExitHandler(ShutDown)
	return logger
}


func ShutDown () {
	log.Exit(0)
}