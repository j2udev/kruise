package kruise

// UsageTemplate is used to override the cobra UsageTemplate to facilitate
// options
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
