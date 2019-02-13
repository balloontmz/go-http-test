package main

import (
	"fmt"
	"log"			// 应该是用于日志
	"net/http"
	"io/ioutil"
	"html/template"
)

type Page struct {
	Title string
	Body []byte
}

func (p *Page) save() error  {						// 结构体的所属方法
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

/** 渲染模板的函数 */
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page)  {			
	t, err := template.ParseFiles(tmpl + ".html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, p)
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
	if err != nil {	// 如果页面不存在,重定向到编辑页面
		http.Redirect(w, r, "/edit/" + title, http.StatusFound)
		return
	}
	// fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)
	renderTemplate(w, "template/view", p)
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
	renderTemplate(w, "template/edit", p)
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