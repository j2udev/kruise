package kruise

import (
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/thoas/go-funk"
)

type (
	// Kommand wraps a cobra Command to give the Kruise CLI greater flexibility
	//
	// cobra Command: https://pkg.go.dev/github.com/spf13/cobra#Command
	Kommand struct {
		Cmd      *cobra.Command
		Opts     *Options
		Profiles *Profiles
	}

	// KommandBuilder is used to build a Kruise Kommand
	KommandBuilder struct {
		Kommand
	}

	// IKommandBuilder defines the functions for a Kruise Kommand
	IKommandBuilder interface {
		WithAliases(a []string) IKommandBuilder
		WithShortDescription(d string) IKommandBuilder
		WithLongDescription(d string) IKommandBuilder
		WithExample(e string) IKommandBuilder
		WithOptions(o Options) IKommandBuilder
		WithProfiles(p Profiles) IKommandBuilder
		WithSubKommands(k ...Kommand) IKommandBuilder
		WithFlags(f *pflag.FlagSet) IKommandBuilder
		WithPersistentFlags(f *pflag.FlagSet) IKommandBuilder
		WithValidArgs(a []string) IKommandBuilder
		WithArgs(a cobra.PositionalArgs) IKommandBuilder
		WithArgAliases(a []string) IKommandBuilder
		WithKruiseTemplate() IKommandBuilder
		WithUsageTemplate(t string) IKommandBuilder
		WithHelpTemplate(t string) IKommandBuilder
		WithUsageFunc(f func(*cobra.Command) error) IKommandBuilder
		WithHelpFunc(f func(*cobra.Command, []string)) IKommandBuilder
		WithPreRunFunc(f func(*cobra.Command, []string)) IKommandBuilder
		WithPreRunEFunc(f func(*cobra.Command, []string) error) IKommandBuilder
		WithPersistentPreRunFunc(f func(*cobra.Command, []string)) IKommandBuilder
		WithPersistentPreRunEFunc(f func(*cobra.Command, []string) error) IKommandBuilder
		WithRunFunc(f func(*cobra.Command, []string)) IKommandBuilder
		WithRunEFunc(f func(*cobra.Command, []string) error) IKommandBuilder
		WithPostRunFunc(f func(*cobra.Command, []string)) IKommandBuilder
		WithPostRunEFunc(f func(*cobra.Command, []string) error) IKommandBuilder
		WithPersistentPostRunFunc(f func(*cobra.Command, []string)) IKommandBuilder
		WithPersistentPostRunEFunc(f func(*cobra.Command, []string) error) IKommandBuilder
		Version(v string) IKommandBuilder
		Deprecated(d string) IKommandBuilder
		Build() Kommand
	}
)

// NewKmd returns a new IKommandBuilder
func NewKmd(name string) IKommandBuilder {
	cmd := cobra.Command{
		Use: name,
	}
	kmd := &KommandBuilder{
		Kommand: Kommand{
			Cmd: &cmd,
		},
	}
	return kmd
}

// WithAliases sets the alias list for the underlying cobra Command
func (b *KommandBuilder) WithAliases(aliases []string) IKommandBuilder {
	b.Cmd.Aliases = aliases
	return b
}

// WithShortDescription sets the short description for the underlying cobra
// Command
func (b *KommandBuilder) WithShortDescription(desc string) IKommandBuilder {
	b.Cmd.Short = desc
	return b
}

// WithLongDescription sets the long description for the underlying cobra
// Command
func (b *KommandBuilder) WithLongDescription(desc string) IKommandBuilder {
	b.Cmd.Long = desc
	return b
}

// WithExample sets an example for the underlying cobra Command
func (b *KommandBuilder) WithExample(example string) IKommandBuilder {
	b.Cmd.Example = example
	return b
}

// WithValidArgs sets the valid arguments for the underlying cobra Command
func (b *KommandBuilder) WithValidArgs(validArgs []string) IKommandBuilder {
	b.Cmd.ValidArgs = validArgs
	return b
}

// WithArgs sets the arguments for the underlying cobra Command
func (b *KommandBuilder) WithArgs(args cobra.PositionalArgs) IKommandBuilder {
	b.Cmd.Args = args
	return b
}

// WithArgAliases sets the argument aliases for for the underlying cobra Command
func (b *KommandBuilder) WithArgAliases(aliases []string) IKommandBuilder {
	b.Cmd.ArgAliases = aliases
	return b
}

// WithSubKommands sets the sub commands for for the underlying cobra Command
func (b *KommandBuilder) WithSubKommands(kmds ...Kommand) IKommandBuilder {
	funk.ForEach(kmds, func(kmd Kommand) {
		b.Cmd.AddCommand(kmd.Cmd)
	})
	return b
}

// WithOptions sets the Options for the Kruise Kommand
func (b *KommandBuilder) WithOptions(opts Options) IKommandBuilder {
	b.Opts = &opts
	b.WithKruiseTemplate()
	return b
}

// WithProfiles sets the Profiles for the Kruise Kommand
func (b *KommandBuilder) WithProfiles(profs Profiles) IKommandBuilder {
	b.Profiles = &profs
	b.WithKruiseTemplate()
	return b
}

// WithFlags sets the Flags for for the underlying cobra Command
func (b *KommandBuilder) WithFlags(flags *pflag.FlagSet) IKommandBuilder {
	b.Cmd.Flags().AddFlagSet(flags)
	return b
}

// WithPersistentFlags sets the persistent flags for for the underlying cobra
// Command
func (b *KommandBuilder) WithPersistentFlags(flags *pflag.FlagSet) IKommandBuilder {
	b.Cmd.PersistentFlags().AddFlagSet(flags)
	return b
}

// WithPreRunFunc sets the pre run function for the underlying cobra Command
func (b *KommandBuilder) WithPreRunFunc(f func(*cobra.Command, []string)) IKommandBuilder {
	b.Cmd.PreRun = f
	return b
}

// WithPreRunEFunc sets the pre run error function for the underlying cobra
// Command
func (b *KommandBuilder) WithPreRunEFunc(f func(*cobra.Command, []string) error) IKommandBuilder {
	b.Cmd.PreRunE = f
	return b
}

// WithPersistentPreRunFunc sets the persistent pre run function for the
// underlying cobra Command
func (b *KommandBuilder) WithPersistentPreRunFunc(f func(*cobra.Command, []string)) IKommandBuilder {
	b.Cmd.PersistentPreRun = f
	return b
}

// WithPersistentPreRunEFunc sets the persistent pre run error function for the
// underlying cobra Command
func (b *KommandBuilder) WithPersistentPreRunEFunc(f func(*cobra.Command, []string) error) IKommandBuilder {
	b.Cmd.PersistentPreRunE = f
	return b
}

// WithRunFunc sets the run function for the underlying cobra Command
func (b *KommandBuilder) WithRunFunc(f func(*cobra.Command, []string)) IKommandBuilder {
	b.Cmd.Run = f
	return b
}

// WithRunEFunc sets the run error function for the underlying cobra Command
func (b *KommandBuilder) WithRunEFunc(f func(*cobra.Command, []string) error) IKommandBuilder {
	b.Cmd.RunE = f
	return b
}

// WithPostRunFunc sets the post run function for the underlying cobra Command
func (b *KommandBuilder) WithPostRunFunc(f func(*cobra.Command, []string)) IKommandBuilder {
	b.Cmd.PostRun = f
	return b
}

// WithPostRunEFunc sets the post run error function for the underlying cobra
// Command
func (b *KommandBuilder) WithPostRunEFunc(f func(*cobra.Command, []string) error) IKommandBuilder {
	b.Cmd.PostRunE = f
	return b
}

// WithPersistentPostRunFunc sets the persistent post run function for the
// underlying cobra Command
func (b *KommandBuilder) WithPersistentPostRunFunc(f func(*cobra.Command, []string)) IKommandBuilder {
	b.Cmd.PersistentPostRun = f
	return b
}

// WithPersistentPostERunFunc sets the persistent post run error function for
// the underlying cobra Command
func (b *KommandBuilder) WithPersistentPostRunEFunc(f func(*cobra.Command, []string) error) IKommandBuilder {
	b.Cmd.PersistentPostRunE = f
	return b
}

// WithUsageTemplate sets the usage template for the underlying cobra Command
func (b *KommandBuilder) WithUsageTemplate(template string) IKommandBuilder {
	b.Cmd.SetUsageTemplate(template)
	return b
}

// WithHelpTemplate sets the help template for the underlying cobra Command
func (b *KommandBuilder) WithHelpTemplate(template string) IKommandBuilder {
	b.Cmd.SetHelpTemplate(template)
	return b
}

// WithUsageFunc sets the usage function for the underlying cobra Command
func (b *KommandBuilder) WithUsageFunc(function func(*cobra.Command) error) IKommandBuilder {
	b.Cmd.SetUsageFunc(function)
	return b
}

// WithHelpFunc sets the help function for the underlying cobra Command
func (b *KommandBuilder) WithHelpFunc(function func(*cobra.Command, []string)) IKommandBuilder {
	b.Cmd.SetHelpFunc(function)
	return b
}

// WithKruiseTemplate is a helper function for more easily setting cobra Command
// Usage
func (b *KommandBuilder) WithKruiseTemplate() IKommandBuilder {
	kmd := b.Build()
	b.WithUsageFunc(kmd.UsageFunc()).
		WithHelpFunc(kmd.HelpFunc())
	return b
}

// Version is used to define the version of the underlying cobra Command
func (b *KommandBuilder) Version(version string) IKommandBuilder {
	b.Cmd.Version = version
	return b
}

// Deprecated is used to define whether the underlying cobra Command should be
// listed as deprecated
func (b *KommandBuilder) Deprecated(deprecated string) IKommandBuilder {
	b.Cmd.Deprecated = deprecated
	return b
}

// Build returns a Kommand from a KommandBuilder
func (b *KommandBuilder) Build() Kommand {
	return Kommand{
		Cmd:      b.Cmd,
		Opts:     b.Opts,
		Profiles: b.Profiles,
	}
}

// UsageFunc overrides the default UsageFunc used by cobra to facilitate showing
// command options
func (k Kommand) UsageFunc() (f func(*cobra.Command) error) {
	return func(c *cobra.Command) error {
		w := tabwriter.NewWriter(os.Stdout, 8, 8, 8, ' ', 0)
		err := tmpl(w, k.UsageTemplate(), k)
		if err != nil {
			c.PrintErrln(err)
		}
		return err
	}
}

// HelpFunc overrides the default HelpFunc used by cobra to facilitate showing
// command options
func (k Kommand) HelpFunc() func(*cobra.Command, []string) {
	return func(c *cobra.Command, s []string) {
		w := tabwriter.NewWriter(os.Stdout, 3, 3, 3, ' ', 0)
		err := tmpl(w, k.UsageTemplate(), k)
		if err != nil {
			c.PrintErrln(err)
		}
	}
}

// UsageTemplate is used to override the cobra UsageTemplate to facilitate
// options and other CLI parameters specific to Kruise
func (k Kommand) UsageTemplate() string {
	return `Usage:{{if .Cmd.Runnable}}
  {{.Cmd.UseLine}}{{end}} [options]|[profiles]{{if .Cmd.HasAvailableSubCommands}}
  {{.Cmd.CommandPath}} [command]{{end}}{{if gt (len .Cmd.Aliases) 0}}

Aliases:
  {{.Cmd.NameAndAliases}}{{end}}{{if .Cmd.HasExample}}

Examples:
{{.Cmd.Example}}{{end}}{{if .HasOptions}}

Options:{{range .Opts }}
  {{.Args}}	{{.Desc}}{{end}}{{end}}

Profiles:{{range .Profiles }}
  {{.Args}}	{{.Desc}}
    ↳ Options:	{{range .Items}}{{.}} {{end}}{{end}}{{if .Cmd.HasAvailableLocalFlags}}

Flags:
{{.Cmd.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .Cmd.HasAvailableInheritedFlags}}

Global Flags:
{{.Cmd.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .Cmd.HasHelpSubCommands}}

Additional help topics:{{range .Cmd.Commands}}{{if .Cmd.IsAdditionalHelpTopicCommand}}
  {{rpad .Cmd.CommandPath .Cmd.CommandPathPadding}} {{.Cmd.Short}}{{end}}{{end}}{{end}}{{if .Cmd.HasAvailableSubCommands}}

Use "{{.Cmd.CommandPath}} [command] --help" for more information about a command.{{end}}
`
}

// HasOptions is used to determine if the Kommand has options
func (k Kommand) HasOptions() bool {
	if k.Opts == nil || len(*k.Opts) == 0 {
		return false
	}
	return true
}

// HasProfiles is used to determine if the Kommand has profiles
func (k Kommand) HasProfiles() bool {
	if k.Profiles == nil || len(*k.Profiles) == 0 {
		return false
	}
	return true
}

// Execute is used to execute the underlying cobra Command
func (k Kommand) Execute() error {
	return k.Cmd.Execute()
}
