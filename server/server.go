package server

import (
	"bytes"
	"fmt"
	"github.com/piger/codesearch/index"
	"html/template"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"sync"
)

const staticPrefix = "/static/"

type flushWriter struct {
	f http.Flusher
	w io.Writer
}

func (fw *flushWriter) Write(p []byte) (n int, err error) {
	n, err = fw.w.Write(p)
	if fw.f != nil {
		fw.f.Flush()
	}
	return
}

type appContext struct {
	idx       *index.Index
	Templates map[string]*template.Template
	Lock      sync.Mutex
}

type TemplateContext map[string]interface{}

func (ac *appContext) RenderTemplate(name string, context TemplateContext, w http.ResponseWriter) {
	var buf bytes.Buffer
	t, ok := ac.Templates[name]
	if !ok {
		http.Error(w, fmt.Sprintf("Cannot find template \"%s\"", name), http.StatusInternalServerError)
		return
	}
	err := t.ExecuteTemplate(&buf, "base", context)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	buf.WriteTo(w)
}

func loadTemplates() map[string]*template.Template {
	templates := make(map[string]*template.Template)
	templateNames := []string{
		"index.html",
	}
	for _, name := range templateNames {
		t1 := filepath.Join("server/templates/", name)
		t2 := "server/templates/_base.html"
		t3 := "server/templates/_extra.html"
		t := template.Must(template.New(name).ParseFiles(t1, t2, t3))
		templates[name] = t
	}

	return templates
}

type appHandler struct {
	*appContext
	h func(*appContext, http.ResponseWriter, *http.Request) (int, error)
}

func (ah appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	status, err := ah.h(ah.appContext, w, r)
	if err != nil {
		switch status {
		case http.StatusNotFound:
			http.NotFound(w, r)
		case http.StatusInternalServerError:
			http.Error(w, http.StatusText(status), status)
		default:
			http.Error(w, http.StatusText(status), status)
		}
	}
}

func RunServer(addr string, index *index.Index) {
	context := &appContext{
		idx:       index,
		Templates: loadTemplates(),
	}

	http.Handle(staticPrefix, http.StripPrefix(staticPrefix, http.FileServer(http.Dir("server/static"))))
	http.Handle("/", appHandler{context, indexHandler})
	http.Handle("/search", appHandler{context, searchHandler})

	log.Printf("Listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
