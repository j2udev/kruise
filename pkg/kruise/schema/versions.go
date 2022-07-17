package schema

type (
	IVersionedConfig interface {
		GetVersion()
	}
)

var Versions = []string{
	"v1alpha1",
}
