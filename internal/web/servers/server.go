package servers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/sessions"

	"GameJamPlatform/internal/log"
	"GameJamPlatform/internal/services/gamejams"
	"GameJamPlatform/internal/services/sessionprovider"
	"GameJamPlatform/internal/services/users"
	"GameJamPlatform/internal/web/templatemanager"
)

type server struct {
	sessionProvider sessionprovider.SessionProvider
	cookieStore     *sessions.CookieStore

	service *gamejams.Service
	users   users.Users
	tm      templatemanager.TemplateManager

	cfg *Config
}

func NewServer(service *gamejams.Service, tm templatemanager.TemplateManager, users users.Users, sessionProvider sessionprovider.SessionProvider, cfg *Config) Server {
	return &server{
		sessionProvider: sessionProvider,
		cookieStore:     sessions.NewCookieStore([]byte(cfg.SessionKey)),

		service: service,
		users:   users,
		tm:      tm,

		cfg: cfg,
	}
}

func (s *server) Run() error {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(log.LoggerMiddleware())
	r.Use(middleware.Recoverer)
	r.Use(s.authMiddleware)

	fs := http.FileServer(http.Dir(s.cfg.StaticDir))

	r.Route("/", func(r chi.Router) {
		r.Get("/", s.indexHandlerGet())

		// Jams
		r.Get("/jams", s.jamsListHandlerGet())
		r.Get("/jam/new", s.jamNewHandlerGet())
		r.Post("/jam/new", s.jamNewHandlerPost())

		r.Get("/jams/{jamURL}", s.jamOverviewHandlerGet())
		r.Get("/jams/{jamURL}/entries", s.jamEntriesHandlerGet())
		r.Get("/jams/{jamID}/delete", s.jamDeleteHandlerGet())

		r.Get("/jams/{jamID}/edit", s.jamEditHandlerGet())
		r.Post("/jams/{jamID}/edit", s.jamEditHandlerPost())

		// Games
		r.Get("/jams/{jamURL}/game/new", s.gameNewHandlerGet())
		r.Post("/jams/{jamURL}/game/new", s.gameNewHandlerPost())

		r.Get("/jams/{jamURL}/games/{gameURL}/edit", s.gameEditHandlerGet())
		r.Post("/jams/{jamURL}/games/{gameURL}/edit", s.gameEditHandlerPost())

		r.Get("/jams/{jamURL}/games/{gameURL}", s.gameOverviewHandlerGet())
		r.Get("/jams/{jamURL}/games/{gameURL}/ban", s.gameBanHandlerGet())

		// Users
		r.Get("/user/new", s.userNewHandlerGet())
		r.Post("/user/new", s.userNewHandlerPost())

		r.Get("/user/login", s.loginHandlerGet())
		r.Post("/user/login", s.loginHandlerPost())
		r.Get("/user/logout", s.logoutHandlerGet())

		//r.Get("/users/{username}", s.userProfileHandler())
		//r.Get("/users/{username}/edit", s.userEditHandler())
		//r.Post("/users/{username}/edit", s.userUpdateHandler())

		//r.Get("/jams/{jamURL}/results", s.jamResultsHandler())

		r.Handle("/static/*", http.StripPrefix("/static/", fs))
	})

	return http.ListenAndServe(s.cfg.HostURI, r)
}
