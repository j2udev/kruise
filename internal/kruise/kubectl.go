package kruise

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"os/exec"
	"strings"

	"github.com/j2udev/kruise/internal/schema/latest"
	"github.com/spf13/pflag"
)

type (
	// KubectlDeployment encapsulates Helm objects like KubectlGenericSecrets,
	// KubectlDockerRegistrySecrets and KubectlManifests for a given deployment
	KubectlDeployment latest.KubectlDeployment
	// KubectlManifest represents information about a Kubectl manifest
	KubectlManifest latest.KubectlManifest
	// KubectlGenericSecret represents information about a generic Kubernetes
	// secret
	// The Namespaces field is used to support creating the same secret across
	// multiple namespaces and only prompting the user once.
	KubectlGenericSecret struct {
		latest.KubectlGenericSecret
		Namespaces []string
	}
	// KubectlDockerRegistrySecret represents information about a docker-registry
	// Kubernetes secret
	// The Namespaces field is used to support creating the same secret across
	// multiple namespaces and only prompting the user once.
	KubectlDockerRegistrySecret struct {
		latest.KubectlDockerRegistrySecret
		Namespaces []string
	}
	// KubectlDeployments represents a slice of KubectlDeployment objects
	KubectlDeployments []KubectlDeployment
	// KubectlManifests represents a slice of KubectlManifest objects
	KubectlManifests []KubectlManifest
	// KubectlGenericSecrets represents a slice of KubectlGenericSecret objects
	KubectlGenericSecrets []KubectlGenericSecret
	// KubectlDockerRegistrySecrets represents a slice of
	// KubectlDockerRegistrySecret objects
	KubectlDockerRegistrySecrets []KubectlDockerRegistrySecret
)

// Install is used to execute a Kubectl apply command
func (m KubectlManifest) Install(fs *pflag.FlagSet) {
	d, err := fs.GetBool("dry-run")
	Fatal(err)
	if !d {
		checkKubectl()
	}
	Debug(kubectlCreateNamespace(d, m.Namespace))
	Error(kubectlExecute(d, m.installArgs(fs)))
}

// Install is used to execute a Kubectl create generic secret command
func (s KubectlGenericSecret) Install(fs *pflag.FlagSet) {
	d, err := fs.GetBool("dry-run")
	Fatal(err)
	if !d {
		checkKubectl()
	}
	for _, ns := range s.Namespaces {
		err = kubectlCreateNamespace(d, ns)
		if err != nil {
			Logger.Debug(err)
		}
	}
	// for now, just overwrite any existing secret
	s.Uninstall(fs)
	switch len(s.Namespaces) {
	case 0:
		fmt.Printf("Creating generic secret %s in the default namespace\n", s.Name)
	case 1:
		fmt.Printf("Creating generic secret %s in the %s namespace\n", s.Name, s.Namespace)
	default:
		fmt.Printf("Creating generic secret %s in the %s namespaces\n", s.Name, s.Namespaces)
	}
	args := s.installArgs(fs)
	for _, a := range args {
		Logger.Debugf("%s %s", "kubectl", strings.Join(a, " "))
		err = kubectlExecute(d, a)
		if err != nil {
			Logger.Error(err)
		}
	}
}

// Install is used to execute a Kubectl create docker-registry secret command
func (s KubectlDockerRegistrySecret) Install(fs *pflag.FlagSet) {
	d, err := fs.GetBool("dry-run")
	Fatal(err)
	if !d {
		checkKubectl()
	}
	for _, ns := range s.Namespaces {
		err = kubectlCreateNamespace(d, ns)
		if err != nil {
			Logger.Debug(err)
		}
	}
	// for now, just overwrite any existing secret
	s.Uninstall(fs)
	switch len(s.Namespaces) {
	case 0:
		fmt.Printf("Creating docker-registry secret %s in the default namespace\n", s.Name)
	case 1:
		fmt.Printf("Creating docker-registry secret %s in the %s namespace\n", s.Name, s.Namespace)
	default:
		fmt.Printf("Creating docker-registry secret %s in the %s namespaces\n", s.Name, s.Namespaces)
	}
	args := s.installArgs(fs)
	for _, a := range args {
		Logger.Debugf("%s %s", "kubectl", strings.Join(a, " "))
		err = kubectlExecute(d, a)
		if err != nil {
			Logger.Error(err)
		}
	}
}

// Uninstall is used to execute a Kubectl delete command
func (m KubectlManifest) Uninstall(fs *pflag.FlagSet) {
	d, err := fs.GetBool("dry-run")
	Fatal(err)
	if !d {
		checkKubectl()
	}
	Warn(kubectlExecute(d, m.uninstallArgs(fs)))
}

// Uninstall is used to execute a Kubectl delete secret command
func (s KubectlGenericSecret) Uninstall(fs *pflag.FlagSet) {
	d, err := fs.GetBool("dry-run")
	Fatal(err)
	if !d {
		checkKubectl()
	}
	args := s.uninstallArgs(fs)
	for _, a := range args {
		Logger.Debugf("%s %s", "kubectl", strings.Join(a, " "))
		err = kubectlDeleteSecret(d, a)
		if err != nil {
			Logger.Debug(err)
		}
	}
}

// Uninstall is used to execute a Kubectl delete secret command
func (s KubectlDockerRegistrySecret) Uninstall(fs *pflag.FlagSet) {
	d, err := fs.GetBool("dry-run")
	Fatal(err)
	if !d {
		checkKubectl()
	}
	args := s.uninstallArgs(fs)
	for _, a := range args {
		Logger.Debugf("%s %s", "kubectl", strings.Join(a, " "))
		err = kubectlDeleteSecret(d, a)
		if err != nil {
			Logger.Debug(err)
		}
	}
}

// GetPriority is used to get the priority of the installer
func (m KubectlManifest) GetPriority() int {
	return m.Priority
}

// GetPriority is used to get the priority of the installer
func (s KubectlGenericSecret) GetPriority() int {
	// for now, kubectl secrets are just installed first
	return 0
}

// GetPriority is used to get the priority of the installer
func (s KubectlDockerRegistrySecret) GetPriority() int {
	// for now, kubectl secrets are just installed first
	return 0
}

// IsInit is used to determine whether the installer should be installed during
// initialization
func (m KubectlManifest) IsInit() bool {
	return m.Init
}

// IsInit is used to determine whether the installer should be installed during
// initialization
func (m KubectlGenericSecret) IsInit() bool {
	return m.Init
}

// IsInit is used to determine whether the installer should be installed during
// initialization
func (m KubectlDockerRegistrySecret) IsInit() bool {
	return m.Init
}

// newKubectlDeployment is a helper function for dealing with the
// latest.KubectlDeployment to KubectlDeployment type definition
func newKubectlDeployment(dep latest.KubectlDeployment) KubectlDeployment {
	return KubectlDeployment(dep)
}

// newKubectlManifest is a helper function for dealing with the
// latest.KubectlManifest to KubectlManifest type definition
func newKubectlManifest(man latest.KubectlManifest) KubectlManifest {
	return KubectlManifest(man)
}

// newKubectlGenericSecret is a helper function for dealing with the
// latest.KubectlGenericSecret to KubectlGenericSecret type definition
func newKubectlGenericSecret(sec latest.KubectlGenericSecret) KubectlGenericSecret {
	// return KubectlGenericSecret(sec)
	if sec.Namespace == "" {
		sec.Namespace = "default"
	}
	return KubectlGenericSecret{sec, []string{sec.Namespace}}
}

// newKubectlDockerRegistrySecret is a helper function for dealing with the
// latest.KubectlDockerRegistrySecret to KubectlDockerRegistrySecret type
// definition
func newKubectlDockerRegistrySecret(sec latest.KubectlDockerRegistrySecret) KubectlDockerRegistrySecret {
	// return KubectlDockerRegistrySecret(sec)
	if sec.Namespace == "" {
		sec.Namespace = "default"
	}
	return KubectlDockerRegistrySecret{sec, []string{sec.Namespace}}
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

// newKubectlGenericSecrets is a helper function for dealing with the
// latest.KubectlGenericSecrets to KubectlGenericSecrets type definition
func newKubectlGenericSecrets(secs []latest.KubectlGenericSecret) KubectlGenericSecrets {
	var s KubectlGenericSecrets
	for _, sec := range secs {
		s = append(s, newKubectlGenericSecret(sec))
	}
	return s
}

// newKubectlDockerRegistrySecrets is a helper function for dealing with the
// latest.KubectlDockerRegistrySecrets to KubectlDockerRegistrySecrets type
// definition
func newKubectlDockerRegistrySecrets(secs []latest.KubectlDockerRegistrySecret) KubectlDockerRegistrySecrets {
	var s KubectlDockerRegistrySecrets
	for _, sec := range secs {
		s = append(s, newKubectlDockerRegistrySecret(sec))
	}
	return s
}

// getKubectlManifests is a helper function for grabbing the KubectlManifests
// from a KubectlDeployment
func (d KubectlDeployment) getKubectlManifests() KubectlManifests {
	return newKubectlManifests(d.Manifests)
}

// getKubectlGenericSecrets is a helper function for grabbing the
// KubectlGenericSecrets from a KubectlDeployment
func (d KubectlDeployment) getKubectlGenericSecrets() KubectlGenericSecrets {
	return newKubectlGenericSecrets(d.Secrets.Generic)
}

// getKubectlDockerRegistrySecrets is a helper function for grabbing the
// KubectlDockerRegistrySecrets from a KubectlDeployment
func (d KubectlDeployment) getKubectlDockerRegistrySecrets() KubectlDockerRegistrySecrets {
	return newKubectlDockerRegistrySecrets(d.Secrets.DockerRegistry)
}

// installArgs is used to build Kubectl apply CLI args given a FlagSet
func (m KubectlManifest) installArgs(fs *pflag.FlagSet) []string {
	args := []string{"apply", "--namespace", m.Namespace}
	for _, p := range m.Paths {
		args = append(args, "-f", p)
	}
	return args
}

// installArgs is used to build Kubectl create generic secret CLI args given a
// FlagSet
func (s KubectlGenericSecret) installArgs(fs *pflag.FlagSet) [][]string {
	d, err := fs.GetBool("dry-run")
	Fatal(err)
	var iargs [][]string
	var largs []string
	v := "***"
	for _, l := range s.Literal {
		if l.Val != "" {
			largs = append(largs, "--from-literal", fmt.Sprintf("%s=%s", l.Key, l.Val))
		} else {
			if !d {
				v = sensitiveInputPrompt(fmt.Sprintf("Please enter a value for key: %s", l.Key))
			}
			largs = append(largs, "--from-literal", fmt.Sprintf("%s=%s", l.Key, v))
		}
	}
	for _, ns := range s.Namespaces {
		args := []string{"create", "secret", "generic", s.Name}
		if ns != "" && ns != "default" {
			args = append(args, "--namespace", ns)
		}
		args = append(args, largs...)
		iargs = append(iargs, args)
	}
	return iargs
}

// installArgs is used to build Kubectl create docker-registry secret CLI args
// given a FlagSet
func (s KubectlDockerRegistrySecret) installArgs(fs *pflag.FlagSet) [][]string {
	d, err := fs.GetBool("dry-run")
	Fatal(err)
	var iargs [][]string
	var dargs []string
	u := "***"
	p := "***"
	dargs = append(dargs, "--docker-server", s.Registry)
	if !d {
		u = normalInputPrompt("Please enter a username")
		p = sensitiveInputPrompt("Please enter a password")
	}
	dargs = append(dargs, "--docker-username", u, "--docker-password", string(p))
	for _, ns := range s.Namespaces {
		args := []string{"create", "secret", "docker-registry", s.Name}
		if ns != "" && ns != "default" {
			args = append(args, "--namespace", ns)
		}
		args = append(args, dargs...)
		iargs = append(iargs, args)
	}
	return iargs
}

// uninstallArgs is used to build Kubectl delete CLI args given a FlagSet
func (m KubectlManifest) uninstallArgs(fs *pflag.FlagSet) []string {
	args := []string{"delete", "--namespace", m.Namespace}
	for _, p := range m.Paths {
		args = append(args, "-f", p)
	}
	return args
}

// uninstallArgs is used to build Kubectl delete secret CLI args given a FlagSet
func (s KubectlGenericSecret) uninstallArgs(fs *pflag.FlagSet) [][]string {
	var uargs [][]string
	for _, ns := range s.Namespaces {
		args := []string{"delete", "secret", s.Name}
		if ns != "" && ns != "default" {
			args = append(args, "--namespace", ns)
		}
		uargs = append(uargs, args)
	}
	return uargs
}

// uninstallArgs is used to build Kubectl delete secret CLI args given a FlagSet
func (s KubectlDockerRegistrySecret) uninstallArgs(fs *pflag.FlagSet) [][]string {
	var uargs [][]string
	for _, ns := range s.Namespaces {
		args := []string{"delete", "secret", s.Name}
		if ns != "" && ns != "default" {
			args = append(args, "--namespace", ns)
		}
		uargs = append(uargs, args)
	}
	return uargs
}

// hash is used to facilitate storing KubectlManifests in a map
func (m *KubectlManifest) hash() string {
	h := sha1.New()
	for _, p := range m.Paths {
		h.Write([]byte(p))
	}
	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}

// hash is used to facilitate storing KubectlGenericSecrets in a map
func (s *KubectlGenericSecret) hash() string {
	h := sha1.New()
	h.Write([]byte(s.Name))
	for _, l := range s.Literal {
		h.Write([]byte(l.Key))
		h.Write([]byte(l.Val))
	}
	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}

// hash is used to facilitate storing KubectlDockerRegistrySecrets in a map
func (s *KubectlDockerRegistrySecret) hash() string {
	h := sha1.New()
	h.Write([]byte(s.Name))
	h.Write([]byte(s.Registry))
	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}

// kubectlCreateNamespace is used to execute a kubectl create namespace command
// hides unnecessary output
func kubectlCreateNamespace(dry bool, n string) error {
	return NewCmd("kubectl").
		WithArgs([]string{"create", "namespace", n}).
		WithDryRun(dry).
		WithNoStdOut().
		Build().
		Execute()
}

// kubectlDeleteSecret is used to execute a kubectl delete secret command
// hides unnecessary output
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
