package views

import (
	"html/template"
)

type View struct {
	Template *template.Template
}

// create new view executable template parsed with layouts
func NewView(files ...string) *View {
	// apeend fileName wihh layout files
	files = append(files, "views/layouts/footer.gohtml")

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
