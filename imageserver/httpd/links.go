package httpd

import (
	"fmt"
	"io"
)

func writeLinks(writer io.Writer) {
	fmt.Fprintln(writer, `<a href="/">Status</a>`)
	fmt.Fprintln(writer, `<a href="listImages">Images</a>`)
	fmt.Fprintln(writer, `<p>`)
}
