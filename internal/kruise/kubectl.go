package kruise

import (
	"os/exec"

	"github.com/j2udev/kruise/internal/kruise/schema/latest"
	"github.com/spf13/pflag"
)

type (
	// KubectlDeployment encapsulates Helm objects like KubectlSecrets and
	// KubectlManifests for a given deployment
	KubectlDeployment latest.KubectlDeployment
	// KubectlSecret represents information about a Kubectl secret
	KubectlSecret latest.KubectlSecret
	// KubectlManifest represents information about a Kubectl manifest
	KubectlManifest latest.KubectlManifest
	// KubectlDeployments represents a slice of KubectlDeployment objects
	KubectlDeployments []KubectlDeployment
	// KubectlSecrets represents a slice of KubectlSecret objects
	KubectlSecrets []KubectlSecret
	// KubectlManifest represents a slice of KubectlManifest objects
	KubectlManifests []KubectlManifest
)

// Install is used to execute a Kubectl apply command
func (m KubectlManifest) Install(fs *pflag.FlagSet) {
	d, err := fs.GetBool("shallow-dry-run")
	Fatal(err)
	if !d {
		checkKubectl()
	}
	Debug(kubectlCreateNamespace(d, m.Namespace))
	Error(kubectlExecute(d, m.installArgs(fs)))
}

// Uninstall is used to execute a Kubectl delete command
func (m KubectlManifest) Uninstall(fs *pflag.FlagSet) {
	d, err := fs.GetBool("shallow-dry-run")
	Fatal(err)
	if !d {
		checkKubectl()
	}
	Warn(kubectlExecute(d, m.uninstallArgs(fs)))
}

// GetPriority is used to get the priority of the installer
func (m KubectlManifest) GetPriority() int {
	return m.Priority
}

// Install is used to execute a Kubectl create secret command
func (s KubectlSecret) Install(fs *pflag.FlagSet) {
	d, err := fs.GetBool("shallow-dry-run")
	Fatal(err)
	if !d {
		checkKubectl()
	}
	Debug(kubectlCreateNamespace(d, s.Namespace))
	// for now, just overwrite any existing secret
	s.Uninstall(fs)
	Error(kubectlExecute(d, s.installArgs(fs)))
}

// Install is used to execute a Kubectl delete secret command
func (s KubectlSecret) Uninstall(fs *pflag.FlagSet) {
	d, err := fs.GetBool("shallow-dry-run")
	Fatal(err)
	if !d {
		checkKubectl()
	}
	Debug(kubectlDeleteSecret(d, s.uninstallArgs(fs)))
}

// GetPriority is used to get the priority of the installer
func (s KubectlSecret) GetPriority() int {
	// For now, KubectlSecrets are just installed first
	return 0
}

// newKubectlDeployment is a helper function for dealing with the latest.KubectlDeployment
// to KubectlDeployment type definition
func newKubectlDeployment(dep latest.KubectlDeployment) KubectlDeployment {
	return KubectlDeployment(dep)
}

// newKubectlSecret is a helper function for dealing with the latest.KubectlSecret
// to KubectlSecret type definition
func newKubectlSecret(sec latest.KubectlSecret) KubectlSecret {
	return KubectlSecret(sec)
}

// newKubectlSecrets is a helper function for dealing with the latest.KubectlSecret
// to KubectlSecret type definition
func newKubectlSecrets(secs []latest.KubectlSecret) KubectlSecrets {
	var s KubectlSecrets
	for _, sec := range secs {
		s = append(s, newKubectlSecret(sec))
	}
	return s
}

// newKubectlManifest is a helper function for dealing with the latest.KubectlManifest
// to KubectlManifest type definition
func newKubectlManifest(man latest.KubectlManifest) KubectlManifest {
	return KubectlManifest(man)
}

// newKubectlManifests is a helper function for dealing with the latest.KubectlManifest
// to KubectlManifest type definition
func newKubectlManifests(mans []latest.KubectlManifest) KubectlManifests {
	var m KubectlManifests
	for _, man := range mans {
		m = append(m, newKubectlManifest(man))
	}
	return m
}

// getKubectlSecrets is a helper function for grabbing the KubectlSecrets
// from a KubectlDeployment
func (d KubectlDeployment) getKubectlSecrets() KubectlSecrets {
	return newKubectlSecrets(d.Secrets)
}

// getKubectlManifests is a helper function for grabbing the KubectlManifests
// from a KubectlDeployment
func (d KubectlDeployment) getKubectlManifests() KubectlManifests {
	return newKubectlManifests(d.Manifests)
}

// installArgs is used to build Kubectl apply CLI args given a FlagSet
func (m KubectlManifest) installArgs(fs *pflag.FlagSet) []string {
	args := []string{"apply", "--namespace", m.Namespace}
	for _, p := range m.Paths {
		args = append(args, "-f", p)
	}
	return args
}

// uninstallArgs is used to build Kubectl delete CLI args given a FlagSet
func (m KubectlManifest) uninstallArgs(fs *pflag.FlagSet) []string {
	args := []string{"delete", "--namespace", m.Namespace}
	for _, p := range m.Paths {
		args = append(args, "-f", p)
	}
	return args
}

// installArgs is used to build Kubectl create secret CLI args given a FlagSet
func (s KubectlSecret) installArgs(fs *pflag.FlagSet) []string {
	sdr, err := fs.GetBool("shallow-dry-run")
	Fatal(err)
	args := []string{"create", "secret"}
	if s.Namespace != "" {
		args = append(args, "--namespace", s.Namespace)
	}
	switch s.Type {
	case "docker-registry":
		u := "***"
		p := []byte("***")
		args = append(args, "docker-registry", s.Name, "--docker-server", s.Registry)
		if !sdr {
			var up, pp string
			up = "Please enter your username for the " + s.Registry + " container registry"
			pp = "Please enter your password for the " + s.Registry + " container registry"
			un, pw, err := credentialPrompt(up, pp)
			Fatal(err)
			u = un
			p = []byte(pw)
		}
		args = append(args,
			"--docker-username", u,
			"--docker-password", string(p))
	default:
		Logger.Fatalf("kubectl secret type: %v not supported", s.Type)
	}
	return args
}

// uninstallArgs is used to build Kubectl delete secret CLI args given a FlagSet
func (s KubectlSecret) uninstallArgs(fs *pflag.FlagSet) []string {
	args := []string{"delete", "secret", s.Name}

	if s.Namespace != "" {
		args = append(args, "--namespace", s.Namespace)
	}
	return args
}

// kubectlCreateNamespace is used to execute a kubectl create namespace command
func kubectlCreateNamespace(dry bool, n string) error {
	return NewCmd("kubectl").
		WithArgs([]string{"create", "namespace", n}).
		WithDryRun(dry).
		WithNoStdOut().
		Build().
		Execute()
}

// kubectlDeleteSecret is used to execute a kubectl delete secret command
func kubectlDeleteSecret(dry bool, args []string) error {
	return NewCmd("kubectl").
		WithArgs(args).
		WithDryRun(dry).
		WithNoStdOut().
		Build().
		Execute()
}

// kubectlExecute is a helper function for executing a Kubectl command given a set of
// args; it will print the command instead of executing it if dry is true
func kubectlExecute(dry bool, args []string) error {
	return NewCmd("kubectl").
		WithArgs(args).
		WithDryRun(dry).
		Build().
		Execute()
}

// checkKubectl is used to determine if the user has the Kubectl CLI installed
func checkKubectl() {
	err := exec.Command("kubectl").Run()
	Fatalf(err, "Kubectl does not appear to be installed")
}
