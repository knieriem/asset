package asset

import (
	"html/template"
)

func ParseHTMLTemplateFiles(t *template.Template, filenames ...string) (*template.Template, error) {
	if t == nil {
		t = new(template.Template)
	}
	return t.ParseFS(ns, filenames...)
}
