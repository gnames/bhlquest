package web

import (
	"embed"
	"fmt"
	"html/template"
	"io"
	"path"
	"strings"

	"github.com/labstack/echo/v4"
)

//go:embed templates
var tmpls embed.FS

// echoTempl implements echo.Renderer interface.
type echoTempl struct {
	templates *template.Template
}

// Render implements echo.Renderer interface.
func (t *echoTempl) Render(
	w io.Writer,
	name string,
	data interface{},
	c echo.Context,
) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func NewTemplate() (*echoTempl, error) {
	t, err := parseFiles()
	if err != nil {
		return nil, fmt.Errorf("cannot parse templates: %w", err)
	}
	return &echoTempl{t}, nil
}

func parseFiles() (*template.Template, error) {
	var err error
	var t *template.Template

	var filenames []string
	dir := "templates"
	entries, _ := tmpls.ReadDir(dir)
	for i := range entries {
		if entries[i].Type().IsRegular() {
			filenames = append(
				filenames,
				fmt.Sprintf("%s/%s", dir, entries[i].Name()),
			)
		}
	}

	for _, filename := range filenames {
		name := path.Base(filename)
		var tmpl *template.Template
		if t == nil {
			t = template.New(name)
		}
		if name == t.Name() {
			tmpl = t
		} else {
			tmpl = t.New(name)
		}
		addFuncs(tmpl)
		_, err = tmpl.ParseFS(tmpls, filename)
		if err != nil {
			return nil, err
		}
	}
	return t, nil
}

func addFuncs(tmpl *template.Template) {
	tmpl.Funcs(template.FuncMap{
		"isEven": func(i int) bool {
			return i%2 == 0
		},
		"formatFloat": func(f float64) string {
			return fmt.Sprintf("%03f", f)
		},
		"upper": func(s string) string {
			return strings.ToUpper(s)
		},
		"formatText": func(s string) template.HTML {
			res := strings.Replace(template.HTMLEscapeString(s), "\n", "<br>", -1)
			res = strings.Replace(res, "&lt;em&gt;", "<em>", -1)
			res = strings.Replace(res, "&lt;/em&gt;", "</em>", -1)

			return template.HTML(res)
		},
	})
}
