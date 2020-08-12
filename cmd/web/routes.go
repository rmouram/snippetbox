package main

import (
	"github.com/bmizerany/pat"
	"net/http"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {

	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	mux := pat.New()
	mux.Get("/", http.HandlerFunc(app.home))
	mux.Get("/snippet/create", http.HandlerFunc(app.createSnippetForm))
	mux.Post("/snippet/create", http.HandlerFunc(app.createSnippet))
	//para o get em /snippet/... funcionar o /snippet/:id precisa ficar abaixo das outras
	mux.Get("/snippet/:id", http.HandlerFunc(app.showSnippet))

	fileServe := http.FileServer(http.Dir("../../ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServe))

	return standardMiddleware.Then(mux)
}
