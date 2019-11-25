package server

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/go-x-pkg/log"
	daemon "github.com/sevlyar/go-daemon"
	"github.com/spf13/pflag"
)

type App struct {
	flags     flags
	daemonCtx *daemon.Context
	ctx       AppContext
}

func (a *App) daemonize() error {
	cfg := &a.ctx.cfg().Daemon

	a.daemonCtx = &daemon.Context{
		Umask: cfg.Umask,
		Args:  os.Args,
	}

	if v := cfg.Pidfile; v != "" {
		a.daemonCtx.PidFileName = v
		a.daemonCtx.PidFilePerm = cfg.PidfileMode
	}

	if v := cfg.WorkDir; v != "" {
		a.daemonCtx.WorkDir = v
	}

	child, err := a.daemonCtx.Reborn()
	if err != nil {
		return fmt.Errorf("daemon reborn failed: %w", err)
	}

	if child != nil { // parent exit
		os.Exit(0)
	}

	return nil
}

func (a *App) run() (outErr error) {
	if a.ctx.cfg().Daemonize {
		if err := a.daemonize(); err != nil {
			return fmt.Errorf("daemonization failed: %w", err)
		}

		if pf := a.ctx.cfg().Daemon.Pidfile; pf != "" {
			logFn(log.Info, "writing pidfile to %q as %s", pf, a.ctx.cfg().Daemon.PidfileMode)
		}
	}

	timeoutWorkersDone := a.ctx.cfg().Timeout.WorkersDone

	serverErrChan := make(chan struct{})
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		defer cancel()

		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan,
			syscall.SIGINT, os.Interrupt, // CTRL-C
			syscall.SIGTERM,
			syscall.SIGQUIT,

			syscall.SIGUSR1, // postrotate
			syscall.SIGHUP,  // reload
		)

		for {
			select {
			case sig := <-signalChan:
				switch sig {
				case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
					logFn(log.Info, "(SIGINT SIGTERM SIGQUIT) will shutdown")

					if a.daemonCtx != nil {
						a.daemonCtx.Release()
						a.daemonCtx = nil
					}

					return

				case syscall.SIGUSR1: // postrotate
					logFn(log.Info, "(SIGUSR1) postrotate logs reloading")
					if err := a.initialize(actionMain, false, true); err != nil {
						logFn(log.Error, "(SIGUSR1) postrotate error: %w", err)
					} else {
						logFn(log.Info, "(SIGUSR1) postrotate OK")
					}

				case syscall.SIGHUP:
					logFn(log.Info, "(SIGHUP) reloading configuration file")
					if err := a.initialize(actionMain, true, true); err != nil {
						logFn(log.Error, "(SIGHUP) reload error: %w", err)
					} else {
						a.ctx.cfg().DumpToFn(func(buf *bytes.Buffer) {
							logFn(log.Debug,
								"(SIGHUP) using configuration:\n%s", buf.String())
						})
						logFn(log.Info, "(SIGHUP) reload OK")
					}
				}
			case <-serverErrChan:
				return
			}
		}
	}()

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()

		// start some worker here
	}()

	logFn(log.Info,
		"application with (:pid %d) started", os.Getpid())

	a.ctx.cfg().DumpToFn(func(buf *bytes.Buffer) {
		logFn(log.Debug,
			"using configuration:\n%s", buf.String())
	})

	wg.Add(1)
	go func() {
		// start some HTTP servers here
	}()

	wg.Wait()

	return outErr
}

// apply flags, reload config, reload loggers
func (a *App) initialize(actn action, resetConfig bool, resetLoggers bool) error {
	cfg := new(config)

	if err := cfg.build(&a.flags, actn); err != nil {
		return err
	} else if resetConfig {
		a.ctx.setCfg(cfg)
	}

	if loggers, err := initLoggers(cfg.Log); err != nil {
		return err
	} else if resetLoggers {
		log.Close( // close old logger if any
			a.ctx.setLoggers(
				loggers))
	}

	return nil
}

func (a *App) main(actn action, args []string, flags *pflag.FlagSet) (err error) {
	a.flags.set = flags
	a.ctx.init()

	if err := a.initialize(actn, true, true); err != nil {
		return fmt.Errorf("intialization failed: %w", err)
	}

	if a.flags.printConfig {
		actn = actionPrintConfig
	}

	switch actn {
	case actionMain:
		err = a.run()
	case actionPrintConfig:
		a.doPrintConfig()
	case actionVersion:
		a.doVersion()
	}

	return err
}

func (a *App) Main(actn action, args []string, flags *pflag.FlagSet) error {
	return a.main(actn, args, flags)
}
