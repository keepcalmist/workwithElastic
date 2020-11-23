package main

import (
	"fmt"
	"github.com/keepcalmist/workwithElastic/pkg/Server"
	"github.com/keepcalmist/workwithElastic/pkg/config"
	"github.com/keepcalmist/workwithElastic/pkg/storage"
	"log"
)

func init(){
	_ = config.New()
}

func main(){
	go Server.Run()
	cl, err := storage.InitElastic()
	if err != nil {
		log.Println(err)
	}
	lol, err := cl.Ping()
	fmt.Println(lol)

}


