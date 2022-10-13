package spikeweb

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
	engine := new(Engine)
	log.Fatal(http.ListenAndServe(":9999", engine))
}
