package kruise

import (
	"strings"

	"github.com/j2udevelopment/kruise/pkg/kruise/schema/latest"
)

type (
	// Profile represents the arguments and description for a CLI profile
	Profile struct {
		latest.Profile
		Arguments string
	}
	//Profiles represents a slice of Profile objects
	Profiles []Profile
)

// newProfile is a helper function used to create a new Profile object from a
// slice of arguments and a description
func newProfile(name string, prof latest.Profile) Profile {
	return Profile{
		prof,
		name + ", " + strings.Join(prof.Aliases, ", "),
	}
}
