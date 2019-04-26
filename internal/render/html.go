package render

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var tpl *template.Template

func init() {
	tpl = template.Must(findAndParseTemplates("internal/templates/", nil))
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

// findAndParseTemplates finds all template files in the root and sub
// directories and parses them. Making them accessible via "dir/name.html".
// https://stackoverflow.com/a/50581032
func findAndParseTemplates(rootDir string, funcMap template.FuncMap) (*template.Template, error) {
	cleanRoot := filepath.Clean(rootDir)
	pfx := len(cleanRoot) + 1
	root := template.New("")

	err := filepath.Walk(cleanRoot, func(path string, info os.FileInfo, e1 error) error {
		if !info.IsDir() && strings.HasSuffix(path, ".html") {
			if e1 != nil {
				return e1
			}

			b, e2 := ioutil.ReadFile(path)
			if e2 != nil {
				return e2
			}

			name := path[pfx:]
			t := root.New(name).Funcs(funcMap)
			t, e2 = t.Parse(string(b))
			if e2 != nil {
				return e2
			}
		}

		return nil
	})

	return root, err
}
