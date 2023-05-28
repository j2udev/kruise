package version

type (
	IVersionedConfig interface {
		GetVersion() string
	}

	Version struct {
		APIVersion string
		Factory    func() IVersionedConfig
	}
)
