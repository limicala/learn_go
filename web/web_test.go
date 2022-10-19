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
	web.GET("/", func(c *Context) {
		c.String(http.StatusOK, "Hello Spike")
	})
	web.POST("/login", func(c *Context) {
		c.JSON(http.StatusOK, map[string]interface{}{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})
	log.Fatal(web.Run(":9999"))
}
