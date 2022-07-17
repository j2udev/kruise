package kruise

import (
	"github.com/j2udevelopment/kruise/pkg/kruise/schema/latest"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/thoas/go-funk"
)

type (
	Kommand struct {
		Cmd  *cobra.Command
		Opts *[]Option
	}

	Option struct {
		latest.Option
	}

	KommandBuilder struct {
		Kommand
	}

	IKommandBuilder interface {
		WithAliases(a []string) IKommandBuilder
		WithShortDescription(d string) IKommandBuilder
		WithLongDescription(d string) IKommandBuilder
		WithExample(e string) IKommandBuilder
		WithOptions(o []Option) IKommandBuilder
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

func (b *KommandBuilder) WithAliases(aliases []string) IKommandBuilder {
	b.Cmd.Aliases = aliases
	return b
}

func (b *KommandBuilder) WithShortDescription(desc string) IKommandBuilder {
	b.Cmd.Short = desc
	return b
}

func (b *KommandBuilder) WithLongDescription(desc string) IKommandBuilder {
	b.Cmd.Long = desc
	return b
}

func (b *KommandBuilder) WithExample(example string) IKommandBuilder {
	b.Cmd.Example = example
	return b
}

func (b *KommandBuilder) WithValidArgs(validArgs []string) IKommandBuilder {
	b.Cmd.ValidArgs = validArgs
	return b
}

func (b *KommandBuilder) WithArgs(args cobra.PositionalArgs) IKommandBuilder {
	b.Cmd.Args = args
	return b
}

func (b *KommandBuilder) WithArgAliases(aliases []string) IKommandBuilder {
	b.Cmd.ArgAliases = aliases
	return b
}

func (b *KommandBuilder) WithSubKommands(kmds ...Kommand) IKommandBuilder {
	funk.ForEach(kmds, func(kmd Kommand) {
		b.Cmd.AddCommand(kmd.Cmd)
	})
	return b
}

func (b *KommandBuilder) WithOptions(opts []Option) IKommandBuilder {
	b.Opts = &opts
	b.WithKruiseTemplate()
	return b
}

func (b *KommandBuilder) WithFlags(flags *pflag.FlagSet) IKommandBuilder {
	b.Cmd.Flags().AddFlagSet(flags)
	return b
}

func (b *KommandBuilder) WithPersistentFlags(flags *pflag.FlagSet) IKommandBuilder {
	b.Cmd.PersistentFlags().AddFlagSet(flags)
	return b
}

func (b *KommandBuilder) WithPreRunFunc(f func(*cobra.Command, []string)) IKommandBuilder {
	b.Cmd.PreRun = f
	return b
}

func (b *KommandBuilder) WithPreRunEFunc(f func(*cobra.Command, []string) error) IKommandBuilder {
	b.Cmd.PreRunE = f
	return b
}

func (b *KommandBuilder) WithPersistentPreRunFunc(f func(*cobra.Command, []string)) IKommandBuilder {
	b.Cmd.PersistentPreRun = f
	return b
}

func (b *KommandBuilder) WithPersistentPreRunEFunc(f func(*cobra.Command, []string) error) IKommandBuilder {
	b.Cmd.PersistentPreRunE = f
	return b
}

func (b *KommandBuilder) WithRunFunc(f func(*cobra.Command, []string)) IKommandBuilder {
	b.Cmd.Run = f
	return b
}

func (b *KommandBuilder) WithRunEFunc(f func(*cobra.Command, []string) error) IKommandBuilder {
	b.Cmd.RunE = f
	return b
}

func (b *KommandBuilder) WithPostRunFunc(f func(*cobra.Command, []string)) IKommandBuilder {
	b.Cmd.PostRun = f
	return b
}

func (b *KommandBuilder) WithPostRunEFunc(f func(*cobra.Command, []string) error) IKommandBuilder {
	b.Cmd.PostRunE = f
	return b
}

func (b *KommandBuilder) WithPersistentPostRunFunc(f func(*cobra.Command, []string)) IKommandBuilder {
	b.Cmd.PersistentPostRun = f
	return b
}

func (b *KommandBuilder) WithPersistentPostRunEFunc(f func(*cobra.Command, []string) error) IKommandBuilder {
	b.Cmd.PersistentPostRunE = f
	return b
}

func (b *KommandBuilder) WithUsageTemplate(template string) IKommandBuilder {
	b.Cmd.SetUsageTemplate(template)
	return b
}

func (b *KommandBuilder) WithHelpTemplate(template string) IKommandBuilder {
	b.Cmd.SetHelpTemplate(template)
	return b
}

func (b *KommandBuilder) WithUsageFunc(function func(*cobra.Command) error) IKommandBuilder {
	b.Cmd.SetUsageFunc(function)
	return b
}

func (b *KommandBuilder) WithHelpFunc(function func(*cobra.Command, []string)) IKommandBuilder {
	b.Cmd.SetHelpFunc(function)
	return b
}

func (b *KommandBuilder) WithKruiseTemplate() IKommandBuilder {
	b.WithUsageTemplate(UsageTemplate()).
		WithHelpTemplate(UsageTemplate())
	kmd := b.Build()
	b.WithUsageFunc(UsageFunc(kmd)).
		WithHelpFunc(HelpFunc(kmd))
	return b
}

func (b *KommandBuilder) Version(version string) IKommandBuilder {
	b.Cmd.Version = version
	return b
}

func (b *KommandBuilder) Deprecated(deprecated string) IKommandBuilder {
	b.Cmd.Deprecated = deprecated
	return b
}

func (b *KommandBuilder) Build() Kommand {
	return Kommand{
		Cmd:  b.Cmd,
		Opts: b.Opts,
	}
}

func (k Kommand) Execute() error {
	return k.Cmd.Execute()
}
