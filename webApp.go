package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

var templates = template.Must(template.ParseFiles("template/edit.html", "template/view.html", "template/create.html", "template/static/head.html", "template/static/navbar.html", "template/static/footer.html"))
var validPath = regexp.MustCompile("^/(edit|save|view|create)/([a-zA-Z0-9]+)$")

type Page struct {
	Title string `json:"title,omitempty"`
	Body  []byte `json:"body"`
}

var pages []Page

func (p *Page) save() error {
	filename := p.Title + ".txt"
	pages = append(pages, *p)
	return ioutil.WriteFile("data/"+filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := "data/" + title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func createNewHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	title := r.Form["new_title"]
	body := r.Form["new_body"]

	p := &Page{title[0], []byte(body[0])}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title[0], http.StatusFound)
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			fn(w, r, "/")
			return
		}
		fn(w, r, m[2])
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	p := &Page{}
	renderTemplate(w, "create", p)
}

func rootHandler(w http.ResponseWriter, r *http.Request, title string) {
	myTitle := "FrontPage"
	http.Redirect(w, r, "/view/"+myTitle, http.StatusFound)
}

func goGetPagesEndpoint(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: goGetPagesEndpoint")
	json.NewEncoder(w).Encode(pages)
}

func main() {
	http.HandleFunc("/", makeHandler(rootHandler))
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	http.HandleFunc("/create/", createHandler)
	http.HandleFunc("/createNew/", createNewHandler)

	// rest
	http.HandleFunc("/pages/", goGetPagesEndpoint)

	log.Fatal(http.ListenAndServe(":8080", nil))
}