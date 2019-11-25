package server

import (
	"github.com/spf13/cobra"
)

type action uint8

const (
	actionUnknown action = iota
	actionMain
	actionPrintConfig
	actionVersion
	// action-migrate, action-show-stats, etc
)

var actionText = map[action]string{
	actionUnknown:     "unknown",
	actionMain:        "main",
	actionPrintConfig: "print-config",
	actionVersion:     "version",
}

func (a action) String() string { return actionText[a] }

func (a *App) applyMainFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&a.flags.pathConfig,
		"config", "c",
		defaultPathConfig, "path to configuration file")

	cmd.Flags().BoolVar(&a.flags.printConfig,
		"print-config", false,
		"print config and exit")

	cmd.Flags().BoolVar(&a.flags.foreground,
		"foreground", false,
		"force app to run in foreground (don't daemonize)")

	cmd.Flags().BoolVar(&a.flags.daemonize,
		"daemonize", false,
		"force app to run in background (daemonize)")

	cmd.Flags().BoolVar(&a.flags.daemonize,
		"background", false,
		"force app to run in background (daemonize)")

	cmd.Flags().DurationVar(&a.flags.timeout.workersDone,
		"timeout-workers-done", defaultTimeoutWorkersDone,
		"wait duration for workers to shutdown gracefully. otherwise force shutdown")

	cmd.Flags().BoolVar(&a.flags.logDisableConsole,
		"log-disable-console", false,
		"disable console stderr/stdout logging")

	cmd.Flags().BoolVar(&a.flags.logDisableFile,
		"log-disable-file", false,
		"disable file logging")

	cmd.Flags().BoolVar(&a.flags.vv,
		"vv", false,
		"all loggers log level forced to debug (very verbose)")

	cmd.Flags().BoolVar(&a.flags.vvv,
		"vvv", false,
		"all loggers log level forced to trace (very very verbose)")
}

func (a *App) CmdSetRunFnAndFlags(cmd *cobra.Command) {
	a.applyMainFlags(cmd)

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return a.Main(actionMain, nil, cmd.Flags())
	}

	cmdPrintConfig := &cobra.Command{
		Use:           "print-config",
		Short:         "log down config and exit",
		Aliases:       []string{"config"},
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.Main(actionPrintConfig, nil, cmd.Flags())
		},
	}
	a.applyMainFlags(cmdPrintConfig)

	cmdVersion := &cobra.Command{
		Use:           "version",
		Short:         "show version and exit",
		Aliases:       []string{"config"},
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.Main(actionVersion, nil, cmd.Flags())
		},
	}
	a.applyMainFlags(cmdVersion)

	cmd.AddCommand(
		cmdPrintConfig,
		cmdVersion,
	)

	// ... add other commands here if any
}

func (a *App) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "server [...]",
		Short:   "server, media server, http-server",
		Aliases: []string{"s", "server"},

		// SilenceUsage is an option to silence usage when an error occurs.
		SilenceUsage: true,

		// SilenceErrors is an option to quiet errors down stream.
		SilenceErrors: true,
	}

	a.CmdSetRunFnAndFlags(cmd)

	return cmd
}
