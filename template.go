package asset

import (
	"path"
	"text/template"
)

func ParseTemplateFiles(t *template.Template, filenames ...string) (tnew *template.Template, err error) {
	for _, name := range filenames {
		tpl, err1 := FileString(name)
		if err1 != nil {
			err = err1
			break
		}
		t, err = t.New(path.Base(name)).Parse(tpl)
		if err != nil {
			break
		}
	}
	tnew = t
	return
}
