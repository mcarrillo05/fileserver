package fileserver

import (
	"html/template"
	"math"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

type typeFile int

const (
	dirType typeFile = iota
	fileType
)

var (
	sizes = [...]string{"Bytes", "KB", "MB", "GB", "TB"}
)

type itemsTemplate struct {
	Root  item
	Items []item
}

type item struct {
	Name       string
	Size       int64
	SizeString string
	Date       string
	Type       typeFile
}

type fileServer struct {
	root       string
	countFiles bool
	template   *template.Template
}

func (f fileServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	dir := req.FormValue("path")
	searchPath := path.Join(f.root, dir)
	if !strings.HasPrefix(searchPath, f.root) {
		return
	}
	items := getItems(searchPath, f.countFiles)
	if len(items) > 0 {
		if items[0].Type == fileType {
			http.ServeFile(w, req, searchPath)
		} else {
			tmpl := template.Must(f.template, nil)
			err := tmpl.Execute(w, itemsTemplate{
				Root:  items[0],
				Items: items[1:],
			})
			if err != nil {
				w.Write([]byte(err.Error()))
			}
		}
	}
}

func getItems(root string, countFiles bool) (items []item) {
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		i := item{
			Name: info.Name(),
		}
		if !info.IsDir() {
			i.Type = fileType
			i.Size = info.Size()
			i.Date = info.ModTime().Format("2006-01-02 03:04PM")
		} else {
			i.Name += "/"
		}
		items = append(items, i)
		if info.IsDir() && path != root {
			return filepath.SkipDir
		}
		return nil
	})
	for i := range items {
		if countFiles {
			if i == 0 {
				items[i].countFiles(root)
			} else {
				items[i].countFiles(filepath.Join(root, items[i].Name))
			}
		}
		items[i].formatSize()
	}
	return items
}

func (i *item) countFiles(root string) {
	if i.Type == fileType {
		return
	}
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			i.Size++
		}
		return nil
	})
}

//FileServer creates a new http Handler using a template to list files.
func FileServer(root string, countFiles bool, tmpl string) http.Handler {
	t, err := template.ParseFiles(tmpl)
	if err != nil {
		panic(err)
	}
	return fileServer{root: root, countFiles: countFiles, template: t}
}

func (i *item) formatSize() {
	if i.Type == dirType {
		i.SizeString = strconv.FormatInt(i.Size, 10) + " archivo(s)"
	} else {
		if i.Size == 0 {
			i.SizeString = "0 Bytes"
			return
		}
		idx := int(math.Log(float64(i.Size)) / math.Log(1024))
		if idx == 0 {
			i.SizeString = strconv.FormatInt(i.Size, 10) + " " + sizes[idx]
		} else {
			i.SizeString = strconv.FormatFloat(float64(i.Size)/math.Pow(1024, float64(idx)), 'f', 2, 64) + " " + sizes[idx]
		}
	}
}
