package kruise

import (
	"strings"

	"github.com/j2udev/kruise/internal/kruise/schema/latest"
)

type (
	// Profile represents the arguments and description for a CLI profile
	Profile struct {
		latest.Profile
		Args string
		Desc string
	}
	//Profiles represents a slice of Profile objects
	Profiles []Profile
)

// newProfile is a helper function used to create a new Profile object from a
// slice of arguments and a description
func newProfile(name string, prof latest.Profile) Profile {
	var aliases []string
	aliases = append(aliases, name)
	aliases = append(aliases, prof.Aliases...)
	return Profile{
		prof,
		strings.Join(aliases, ","),
		prof.Description.Deploy,
	}
}

// getValidArgs is used to get arguments from a slice of Options
func (p Profiles) getValidArgs() []string {
	var valid []string
	for _, prof := range p {
		valid = append(valid, strings.Split(prof.Args, ", ")...)
	}
	return valid
}
