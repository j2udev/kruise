package kruise

type (
	IInstaller interface {
		Install(dryRun bool) error
		Uninstall(dryRun bool) error
	}
)
