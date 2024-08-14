package mails

import (
	"bytes"
	"html/template"
)

type Template interface {
	Template() string
}

func Parse(temp Template) string {
	body, err := template.New("api").Parse(temp.Template())
	if err != nil {
		return ""
	}
	buf := new(bytes.Buffer)
	if err = body.Execute(buf, temp); err != nil {
		return ""
	}
	return buf.String()
}
