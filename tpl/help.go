package tpl

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"
	"text/template"
	"unicode"

	c "github.com/j2udevelopment/kruise/pkg/config"
	"github.com/spf13/cobra"
)

var templateFuncs = template.FuncMap{
	"trim":                    strings.TrimSpace,
	"trimRightSpace":          trimRightSpace,
	"trimTrailingWhitespaces": trimRightSpace,
	"rpad":                    rpad,
}

func trimRightSpace(s string) string {
	return strings.TrimRightFunc(s, unicode.IsSpace)
}

// rpad adds padding to the right of a string.
func rpad(s string, padding int) string {
	template := fmt.Sprintf("%%-%ds", padding)
	return fmt.Sprintf(template, s)
}

// tmpl executes the given template text on data, writing the result to w.
func tmpl(w io.Writer, text string, data interface{}) error {
	t := template.New("tmpl")
	t.Funcs(templateFuncs)
	template.Must(t.Parse(text))
	return t.Execute(w, data)
}

// UsageTemplate returns usage template for the command.
func UsageTemplate() string {
	return `Usage:{{if .Cmd.Runnable}}
  {{.Cmd.UseLine}}{{end}}{{if .Cmd.HasAvailableSubCommands}}
  {{.Cmd.CommandPath}} [command]{{end}}{{if gt (len .Cmd.Aliases) 0}}

Aliases:
  {{.Cmd.NameAndAliases}}{{end}}{{if .Cmd.HasExample}}

Examples:
{{.Cmd.Example}}{{end}}

Available Options:{{range .Opts }}
  {{.Arguments}}	{{.Description}}{{end}}{{if .Cmd.HasAvailableLocalFlags}}

Flags:
{{.Cmd.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .Cmd.HasAvailableInheritedFlags}}

Global Flags:
{{.Cmd.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .Cmd.HasHelpSubCommands}}

Additional help topics:{{range .Cmd.Commands}}{{if .Cmd.IsAdditionalHelpTopicCommand}}
  {{rpad .Cmd.CommandPath .Cmd.CommandPathPadding}} {{.Cmd.Short}}{{end}}{{end}}{{end}}{{if .Cmd.HasAvailableSubCommands}}

Use "{{.Cmd.CommandPath}} [command] --help" for more information about a command.{{end}}
`
}

// UsageFunc overrides the default UsageFunc used by Cobra to facilitate showing command options
func UsageFunc(wrapper c.CommandWrapper) (f func(*cobra.Command) error) {
	return func(c *cobra.Command) error {
		w := tabwriter.NewWriter(os.Stdout, 8, 8, 8, ' ', 0)
		err := tmpl(w, UsageTemplate(), wrapper)
		if err != nil {
			c.PrintErrln(err)
		}
		w.Flush()
		return err
	}
}
