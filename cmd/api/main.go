package main

import (
	"fmt"
	"github.com/keepcalmist/workwithElastic/pkg/config"
)

func main() {
	cfg := config.New()
	fmt.Println(cfg.GetString("hui"))
}
