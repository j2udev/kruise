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
