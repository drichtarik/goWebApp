package handlers

import (
	"github.com/gorilla/mux"
	"html/template"
	"io/ioutil"
	"net/http"
	"regexp"
)

var templates = template.Must(template.ParseFiles("template/edit.html", "template/view.html", "template/create.html", "template/static/head.html", "template/static/navbar.html", "template/static/footer.html"))
var validPath = regexp.MustCompile("^/(edit|savePage|view|create)/([a-zA-Z0-9]+)$")

type Page struct {
	Title string `json:"title,omitempty"`
	Body  []byte `json:"body"`
}

var Pages []Page

func (p *Page) savePage() error {
	filename := p.Title + ".txt"
	Pages = append(Pages, *p)
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

func CacheDefaultPages() []Page {
	names := []string{"FrontPage", "TestPage"}
	for i := 0; i < 2; i++ {
		p := &Page{}
		p, _ = loadPage(names[i])
		Pages = append(Pages, *p)
	}
	return Pages
}

func pageRenderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func pageMakeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			fn(w, r, "/")
			return
		}
		fn(w, r, m[2])
	}
}

func pageRootHandler(w http.ResponseWriter, r *http.Request, title string) {
	myTitle := "FrontPage"
	http.Redirect(w, r, "/view/"+myTitle, http.StatusFound)
}

func pageViewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	pageRenderTemplate(w, "view", p)
}

func pageEditHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	pageRenderTemplate(w, "edit", p)
}

func pageSaveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.savePage()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func pageCreateHandler(w http.ResponseWriter, r *http.Request) {
	p := &Page{}
	pageRenderTemplate(w, "create", p)
}

func pageCreateNewHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	title := r.Form["new_title"]
	body := r.Form["new_body"]

	p := &Page{title[0], []byte(body[0])}
	err := p.savePage()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title[0], http.StatusFound)
}

func BootAllPageHandlers(router *mux.Router) {
	CacheDefaultPages()
	http.HandleFunc("/", pageMakeHandler(pageRootHandler))
	http.HandleFunc("/view/", pageMakeHandler(pageViewHandler))
	http.HandleFunc("/edit/", pageMakeHandler(pageEditHandler))
	http.HandleFunc("/save/", pageMakeHandler(pageSaveHandler))
	http.HandleFunc("/create/", pageCreateHandler)
	http.HandleFunc("/createNew/", pageCreateNewHandler)
}
