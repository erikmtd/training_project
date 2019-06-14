package web

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"sync"

	"fmt"

	"github.com/training_project/core/cache"
	"github.com/training_project/core/mq"
	"github.com/training_project/handler"
	"github.com/training_project/service/user"
)

type web struct {
}

func (web *web) Start() {
	fmt.Println("Web Running")

	//http.HandleFunc("/", index)
	http.HandleFunc("/list", list)

	server := http.Server{
		Addr: ":9793",
	}
	server.ListenAndServe()
}

var lock = &sync.Mutex{}
var instance handler.Handler

func New() handler.Handler {
	lock.Lock()
	defer lock.Unlock()

	if instance == nil {
		instance = &web{}
	}
	return instance
}

func list(w http.ResponseWriter, r *http.Request) {
	q := mq.New()
	qData := struct{}{}
	q.Publish("devel-go.tkpd:4150", "tech_cur_nsq_0619_erik", qData)

	tmpl, err := findTemplate("index.html")
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("Params => ", r.URL.Query())

	var name string = ""
	if pName := r.URL.Query()["name"]; pName != nil && len(pName) > 0 {
		name = pName[0]
	}

	data := &struct {
		Vc   int
		User []user.User
	}{
		User: user.New().List(name),
	}

	visitor := cache.New().Get("training_project_0619_erik")
	fmt.Println(visitor)
	if visitor != "" {
		count, err := strconv.Atoi(visitor)
		if err == nil {
			data.Vc = count + 1
		} else {
			log.Fatal(err)
		}
	} else {
		data.Vc = 1
	}

	tmpl.Execute(w, data)
	return
}

func findTemplate(page string) (tmpl *template.Template, err error) {
	tmpl, err = template.ParseFiles(fmt.Sprintf("../../handler/web/template/%s", page))
	return
}
