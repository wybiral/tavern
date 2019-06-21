package app

import (
	"net"
	"net/http"

	"github.com/gorilla/mux"
)

// App manages main application.
type App struct {
	Config   *Config
	Router   *mux.Router
	Listener net.Listener
}

// NewApp returns new App from Config.
func NewApp(c *Config) (*App, error) {
	if c == nil {
		c = DefaultConfig()
	}
	// setup router
	r := mux.NewRouter().StrictSlash(true)
	fs := http.FileServer(http.Dir("./public/"))
	r.PathPrefix("/").Handler(fs)
	// setup listener
	ln, err := newListener(c.Server)
	if err != nil {
		return nil, err
	}
	// return App
	a := &App{
		Config:   c,
		Router:   r,
		Listener: ln,
	}
	return a, nil
}

// Run starts App server.
func (a *App) Run() error {
	return http.Serve(a.Listener, a.Router)
}
