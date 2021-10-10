package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func init() {
	os.Setenv("Version", "1.0.0")
}

func helloworld(w http.ResponseWriter, r *http.Request) {
	status_code := 200
	if r.URL.Path == "/" {
		for k, v := range r.Header {
			//fmt.Printf("%s: %s\n", k, v)
			w.Header().Add(k, v[0])
		}
		w.Header().Set("Version", os.Getenv("Version"))
		fmt.Println(w.Header())
		w.Write([]byte("Hello world"))
	} else if r.URL.Path == "/healthz" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("200"))
	} else {
		status_code = 404
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Error 404, Page not found"))
	}
	log.Printf("Client IP address: %s, Status code: %d", r.RemoteAddr, status_code)
}

//func health(w http.ResponseWriter, r *http.Request) {
//	w.WriteHeader(http.StatusOK)
//	w.Write([]byte("200"))
//	log.Printf("Client IP address: %s, Status code: %d", r.RemoteAddr, http.StatusOK)
//}

func main() {
	http.HandleFunc("/", helloworld)
//	http.HandleFunc("/healthz", health)
	log.Println("Starting HTTP server on port 80...")
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal(err, "\nHttp server failed to start...")
	}
}
