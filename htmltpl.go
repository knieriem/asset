package asset

import (
	"fmt"
	"html/template"
)

func ParseHTMLTemplateFiles(t *template.Template, filenames ...string) (*template.Template, error) {
	if t == nil {
		t = new(template.Template)
	}
	fmt.Println(ns)
	return t.ParseFS(ns, filenames...)
}
