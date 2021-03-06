package herd

import (
	"bufio"
	"fmt"
	"github.com/Symantec/Dominator/lib/html"
	"net/http"
)

func statusHandler(w http.ResponseWriter, req *http.Request) {
	writer := bufio.NewWriter(w)
	defer writer.Flush()
	fmt.Fprintln(writer, "<title>Dominator status page</title>")
	fmt.Fprintln(writer, "<body>")
	fmt.Fprintln(writer, "<center>")
	fmt.Fprintln(writer, "<h1><b>Dominator</b> status page</h1>")
	fmt.Fprintln(writer, "</center>")
	html.WriteHeader(writer)
	fmt.Fprintln(writer, "<h3>")
	httpdHerd.writeHtml(writer)
	for _, htmlWriter := range httpdHerd.htmlWriters {
		htmlWriter.WriteHtml(writer)
	}
	fmt.Fprintln(writer, "</h3>")
	fmt.Fprintln(writer, "<hr>")
	html.WriteFooter(writer)
	fmt.Fprintln(writer, "</body>")
}
