package view

import (
	"fmt"
	"net/http"
	"text/template"
)

type Renderer struct {
	templates *template.Template
}

func NewRenderer(filenames ...string) *Renderer {
	return &Renderer{}
}

func (o *Renderer) ParseFiles(filenames ...string) error {
	var err error
	if o.templates, err = template.ParseFiles(filenames...); err != nil {
		return err
	}

	fmt.Println("Files parsed successfully")
	return nil
}

func (o *Renderer) RenderTemplate(w http.ResponseWriter, tmpl string) error {
	err := o.templates.ExecuteTemplate(w, tmpl+".html", nil)
	if err != nil {
		return err
	}

	return nil
}
