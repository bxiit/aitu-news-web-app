package main

import (
	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
	"net/http"
)

// secureHeaders → servemux → application handler → servemux → secureHeaders
// logRequest ↔ secureHeaders ↔ servemux ↔ application handler
func (app *application) routes() http.Handler {

	// Содержит стандартные middleware, которые будут применяться к каждому запросу,
	// ( 1 ) включая обработку паник, ( 2 ) логгирование и  ( 3 )настройку безопасных заголовков.
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	// dynamicMiddleware содержит middleware, специфичные для динамических маршрутов вашего приложения,
	// например, обработку сеансов.
	dynamicMiddleware := alice.New(app.session.Enable)

	mux := pat.New()

	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))

	mux.Get("/snippet", dynamicMiddleware.ThenFunc(app.showNews))

	// Add the requireAuthentication middleware to the chain.
	mux.Get("/snippet/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.showCreate))

	// Add the requireAuthentication middleware to the chain.
	mux.Post("/snippet/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createNews))
	mux.Post("/snippet/delete", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.delete))
	mux.Post("/snippet/update", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.updateNews))
	mux.Post("/snippet/showUpdate", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.showUpdate))

	mux.Get("/contacts", dynamicMiddleware.ThenFunc(app.contacts))
	mux.Get("/students", dynamicMiddleware.ThenFunc(app.students))
	mux.Get("/staff", dynamicMiddleware.ThenFunc(app.staff))
	mux.Get("/applicants", dynamicMiddleware.ThenFunc(app.applicants))
	mux.Get("/researchers", dynamicMiddleware.ThenFunc(app.researchers))
	mux.Get("/user/manage", dynamicMiddleware.ThenFunc(app.manageUsers))
	mux.Post("/user/changeRole", dynamicMiddleware.ThenFunc(app.changeRole))
	mux.Get("/snippet/:id", dynamicMiddleware.ThenFunc(app.showNews))

	// authentication
	mux.Get("/user/signup", dynamicMiddleware.ThenFunc(app.signupUserForm))
	mux.Post("/user/signup", dynamicMiddleware.ThenFunc(app.signupUser))
	mux.Get("/user/login", dynamicMiddleware.ThenFunc(app.loginUserForm))
	mux.Post("/user/login", dynamicMiddleware.ThenFunc(app.loginUser))
	mux.Get("/user/:id", dynamicMiddleware.ThenFunc(app.showUser))

	mux.Post("/user/logout", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.logoutUser))

	// Leave the static files route unchanged.
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(mux)
}
