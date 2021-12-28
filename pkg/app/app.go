/*
 * Copyright 2021 SuperPony <superponyyy@gmail.com>. All rights reserved.
 * Use of this source code is governed by a MIT style
 * license that can be found in the LICENSE file.
 */

package app

import (
	cliflag "blog-go/pkg/cli/flag"
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"os"
)

var (
	progressMessage = color.GreenString("==>")

	commandDesc = `Welcome to api-server`

	usageTemplate = fmt.Sprintf(`%s{{if .Runnable}}
  %s{{end}}{{if .HasAvailableSubCommands}}
  %s{{end}}{{if gt (len .Aliases) 0}}

%s
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

%s
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

%s{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  %s {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

%s
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

%s
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

%s{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "%s --help" for more information about a command.{{end}}
`,
		color.CyanString("Usage:"),
		color.GreenString("{{.UseLine}}"),
		color.GreenString("{{.CommandPath}} [command]"),
		color.CyanString("Aliases:"),
		color.CyanString("Examples:"),
		color.CyanString("Available Commands:"),
		color.GreenString("{{rpad .Name .NamePadding }}"),
		color.CyanString("Flags:"),
		color.CyanString("Global Flags:"),
		color.CyanString("Additional help topics:"),
		color.GreenString("{{.CommandPath}} [command]"),
	)
)

type App struct {
	use   string // 应用名称
	short string
	long  string
	// options
	options CliOptions
	cmd     *cobra.Command
	// 子命令
	commands []*Command
	// 非标志参数验证函数
	args cobra.PositionalArgs
	// 允许的非标志参数
	validArgs     []string
	runFunc       RunFunc
	silenceUsage  bool
	silenceErrors bool
}

type Option func(*App)

type RunFunc func(basename string) error

func WithLong(desc string) Option {
	return func(app *App) {
		app.long = desc
	}
}

func WithFlags(flags CliOptions) Option {
	return func(app *App) {
		app.options = flags
	}
}

func WithArgs(args cobra.PositionalArgs) Option {
	return func(app *App) {
		app.args = args
	}
}

func WithValidArgs(validArgs []string) Option {
	return func(app *App) {
		app.validArgs = validArgs
	}
}

func WithRunFunc(runFunc RunFunc) Option {
	return func(app *App) {
		app.runFunc = runFunc
	}
}

func WithSilenceUsage(silenceUsage bool) Option {
	return func(app *App) {
		app.silenceUsage = silenceUsage
	}
}

func WithSilenceErrors(silenceErrors bool) Option {
	return func(app *App) {
		app.silenceUsage = silenceErrors
	}
}

// NewApp 用户创建新的应用
// 	use 命令名称
// 	short 短介绍
func NewApp(use string, short string, opts ...Option) *App {
	app := &App{
		use:           use,
		short:         short,
		silenceUsage:  true,
		silenceErrors: true,
	}

	for _, opt := range opts {
		opt(app)
	}

	app.buildCmd()

	return app
}

func (a *App) buildCmd() {
	cmd := cobra.Command{
		Use:           FormatUseName(a.use),
		Short:         a.short,
		Long:          a.long,
		SilenceUsage:  a.silenceUsage,
		SilenceErrors: a.silenceErrors,
		Args:          a.args,
	}
	cmd.SetOut(os.Stdout)
	cmd.SetErr(os.Stderr)
	cliflag.InitFlags(cmd.Flags())
	cmd.Flags().SortFlags = true

	// 如果子命令不为空，则追加子命令
	if len(a.commands) > 0 {
		for _, command := range a.commands {
			cmd.AddCommand(command.cobraCommand())
		}

		cmd.SetHelpCommand(helpCommand(FormatUseName(a.use)))
	}

	if a.runFunc != nil {
		cmd.RunE = a.runE
	}

	var namedFlagSets cliflag.NamedFlagSets
	if a.options != nil {
		namedFlagSets = a.options.Flags()
		fs := cmd.Flags()
		for _, f := range namedFlagSets.FlagSets {
			fs.AddFlagSet(f)
		}
	}

	addCmdTemplate(&cmd, namedFlagSets)

	a.cmd = &cmd
}

func (a *App) runE(cmd *cobra.Command, args []string) error {
	// cliflag.PrintFlags(cmd.Flags())

	if a.options != nil {
		if err := a.applyOptionRules(); err != nil {
			return err
		}
	}

	if a.runFunc != nil {
		return a.runFunc(a.use)
	}

	return nil
}

func (a *App) applyOptionRules() error {
	if completableOptions, ok := a.options.(CompletableOptions); ok {
		if err := completableOptions.Complete(); err != nil {
			return err
		}
	}
	// TODO 需要一个 error 包
	if errs := a.options.Validate(); len(errs) > 0 {
		return errs[0]
	}

	// if printableOptions, ok := a.options.(PrintableOptions); ok && !a.silence {
	// 	fmt.Printf("%v Config: `%s`", progressMessage, printableOptions.String())
	//
	// }

	return nil
}

func (a *App) Run() {
	if err := a.cmd.Execute(); err != nil {
		fmt.Printf("%s \n", color.RedString("Error: %v", err.Error()))
		os.Exit(1)
	}
}

func (a App) Command() *cobra.Command {
	return a.cmd
}
