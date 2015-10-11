package main

import (
	"html/template"
	"net"
	"net/http"

	"github.com/zenazn/goji"
	"github.com/zenazn/goji/graceful/listener"
	"github.com/zenazn/goji/web"
)

func handler(c web.C, w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("templates/index.html"))
	data := struct {
		Title string
	}{
		"goof",
	}

	t.Execute(w, data)
}

func main() {
	rawListener, err := net.Listen("tcp", ":8081")
	if err != nil {
		panic(err)
	}

	gracefulListener := listener.Wrap(rawListener, listener.Automatic)

	goji.Get("/", handler)
	goji.DefaultMux.Handle(
		"/static/*",
		http.StripPrefix("/static/", http.FileServer(http.Dir("static"))),
	)

	goji.ServeListener(gracefulListener)
}
