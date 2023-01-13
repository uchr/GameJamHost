package servers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"GameJamPlatform/internal/log"
	"GameJamPlatform/internal/services"
	"GameJamPlatform/internal/templates"
)

type server struct {
	service *services.Service
	tmpl    templates.Templates
	cfg     *Config
}

func NewServer(service *services.Service, tmpl templates.Templates, cfg *Config) Server {
	return &server{
		service: service,
		tmpl:    tmpl,
		cfg:     cfg,
	}
}

func (s *server) Run() error {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(log.LoggerMiddleware())
	r.Use(middleware.Recoverer)

	fs := http.FileServer(http.Dir(s.cfg.StaticDir))

	r.Route("/", func(r chi.Router) {
		r.Get("/jams", s.jamsListHandler())
		r.Get("/jam/new", s.jamNewHandler())
		r.Post("/jam/new", s.jamCreateHandler())

		r.Get("/jams/{jamURL}", s.jamOverviewHandler())
		r.Get("/jams/{jamURL}/entries", s.jamEntriesHandler())

		r.Get("/jams/{jamID}/edit", s.jamEditHandler())
		r.Get("/jams/{jamID}/delete", s.jamDeleteHandler())
		r.Post("/jams/{jamID}/edit", s.jamUpdateHandler())

		r.Get("/jams/{jamURL}/game/new", s.gameNewHandler())
		r.Post("/jams/{jamURL}/game/new", s.gameCreateHandler())

		r.Get("/jams/{jamURL}/games/{gameURL}", s.gameOverviewHandler())
		r.Get("/jams/{jamURL}/games/{gameURL}/edit", s.gameEditHandler())
		r.Get("/jams/{jamURL}/games/{gameURL}/ban", s.gameBanHandler())
		r.Post("/jams/{jamURL}/games/{gameURL}/edit", s.gameUpdateHandler())

		//r.Get("/jams/{jamURL}/results", s.jamResultsHandler())

		r.Handle("/static/*", http.StripPrefix("/static/", fs))
	})

	return http.ListenAndServe(s.cfg.HostURI, r)
}
