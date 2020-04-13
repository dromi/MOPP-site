package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"

	"github.com/dromi/MOPP-site/model"
)

type MetaData struct {
	Username string
}

var templates = template.Must(template.ParseFiles(
	"tmpl/frontpage.html",
	"tmpl/calendar.html",
	"tmpl/performers.html",
	"tmpl/signin.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, p interface{}) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {

	remakeDB := flag.Bool("db", false, "recreate db")
	newPerformer := flag.Bool("nu", false, "add new performer")
	name := flag.String("name", "", "name for creating new performers")
	password := flag.String("password", "", "password for creating new performers")
	flag.Parse()

	switch {
	case *remakeDB:
		model.CreateDB()
	case *newPerformer:
		if *name == "" || *password == "" {
			log.Fatal("No username or password provided")
		}
		HashPassword(password)
		model.AddPerformer(model.Performer{Name: *name, Password: *password})
	default:
		http.HandleFunc("/calendar/", authHandler(calendarHandler))
		http.HandleFunc("/performers/", authHandler(performerHandler))
		http.HandleFunc("/", authHandler(frontpageHandler))
		http.HandleFunc("/signin", signinHandler)
		log.Fatal(http.ListenAndServe(":8080", nil))
	}
}
