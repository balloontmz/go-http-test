package main

import (
	"regexp"
	"fmt"
	"log"			
	"errors"
	"net/http"
	"io/ioutil"
	"html/template"
)

// 缓存模板，缓存完成之后，只保存了文件名，路径名没有保留下来
var templates = template.Must(template.ParseFiles("template/edit.html", "template/view.html"))			

var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

type Page struct {
	Title string
	Body []byte
}

// 结构体的所属方法
func (p *Page) save() error  {						
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {			
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return "", errors.New("Invalid Page Title")
	}
	return m[2], nil 
}

/** 渲染模板的函数 */
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page)  {			
	err := templates.ExecuteTemplate(w, tmpl + ".html", p)				// 此时 模板已经被缓存过了
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handler(w http.ResponseWriter, r *http.Request)  {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func viewHandler(w http.ResponseWriter, r *http.Request)  {
	title := r.URL.Path[len("/view/"):]
	p, err := loadPage(title)

	// 如果页面不存在,重定向到编辑页面
	if err != nil {	
		http.Redirect(w, r, "/edit/" + title, http.StatusFound)
		return
	}
	// fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request)  {
	title := r.URL.Path[len("/edit/"):]
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	// fmt 包不再需要
	// fmt.Fprintf(w, "<h1>Editing %s</h1>"+
	// "<form action=\"/save/%s\" method=\"POST\">"+
	// "<textarea name=\"body\">%s</textarea><br>"+
	// "</form>",
	// p.Title, p.Title, p.Body)
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request)  {
	title := r.URL.Path[len("/save/"):]
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/" + title, http.StatusFound)
}

func main()  {
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}