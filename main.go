package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"
)

//templ representa un simple template
type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

//ServeHTTP captura los req HTTP
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, nil)
}

func main() {

	r := newRoom()

	http.Handle("/", &templateHandler{filename: "chat.html"})
	http.Handle("/room", r)

	//get the room going
	go r.run()

	//start the web server
	if err := http.ListenAndServe(":9000", nil); err != nil {
		log.Fatal("ListenerAndServe: ", err)
	}
}
