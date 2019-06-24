package app

import (
	"fmt"
	"net"
	"net/http"

	"github.com/gorilla/mux"
)

// App manages main application.
type App struct {
	Config   *Config
	Router   *mux.Router
	Tor      *tor
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
	// setup tor
	t, err := newTor(c.Tor)
	if err != nil {
		return nil, err
	}
	// setup App
	a := &App{
		Config:   c,
		Router:   r,
		Tor:      t,
		Listener: ln,
	}
	return a, nil
}

// Run starts App server.
func (a *App) Run() error {
	cs := a.Config.Server
	onion, err := a.Tor.OnionKey.Onion()
	if err != nil {
		return err
	}
	onion.Ports[80] = fmt.Sprintf("%s:%d", cs.Host, cs.Port)
	err = a.Tor.Controller.AddOnion(onion)
	if err != nil {
		return err
	}
	return http.Serve(a.Listener, a.Router)
}
