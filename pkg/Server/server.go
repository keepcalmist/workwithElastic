package Server

import (
	"github.com/gorilla/mux"
	"github.com/keepcalmist/workwithElastic/pkg/config"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run(status chan int,log *log.Logger) {
	r := initRouter()
	serv := initServer(r)
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt,syscall.SIGINT, syscall.SIGTERM)

	go func(){
		if err := serv.ListenAndServe(); err != nil {
			log.Panic("Something wrong with server")
		}
	}()

	waitForStop := <- signals

	signal.Stop(signals)

	log.Println("Stop signal: ", waitForStop)

}


//Router
func initRouter() *mux.Router{
	r := mux.NewRouter()
	r.Methods(http.MethodOptions)
	r.HandleFunc("/status", statusHandler()).Methods(http.MethodGet)
	return r
}

func initServer(r *mux.Router) *http.Server {
	return &http.Server{
		Addr:              ":"+config.GetPort(),
		Handler:           r,
		TLSConfig:         nil,
		ReadTimeout:       10*time.Second,
		ReadHeaderTimeout: 0,
		WriteTimeout:      10*time.Second,
	}
}


