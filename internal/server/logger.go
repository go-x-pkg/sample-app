package server

import (
	"github.com/cihub/seelog"
	"github.com/go-x-pkg/log"
	"github.com/go-x-pkg/xseelog"
)

// default global logger
func logger() seelog.LoggerInterface { return seelog.Current }

// default global log func
func logFn(l log.Level, msg string, a ...interface{}) {
	log.Logf(logger(), l, msg, a...)
}

func initLoggers(cfg *xseelog.Config) (log.Loggers, error) {
	loggers, err := cfg.Loggers()
	if err != nil {
		return nil, err
	}

	loggers.ReplaceRoot()
	return loggers, nil
}
