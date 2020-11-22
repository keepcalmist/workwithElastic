package Server

import "net/http"

func statusHandler() func (w http.ResponseWriter, r *http.Request){
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK!"))
	}
}
