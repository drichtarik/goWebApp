package handlers

import (
	"errors"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
	"regexp"
)

var articleTemplates = template.Must(template.ParseFiles("template/edit.html", "template/view.html", "template/create.html", "template/static/head.html", "template/static/navbar.html", "template/static/footer.html"))
var articleValidPath = regexp.MustCompile("^/(edit|save|view|create)/([a-zA-Z0-9]+)$")

type Article struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

var Articles []Article

func populateArticles() []Article {
	Articles = append(Articles, Article{"index", "This is the index page"})
	Articles = append(Articles, Article{"testPage", "This is a test page"})
	return Articles
}

func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
	m := articleValidPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return "", errors.New("invalid page title")
	}
	return m[2], nil // The title is the second subexpression.
}

func loadArticle(title string) (*Article, error) {
	var a Article
	for _, f := range Articles {
		if f.Title == title {
			return &f, nil
		}
	}
	err := errors.New("page does not exist")
	return &a, err
}

func renderTemplate(w http.ResponseWriter, tmpl string, a *Article) {
	err := articleTemplates.ExecuteTemplate(w, tmpl+".html", a)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	title := "index"
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title, err := getTitle(w, r)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	a, err := loadArticle(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", a)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	title, err := getTitle(w, r)
	if err != nil {
		http.HandleFunc("/", rootHandler)
		return
	}
	a, loadErr := loadArticle(title)
	if loadErr != nil {
		a.Title = title
	}
	renderTemplate(w, "edit", a)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	title, titleErr := getTitle(w, r)
	if titleErr != nil {
		return
	}
	body := r.FormValue("body")
	a, loadErr := loadArticle(title)
	if loadErr == nil {
		for i := 0; i < len(Articles); i++ {
			if Articles[i].Title == title {
				Articles[i].Body = body
			}
		}
	} else {
		a.Title = title
		a.Body = body
		Articles = append(Articles, *a)
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	a := &Article{}
	renderTemplate(w, "create", a)
}

func createNewHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	title := r.Form["new_title"]
	body := r.Form["new_body"]

	a := &Article{title[0], body[0]}
	Articles = append(Articles, *a)
	http.Redirect(w, r, "/view/"+title[0], http.StatusFound)
}

func BootAllArticleHandlers(router *mux.Router) {
	populateArticles()
	router.HandleFunc("/", rootHandler)
	router.HandleFunc("/view/", viewHandler)
	router.HandleFunc("/edit/", editHandler)
	router.HandleFunc("/save/", saveHandler)
	router.HandleFunc("/create/", createHandler)
	router.HandleFunc("/createNew/", createNewHandler)
}
