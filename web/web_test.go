package web

import (
	"fmt"
	"log"
	"net/http"
	"testing"
)

func Test_StartServer(t *testing.T) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "URL.PATH = %q\n", r.URL.Path)
	})
	log.Fatal(http.ListenAndServe(":9999", nil))
}

func Test_StartEngine(t *testing.T) {
	web := New()
	web.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello Gee")
	})
	web.Post("/login", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Login Succeed")
	})
	log.Fatal(web.Run(":9999"))
}
