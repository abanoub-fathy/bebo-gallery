package views

import (
	"html/template"
	"net/http"
	"path/filepath"
)

// View type is a struct contain the template of the view
// and also the layout template
type View struct {
	Template *template.Template
	Layout   string
}

var (
	LayoutDir         = "views/layouts/"
	TemplateDir       = "views/"
	TemplateExtension = ".gohtml"
)

// NewView is a constructor function used to create new view
// executable template parsed with layouts
func NewView(layout string, files ...string) *View {
	// adding template path
	addTemplatePath(files)

	// adding template Extension
	addTemplateExtension(files)

	// apend fileName wihh layout files
	files = append(files, getLayoutFiles()...)

	// parse template file with layout files
	t, err := template.ParseFiles(files...)
	if err != nil {
		panic(err)
	}

	// return the view
	return &View{
		Template: t,
		Layout:   layout,
	}
}

// GetLayoutFiles is a func used to return all layout files
func getLayoutFiles() []string {
	// get all files in the layout directory
	layoutFiles, err := filepath.Glob(LayoutDir + "*" + TemplateExtension)
	if err != nil {
		panic(err)
	}

	return layoutFiles
}

// ServeHttp is used to implement the Handler type
// now the *view type can be used as a Handler type
func (view *View) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := view.Render(w, nil); err != nil {
		panic(err)
	}
}

// Render is used to render a view based on the predefined layout
func (view *View) Render(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "text/html")
	return view.Template.ExecuteTemplate(w, view.Layout, data)
}

// addTemplatePath takes in a slice of strings
// representing file paths for templates, and it prepends
// the TemplateDir directory to each string in the slice
//
// Eg the input {"home"} would result in the output
// {"views/home"} if TemplateDir == "views/"
func addTemplatePath(files []string) {
	for i, file := range files {
		files[i] = TemplateDir + file
	}
}

// addTemplateExtension takes in a slice of strings
// representing file paths for templates and it appends
// the TemplateExt extension to each string in the slice
//
// Eg the input {"home"} would result in the output
// {"home.gohtml"} if TemplateExt == ".gohtml"
func addTemplateExtension(files []string) {
	for i, file := range files {
		files[i] = file + TemplateExtension
	}
}
