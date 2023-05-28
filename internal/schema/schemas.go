package schema

import (
	"github.com/j2udev/kruise/internal/schema/latest"
	"github.com/j2udev/kruise/internal/schema/v1alpha1"
	"github.com/j2udev/kruise/internal/schema/v1alpha2"
	"github.com/j2udev/kruise/internal/schema/version"
	"github.com/thoas/go-funk"
)

type (
	Versions []version.Version
)

var (
	SchemaVersionsV1 = Versions{
		{APIVersion: latest.Version, Factory: latest.NewKruiseConfig},
		{APIVersion: v1alpha2.Version, Factory: v1alpha2.NewKruiseConfig},
		{APIVersion: v1alpha1.Version, Factory: v1alpha1.NewKruiseConfig},
	}
)

func ConfigFactory(apiVersion string) version.IVersionedConfig {
	version := funk.Find(SchemaVersionsV1, func(v version.Version) bool {
		return v.APIVersion == apiVersion
	}).(version.Version)
	cfg := version.Factory()
	return cfg
}
