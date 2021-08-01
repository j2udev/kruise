package kruise

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"strings"
	"text/tabwriter"
	"unicode"

	"github.com/spf13/cobra"
)

var templateFuncs = template.FuncMap{
	"trim":                    strings.TrimSpace,
	"trimRightSpace":          trimRightSpace,
	"trimTrailingWhitespaces": trimRightSpace,
	"rpad":                    rpad,
}

// trimRightSpace trims any trailing whitespace
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

// UsageFunc overrides the default UsageFunc used by Cobra to facilitate showing command options
func UsageFunc(kmd Kommand) (f func(*cobra.Command) error) {
	return func(c *cobra.Command) error {
		w := tabwriter.NewWriter(os.Stdout, 8, 8, 8, ' ', 0)
		err := tmpl(w, UsageTemplate(), kmd)
		if err != nil {
			c.PrintErrln(err)
		}
		w.Flush()
		return err
	}
}

// HelpFunc overrides the default HelpFunc used by Cobra to facilitate showing command options
func HelpFunc(kmd Kommand) func(*cobra.Command, []string) {
	return func(c *cobra.Command, s []string) {
		w := tabwriter.NewWriter(os.Stdout, 8, 8, 8, ' ', 0)
		err := tmpl(w, UsageTemplate(), kmd)
		if err != nil {
			c.PrintErrln(err)
		}
		w.Flush()
	}
}
