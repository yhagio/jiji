package views

import (
	"bytes"
	"errors"
	"html/template"
	"io"
	"jiji/middlewares"
	"log"
	"net/http"
	"net/url"
	"path/filepath"

	"github.com/gorilla/csrf"
)

var (
	LayoutDir   string = "views/layouts/"
	TemplateDir string = "views/"
	TemplateExt string = ".gohtml"
)

const (
	AlertLvlError   = "danger"
	AlertLvlWarning = "warning"
	AlertLvlInfo    = "info"
	AlertLvlSuccess = "success"
	AlertMsgGeneric = "Something went wrong. Please try again. Contact us if the problem persists."
)

type View struct {
	Template *template.Template
	Layout   string
}

type Alert struct {
	Level   string
	Message string
}

type PublicError interface {
	error
	Public() string
}

func (d *Data) SetAlert(err error) {
	var msg string

	publicErr, ok := err.(PublicError)
	if ok {
		msg = publicErr.Public()
	} else {
		log.Println(err)
		msg = AlertMsgGeneric
	}

	d.Alert = &Alert{
		Level:   AlertLvlError,
		Message: msg,
	}
}

func (d *Data) AlertError(errMsg string) {
	d.Alert = &Alert{
		Level:   AlertLvlError,
		Message: errMsg,
	}
}

func NewView(layout string, files ...string) *View {
	addTemplatePath(files)
	addTemplateExt(files)
	files = append(files, layoutFiles()...)

	// We are changing how we create our templates, calling
	// New("") to give us a template that we can add a function to
	// before finally passing in files to parse as part of the template.
	t, err := template.New("").Funcs(template.FuncMap{
		"csrfField": func() (template.HTML, error) {
			// If this is called without being replace with a proper implementation
			// returning an error as the second argument will cause our template
			// package to return an error when executed.
			return "", errors.New("csrfField is not implemented")
		},
		"pathEscape": func(s string) string {
			return url.PathEscape(s)
		},
		// Once we have our template with a function we'll pass in files to parse.
	}).ParseFiles(files...)
	if err != nil {
		panic(err)
	}

	return &View{
		Template: t,
		Layout:   layout,
	}
}

func (v *View) Render(w http.ResponseWriter, r *http.Request, data interface{}) {
	w.Header().Set("Content-Type", "text/html")

	var vd Data
	switch d := data.(type) {
	case Data:
		vd = d
	default:
		// Wrap it for unknown data
		vd = Data{
			Yield: data,
		}
	}

	vd.User = middlewares.LookUpUserFromContext(r.Context())

	// Using a buffer here because writing any data to ResponseWriter will
	// result in 200 statusCode and we can’t undo that write.
	// By writing to a buffer first we can confirm that the whole template
	// executes before starting writing any data to ResponseWriter.
	var buff bytes.Buffer
	// We need to create the csrfField using the current http request.
	csrfField := csrf.TemplateField(r)
	tpl := v.Template.Funcs(template.FuncMap{
		// We can also change the return type of our function, since we no longer
		// need to worry about errors.
		"csrfField": func() template.HTML {
			// We can then create this closure that returns the csrfField for
			// any templates that need access to it.
			return csrfField
		},
	})

	err := tpl.ExecuteTemplate(&buff, v.Layout, vd)
	if err != nil {
		http.Error(w, AlertMsgGeneric, http.StatusInternalServerError)
		return
	}
	io.Copy(w, &buff)
}

// http.Handler needs ServeHTTP method
// https://golang.org/pkg/net/http/#Handler
func (v *View) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v.Render(w, r, nil)
}

func layoutFiles() []string {
	files, err := filepath.Glob(LayoutDir + "*" + TemplateExt)
	if err != nil {
		panic(err)
	}
	return files
}

func addTemplatePath(files []string) {
	for i, f := range files {
		files[i] = TemplateDir + f
	}
}

func addTemplateExt(files []string) {
	for i, f := range files {
		files[i] = f + TemplateExt
	}
}
