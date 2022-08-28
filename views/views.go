package views

import (
	"html/template"
	"path/filepath"
)

type View struct {
	Template *template.Template
}

// NewView is a constructor function used to create new view
// executable template parsed with layouts
func NewView(files ...string) *View {
	// get all files in the layout directory
	layoutFiles, err := filepath.Glob("./views/layouts/*.gohtml")
	if err != nil {
		panic(err)
	}

	// apeend fileName wihh layout files
	files = append(files, layoutFiles...)

	// parse template file with layout files
	t, err := template.ParseFiles(files...)
	if err != nil {
		panic(err)
	}

	// return the view
	return &View{
		Template: t,
	}
}
