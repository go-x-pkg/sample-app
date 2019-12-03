package server

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/go-x-pkg/bufpool"
	"github.com/go-x-pkg/dumpctx"
	"github.com/go-x-pkg/fnscli"
	"github.com/go-x-pkg/log"
	"github.com/go-x-pkg/servers"
	"github.com/go-x-pkg/xseelog"
	"github.com/spf13/pflag"
)

const (
	defaultPathConfig string = "/etc/sample-app-server/config.yml"

	defaultDaemonPidfile     string      = "/run/sample-app-server/sample-app-server.pid"
	defaultDaemonPidfileMode os.FileMode = 0644
	defaultDaemonWorkdir     string      = ""

	defaultDirLog string = "/var/log/sampla-app-server"

	defaultServerINETHost string = "0.0.0.0"
	defaultServerINETPort int    = 8000
	defaultServerUNIXAddr string = "/run/sample-app-server/sample-app-server.sock"

	defaultTimeoutWorkersDone time.Duration = 5 * time.Second

	defaultPeriodMemstats time.Duration = 60 * time.Second
)

type flags struct {
	pathConfig string

	printConfig bool

	foreground bool
	daemonize  bool

	serverPort int

	timeout struct {
		workersDone time.Duration
	}

	period struct {
		memstats time.Duration
	}

	logDisableFile    bool
	logDisableConsole bool
	vv                bool
	vvv               bool

	// corresponding flag-set
	set *pflag.FlagSet
}

type config struct {
	Daemonize bool `yaml:"daemonize"`
	Daemon    struct {
		Pidfile     string      `yaml:"pidfile"`
		PidfileMode os.FileMode `yaml:"pidfile-mode"`
		WorkDir     string      `yaml:"workdir"`
		Umask       int         `yaml:"umask"`
	} `yaml:"daemon"`

	Servers servers.Configs `yaml:"servers"`

	Timeout struct {
		WorkersDone time.Duration `yaml:"workers-done"`
	} `yaml:"timeout"`

	Period struct {
		Memstats time.Duration `yaml:"memstats"`
	} `yaml:"period"`

	Log *xseelog.Config `yaml:"log"`
}

func (c *config) fromFile(flags *flags) error {
	path := flags.pathConfig

	return fnscli.DecodeYAMLFromPath(
		c,                                      // destination to decode
		path,                                   // path to yaml config file
		fnscli.IsPFlagSet(flags.set, "config"), // path to config file was forced
	)
}

func (c *config) defaultize() error {
	if c.Daemon.Pidfile == "" {
		c.Daemon.Pidfile = defaultDaemonPidfile
	}

	if c.Daemon.PidfileMode == 0 {
		c.Daemon.PidfileMode = defaultDaemonPidfileMode
	}

	if c.Daemon.WorkDir == "" {
		c.Daemon.WorkDir = defaultDaemonWorkdir
	}

	if e := c.Servers.Defaultize(
		defaultServerINETHost,
		defaultServerINETPort,
		defaultServerUNIXAddr,
	); e != nil {
		return e
	}

	if len(c.Servers) == 0 {
		c.Servers.PushINETIfNotExists(defaultServerINETHost, defaultServerINETPort)
	}

	if c.Timeout.WorkersDone == 0 {
		c.Timeout.WorkersDone = defaultTimeoutWorkersDone
	}

	if c.Period.Memstats == 0 {
		c.Period.Memstats = defaultPeriodMemstats
	}

	if c.Log == nil {
		c.Log = xseelog.NewConfig()
	}

	if c.Log.Dir == "" {
		c.Log.Dir = defaultDirLog
	}

	c.Log.Ensure("app", "", log.Info, log.Critical)
	c.Log.Ensure("http", "(:http)", log.Info, log.Critical)
	c.Log.Ensure("memstats", "(:memstats)", log.Info, log.Critical)

	return nil
}

func (c *config) fromFlags(flags *flags, actn action) error {
	if fnscli.IsPFlagSet(flags.set, "foreground") && flags.foreground {
		c.Daemonize = false
	} else {
		if (fnscli.IsPFlagSet(flags.set, "background") ||
			fnscli.IsPFlagSet(flags.set, "daemonize")) &&
			flags.daemonize {
			c.Daemonize = true
		}
	}

	if fnscli.IsPFlagSet(flags.set, "timeout-worker-done") {
		c.Timeout.WorkersDone = flags.timeout.workersDone
	}

	if flags.logDisableConsole {
		c.Log.DisableConsole = true
	}

	if flags.logDisableFile {
		c.Log.DisableFile = true
	}

	if c.Log.DisableConsole && c.Log.DisableFile {
		c.Log.DisableConsole = false
		c.Log.Quiet()
	}

	if flags.vv {
		c.Log.VV()
	}

	if flags.vvv {
		c.Log.VVV()
	}

	return nil
}

func (c *config) validate(actn action) error {
	return nil
}

func (c *config) build(flags *flags, actn action) error {
	if err := c.fromFile(flags); err != nil {
		return err
	}

	c.defaultize()

	if err := c.fromFlags(flags, actn); err != nil {
		return err
	}

	if err := c.validate(actn); err != nil {
		return err
	}

	return nil
}

func (c *config) DumpToFn(cb func(*bytes.Buffer)) {
	buf := bufpool.NewBuf()
	buf.Reset()
	c.Dump(buf)
	cb(&buf.Buffer)
	buf.Release()
}

func (c *config) Dump(w io.Writer) {
	ctx := dumpctx.Ctx{}
	ctx.Init()

	fmt.Fprintf(w, "%sdaemonize: %t\n", ctx.Indent(), c.Daemonize)
	fmt.Fprintf(w, "%sdaemon:\n", ctx.Indent())
	ctx.Wrap(func() {
		fmt.Fprintf(w, "%spidfile: %s\n", ctx.Indent(), c.Daemon.Pidfile)
		fmt.Fprintf(w, "%spidfile-mode: 0%03o | %s\n", ctx.Indent(), c.Daemon.PidfileMode, c.Daemon.PidfileMode)
		fmt.Fprintf(w, "%sworkdir: %s\n", ctx.Indent(), c.Daemon.WorkDir)
		fmt.Fprintf(w, "%sumask: %03o\n", ctx.Indent(), c.Daemon.Umask)
	})

	fmt.Fprintf(w, "%sservers", ctx.Indent())
	{
		c.Servers.Dump(&ctx, w)
	}

	fmt.Fprintf(w, "%stimeout:\n", ctx.Indent())
	ctx.Wrap(func() {
		fmt.Fprintf(w, "%sworkers-done: %s\n", ctx.Indent(), c.Timeout.WorkersDone)
	})

	fmt.Fprintf(w, "%speriod:\n", ctx.Indent())
	ctx.Wrap(func() {
		fmt.Fprintf(w, "%smemstats: %s\n", ctx.Indent(), c.Period.Memstats)
		ctx.Leave()
	})

	fmt.Fprintf(w, "%slogs:\n", ctx.Indent())
	ctx.Wrap(func() { c.Log.Dump(&ctx, w) })
}
