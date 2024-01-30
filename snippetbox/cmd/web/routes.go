package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createNews)
	mux.HandleFunc("/snippet/showCreate", app.showCreate)
	mux.HandleFunc("/create", app.createNews)
	mux.HandleFunc("/snippet/delete", app.delete)
	mux.HandleFunc("/snippet/update", app.updateNews)
	mux.HandleFunc("/snippet/showUpdate", app.showUpdate)
	mux.HandleFunc("/contacts", app.contacts)
	mux.HandleFunc("/students", app.students)
	mux.HandleFunc("/staff", app.staff)
	mux.HandleFunc("/applicants", app.applicants)
	mux.HandleFunc("/researchers", app.researchers)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return mux
}
