package config

import (
	//"github.com/joho/godotenv"

	"os"
)

type ConfI interface {
	GetPort() string
}

type ConfElastic struct {
	port string
}

func (c *ConfElastic) GetPort() string {
	return c.port
}

func GetPort() string {
	return os.Getenv("PORT")
}

func New() error {
	openEnv()
	return nil
}

func openEnv(){
	//err := godotenv.Load("C:\\Users\\User\\go\\src\\github.com\\keepcalmist\\workwithElastic\\.env")
	//if err != nil {
	//	log.Println(err)
	//}
}



