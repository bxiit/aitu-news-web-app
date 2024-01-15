package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)
	mux.HandleFunc("/snippet/showCreate", app.showCreate)
	mux.HandleFunc("/students", app.students)
	mux.HandleFunc("/staff", app.staff)
	mux.HandleFunc("/applicants", app.applicants)
	mux.HandleFunc("/researchers", app.researchers)
	mux.HandleFunc("/create", app.createSnippet)
	mux.HandleFunc("/contacts", app.contacts)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return mux
}
