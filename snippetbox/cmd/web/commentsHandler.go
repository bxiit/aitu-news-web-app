package main

import (
	"alexedwards.net/snippetbox/pkg/forms"
	"fmt"
	"net/http"
	"strconv"
)

func (app *application) createComment(w http.ResponseWriter, r *http.Request) {
	userId := app.session.Get(r, "authenticatedUserID")
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	newsId, _ := strconv.Atoi(form.Get("id"))
	_, _ = app.comment.InsertComment(userId.(int), newsId, form.Get("text"))

	app.session.Put(r, "commentFlash", "Comment added successfully")
	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", newsId), http.StatusSeeOther)
}

func (app *application) deleteComment(w http.ResponseWriter, r *http.Request) {
	//userId := app.session.Get(r, "authenticatedUserID")
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	commentId, _ := strconv.Atoi(form.Get("commentId"))
	newsId, _ := strconv.Atoi(form.Get("newsId"))
	_ = app.comment.DeleteComment(commentId)

	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", newsId), http.StatusSeeOther)
}
