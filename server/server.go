package server

import (
	"Go-server/db"
	"Go-server/handlers"
	"Go-server/model"
	"Go-server/server/sessions"
	"Go-server/view"
	"log"
	"net/http"
	"time"
)

type Server struct {
	db       *db.PstgDB
	renderer *view.Renderer
	server   *http.Server
	sessions *sessions.Manager

	cfg *model.Config
}

func NewServer() *Server {
	return &Server{}
}

func (o *Server) registerHandlers() {
	http.HandleFunc("/login", o.makeHandler(func(w http.ResponseWriter, r *http.Request) { handlers.Login(w, r, o.db, o.sessions, o.renderer) }))
	http.HandleFunc("/logout", o.makeHandler(func(w http.ResponseWriter, r *http.Request) { handlers.Logout(w, r, o.db, o.sessions) }))
	http.HandleFunc("/account", o.makeHandler(func(w http.ResponseWriter, r *http.Request) { handlers.RenderAccountPage(w, r, o.renderer, o.sessions) }))
	http.HandleFunc("/register", o.makeHandler(func(w http.ResponseWriter, r *http.Request) { handlers.Register(w, r, o.db, o.renderer) }))
}

func (o *Server) Start() error {
	o.cfg = &model.Config{
		ListenPort: ":3000",
		ConnectString: `
			host=localhost
			port=5432 
			user=postgres 
			dbname=gowebapp 
			password=psqlpass 
			sslmode=disable`,
	}

	var err error
	o.db, err = db.NewPstgrDB(o.cfg)
	if err != nil {
		log.Printf("Error initializing database: %v\n", err)
		return err
	}

	o.renderer = view.NewRenderer()
	if err := o.renderer.ParseFiles(
		"view/html/register.html",
		"view/html/login.html",
		"view/html/account.html"); err != nil {
		log.Printf("Error parsing html files: %v\n", err)
	}

	o.server = &http.Server{
		Addr:         o.cfg.ListenPort,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	o.sessions = sessions.NewManager(o.db)
	o.registerHandlers()
	if err := o.server.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func (o *Server) makeHandler(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r)
	}
}
