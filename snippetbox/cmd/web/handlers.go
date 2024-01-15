package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"alexedwards.net/snippetbox/pkg/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	s, err := app.snippets.Latest()
	log.Printf("home after lates")
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Use the new render helper.
	app.render(w, r, "home.page.tmpl", &templateData{
		Snippets: s,
	})
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id")) // Для извлечения поля айди из ссылки
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	s, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	// Use the new render helper.
	app.render(w, r, "show.page.tmpl", &templateData{
		Snippet: s,
	})
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	// Извлекаем данные из формы
	title := r.FormValue("title")
	content := r.FormValue("content")
	expires := "7"
	category := r.FormValue("category")

	// Pass the data to the SnippetModel.Insert() method, receiving the
	// ID of the new record back.
	id, err := app.snippets.Insert(title, content, expires, category)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Redirect the user to the relevant page for the snippet.
	http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)
}

func (app *application) showCreate(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl", &templateData{})
}

func (app *application) contacts(w http.ResponseWriter, r *http.Request) {
	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "contacts.page.tmpl", &templateData{
		Snippets: s,
	})
}

func (app *application) students(w http.ResponseWriter, r *http.Request) {
	s, err := app.snippets.LatestByCategory("Students")
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "students.page.tmpl", &templateData{
		Snippets: s,
	})
}

func (app *application) staff(w http.ResponseWriter, r *http.Request) {
	s, err := app.snippets.LatestByCategory("Staff")
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "staff.page.tmpl", &templateData{
		Snippets: s,
	})
}

func (app *application) researchers(w http.ResponseWriter, r *http.Request) {
	s, err := app.snippets.LatestByCategory("Researchers")
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "researchers.page.tmpl", &templateData{
		Snippets: s,
	})
}

func (app *application) applicants(w http.ResponseWriter, r *http.Request) {
	s, err := app.snippets.LatestByCategory("Applicants")
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "applicants.page.tmpl", &templateData{
		Snippets: s,
	})
}
