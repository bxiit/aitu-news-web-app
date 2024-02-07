package main

import (
	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
	"net/http"
)

// secureHeaders → servemux → application handler → servemux → secureHeaders
// logRequest ↔ secureHeaders ↔ servemux ↔ application handler
func (app *application) routes() http.Handler {
	// Create a middleware chain containing our 'standard' middleware
	// which will be used for every request our application receives.
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	mux := pat.New()
	mux.Get("/", http.HandlerFunc(app.home))
	mux.Get("/snippet", http.HandlerFunc(app.showSnippet))
	mux.Get("/snippet/create", http.HandlerFunc(app.showCreate))
	mux.Post("/snippet/create", http.HandlerFunc(app.createNews))
	//mux.Post("/create", http.HandlerFunc(app.createNews))
	mux.Post("/snippet/delete", http.HandlerFunc(app.delete))
	mux.Post("/snippet/update", http.HandlerFunc(app.updateNews))
	mux.Post("/snippet/showUpdate", http.HandlerFunc(app.showUpdate))
	mux.Get("/contacts", http.HandlerFunc(app.contacts))
	mux.Get("/students", http.HandlerFunc(app.students))
	mux.Get("/staff", http.HandlerFunc(app.staff))
	mux.Get("/applicants", http.HandlerFunc(app.applicants))
	mux.Get("/researchers", http.HandlerFunc(app.researchers))
	mux.Get("/snippet/:id", http.HandlerFunc(app.showSnippet))

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	// Pass the servemux as the 'next' parameter to the secureHeaders middleware.
	// Because secureHeaders is just a function, and the function returns a
	// http.Handler we don't need to do anything else.

	// Return the 'standard' middleware chain followed by the servemux.
	return standardMiddleware.Then(mux)
}
