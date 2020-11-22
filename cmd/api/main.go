package main

import (

	"github.com/keepcalmist/workwithElastic/pkg/config"
	"github.com/keepcalmist/workwithElastic/pkg/Server"
)

func main(){
	_ = config.New()
	Server.Run()
}


