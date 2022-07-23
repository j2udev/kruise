package schema

import (
	"github.com/j2udevelopment/kruise/pkg/kruise/schema/latest"
	"github.com/j2udevelopment/kruise/pkg/kruise/schema/version"
	"github.com/thoas/go-funk"
)

type (
	Versions []version.Version
)

var (
	SchemaVersionsV1 = Versions{
		{APIVersion: latest.Version, Factory: latest.NewKruiseConfig},
	}
)

func ConfigFactory(apiVersion string) version.IVersionedConfig {
	version := funk.Find(SchemaVersionsV1, func(v version.Version) bool {
		return v.APIVersion == apiVersion
	}).(version.Version)
	cfg := version.Factory()
	return cfg
}
