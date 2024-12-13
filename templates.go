package main

import (
	"html/template"
)

var tmpl = make(map[string]*template.Template)

func init() {
	m := template.Must
	p := template.ParseFiles

	// Parsing templates and adding them to the map
	tmpl["index"] = m(p("templates/index.gohtml", "templates/layout.gohtml"))
	tmpl["new_event"] = m(p("templates/new_event.gohtml", "templates/layout.gohtml"))
	tmpl["event_detail"] = m(p("templates/event_detail.gohtml", "templates/layout.gohtml"))
	tmpl["login"] = m(p("templates/login.gohtml", "templates/layout.gohtml"))
	tmpl["register"] = m(p("templates/register.gohtml", "templates/layout.gohtml"))
	tmpl["about"] = m(p("templates/about.gohtml", "templates/layout.gohtml"))
}
