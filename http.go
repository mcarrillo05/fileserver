package fileserver

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"path"
	"path/filepath"
	"strings"
)

type fileServer struct {
	root       string
	countFiles bool
	template   *template.Template
}

//FileServer creates a new http Handler using a template to list files, if template doesn't exists response will be sent using JSON format.
func FileServer(root string, countFiles bool, tmpl string) http.Handler {
	t, err := template.ParseFiles(tmpl)
	if err != nil {
		fmt.Println("serving without template, using JSON response")
	}
	return fileServer{root: root, countFiles: countFiles, template: t}
}

func (f fileServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	searchPath := path.Join(f.root, req.FormValue("path"))
	if !strings.HasPrefix(searchPath, f.root) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	items := GetItems(searchPath, f.countFiles)
	if len(items) > 0 {
		if items[0].Type == FileType {
			w.Header().Add("Content-Disposition", "attachment;filename="+filepath.Base(searchPath))
			http.ServeFile(w, req, searchPath)
		} else {
			if f.template != nil {
				//Template response
				f.serveTemplate(w, items)
			} else {
				//JSON response
				f.serverJSON(w, items)
			}
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func (f fileServer) serveTemplate(w http.ResponseWriter, items []Item) {
	tmpl := template.Must(f.template, nil)
	err := tmpl.Execute(w, itemsResponse{
		Root:  items[0],
		Items: items[1:],
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
}

func (f fileServer) serverJSON(w http.ResponseWriter, items []Item) {
	resp, err := json.Marshal(itemsResponse{
		Root:  items[0],
		Items: items[1:],
	})
	if err != nil {
		w.Write([]byte(err.Error()))
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(resp)
	}
}
