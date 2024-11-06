package main

import (
    "html/template"
)

var tmpl = make(map[string]*template.Template)

func init() {
    m := template.Must
    p := template.ParseFiles
    tmpl["index"] = m(p("templates/index.gohtml", "templates/layout.gohtml"))
    tmpl["new_event"] = m(p("templates/new_event.gohtml", "templates/layout.gohtml"))
    tmpl["event_detail"] = m(p("templates/event_detail.gohtml", "templates/layout.gohtml"))
}
