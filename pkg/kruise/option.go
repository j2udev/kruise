package kruise

// Option struct used to unmarshal kruise configuration and facilitate wrapping
// cobra commands
type Option struct {
	Arguments   string
	Description string
}

// GetOptions gets options associated with a config key argument. Valid keys
// are `deploy` and `delete`
func GetOptions(key string) []Option {
	var options []Option
	for _, dep := range GetHelmDeployments(key) {
		options = append(options, dep.Option)
	}
	return options
}
