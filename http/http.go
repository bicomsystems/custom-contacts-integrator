package http

import (
	"csm/config"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
)

func InitHttp(db Database) {
	h := &Http{db: db}

	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	r.Get("/token", handleGenerateJWT)

	r.Route("/contacts", func(r chi.Router) {
		r.Use(authMiddleware)
		r.Get("/", h.getContactsHandler)           // GET /contacts
		r.Get("/delta", h.getDeltaContactsHandler) // GET /contacts/delta
	})

	go http.ListenAndServe(fmt.Sprintf(":%d", config.Conf.HttpServer.Port), r)
}
