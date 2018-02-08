package main

import (
	"flag"
	"net/http"

	"github.com/mcarrillo05/fileserver"
)

func main() {
	dir := flag.String("d", "", "path to directory (obligatory)")
	port := flag.String("p", "8000", "webserver port (8000 if not set)")
	count := flag.Bool("count-files", false, "if set, directory will be showed with number of files")
	tmpl := flag.String("t", "", "path to template")
	flag.Parse()
	if *dir == "" {
		flag.PrintDefaults()
		panic("directory flag must be set")
	}
	http.Handle("/files.html", fileserver.FileServer(*dir, *count, *tmpl))
	err := http.ListenAndServe(":"+*port, nil)
	if err != nil {
		panic(err)
	}
}
