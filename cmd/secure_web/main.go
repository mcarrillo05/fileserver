package main

import (
	"flag"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mcarrillo05/fileserver"
)

func main() {
	dir := flag.String("d", "", "path to directory (obligatory)")
	port := flag.String("p", "8000", "webserver port (8000 if not set)")
	count := flag.Bool("count-files", false, "if set, directory will be showed with number of files")
	user := flag.String("u", "", "user name")
	pass := flag.String("pass", "", "password")
	tmpl := flag.String("t", "", "path to template")
	flag.Parse()
	if *dir == "" {
		flag.PrintDefaults()
		panic("directory flag must be set")
	}
	if *user == "" {
		flag.PrintDefaults()
		panic("user flag must be set")
	}
	if *pass == "" {
		flag.PrintDefaults()
		panic("pass flag must be set")
	}
	r := gin.Default()
	secure := r.Group("/", gin.BasicAuth(map[string]string{
		*user: *pass,
	}))
	{
		secure.GET("", func(c *gin.Context) {
			c.Redirect(http.StatusTemporaryRedirect, "/files.html")
		})
		secure.GET("files.html", gin.WrapH(fileserver.FileServer(*dir, *count, *tmpl)))
	}
	err := r.Run(":" + *port)
	if err != nil {
		panic(err)
	}
}
