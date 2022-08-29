package views

import (
	"html/template"
	"path/filepath"
)

type View struct {
	Template *template.Template
	Layout   string
}

var (
	LayoutDir         = "./views/layouts/"
	TemplateExtension = ".gohtml"
)

// NewView is a constructor function used to create new view
// executable template parsed with layouts
func NewView(layout string, files ...string) *View {
	// apeend fileName wihh layout files
	files = append(files, GetLayoutFiles()...)

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
func GetLayoutFiles() []string {
	// get all files in the layout directory
	layoutFiles, err := filepath.Glob(LayoutDir + "*" + TemplateExtension)
	if err != nil {
		panic(err)
	}

	return layoutFiles
}
