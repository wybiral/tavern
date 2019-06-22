package app

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/wybiral/torgo"
)

type tor struct {
	Onion      *torgo.Onion
	Controller *torgo.Controller
}

func newTor(ct *TorConfig, cs *ServerConfig) (*tor, error) {
	addr := fmt.Sprintf("%s:%d", ct.Controller.Host, ct.Controller.Port)
	ctrl, err := torgo.NewController(addr)
	if err != nil {
		return nil, err
	}
	if len(ct.Controller.Password) > 0 {
		err = ctrl.AuthenticatePassword(ct.Controller.Password)
	} else {
		err = ctrl.AuthenticateCookie()
		if err != nil {
			err = ctrl.AuthenticateNone()
		}
	}
	if err != nil {
		return nil, err
	}
	t := &tor{
		Controller: ctrl,
	}
	err = t.startOnion(cs)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (t *tor) startOnion(cs *ServerConfig) error {
	addr := fmt.Sprintf("%s:%d", cs.Host, cs.Port)
	t.Onion = &torgo.Onion{}
	t.Onion.Ports = map[int]string{80: addr}
	raw, err := ioutil.ReadFile("onion.key")
	if err != nil {
		return err
	}
	pk := strings.TrimSpace(string(raw))
	parts := strings.SplitN(pk, ":", 2)
	t.Onion.PrivateKeyType = parts[0]
	t.Onion.PrivateKey = parts[1]
	return t.Controller.AddOnion(t.Onion)
}
