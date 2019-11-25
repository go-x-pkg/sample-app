package server

import (
	"sync"

	"github.com/go-x-pkg/log"
)

type AppContext struct {
	cfgLock sync.RWMutex
	config  *config

	loggersLock sync.RWMutex
	lggrs       log.Loggers
}

func (ac *AppContext) setCfg(v *config) *config {
	ac.cfgLock.Lock()
	defer ac.cfgLock.Unlock()
	old := ac.config
	ac.config = v
	return old
}

func (ac *AppContext) cfg() *config {
	ac.cfgLock.RLock()
	defer ac.cfgLock.RUnlock()
	return ac.config
}

func (ac *AppContext) setLoggers(v log.Loggers) log.Loggers {
	ac.loggersLock.Lock()
	defer ac.loggersLock.Unlock()
	old := ac.lggrs
	ac.lggrs = v
	return old
}

func (ac *AppContext) loggers() log.Loggers {
	ac.loggersLock.RLock()
	defer ac.loggersLock.RUnlock()
	return ac.lggrs
}

func (ac *AppContext) init() {
}
