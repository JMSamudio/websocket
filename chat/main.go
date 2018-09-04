package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"websocket/trace"
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
	t.templ.Execute(w, r)
}

func main() {

	var addr = flag.String("addr", ":9000", "the addr of the app")
	flag.Parse()

	r := newRoom()
	r.tracer = trace.New(os.Stdout)

	http.Handle("/", &templateHandler{filename: "chat.html"})
	http.Handle("/room", r)

	//get the room going
	go r.run()

	log.Println("Starting web server on", *addr)

	//start the web server
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenerAndServe: ", err)
	}
}
