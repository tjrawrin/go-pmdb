package render

import (
	"net/http"
	"text/template"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseFiles("web/index.html"))
}

// HTML renders a simple HTML response and sets the content type and status.
func HTML(w http.ResponseWriter, status int, template string, v interface{}) error {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(status)

	err := tpl.ExecuteTemplate(w, template, v)
	if err != nil {
		return err
	}

	return nil
}
