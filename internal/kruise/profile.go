package kruise

import (
	"github.com/j2udev/kruise/internal/schema/latest"
)

type (
	// Profile represents the arguments and description for a CLI profile
	Profile latest.Profile
	//Profiles represents a slice of Profile objects
	Profiles []Profile
)

// newProfile is a helper function used to create a new Profile object from a
// slice of arguments and a description
func newProfile(prof latest.Profile) Profile {
	return Profile(prof)
}

// getValidArgs is used to get arguments from a slice of Options
func (p Profiles) getValidArgs() []string {
	var valid []string
	for _, prof := range p {
		valid = append(valid, prof.Name)
		valid = append(valid, prof.Aliases...)
	}
	return valid
}
