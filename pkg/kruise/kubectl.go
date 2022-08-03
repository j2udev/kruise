package kruise

import (
	"fmt"
	"os/exec"

	"github.com/j2udevelopment/kruise/pkg/kruise/schema/latest"
	"github.com/spf13/pflag"
)

type (
	KubectlDeployment  latest.KubectlDeployment
	KubectlSecret      latest.KubectlSecret
	KubectlManifest    latest.KubectlManifest
	KubectlDeployments []KubectlDeployment
	KubectlSecrets     []KubectlSecret
	KubectlManifests   []KubectlManifest
)

// NewKubectlDeployment is a helper function for dealing with the latest.KubectlDeployment
// to KubectlDeployment type definition
func NewKubectlDeployment(dep latest.KubectlDeployment) KubectlDeployment {
	return KubectlDeployment(dep)
}

// NewKubectlDeployments is a helper function for dealing with the latest.KubectlDeployment
// to KubectlDeployment type definition
func NewKubectlDeployments(deps []latest.KubectlDeployment) KubectlDeployments {
	var d []KubectlDeployment
	for _, dep := range deps {
		d = append(d, NewKubectlDeployment(dep))
	}
	return d
}

// NewKubectlSecret is a helper function for dealing with the latest.KubectlSecret
// to KubectlSecret type definition
func NewKubectlSecret(sec latest.KubectlSecret) KubectlSecret {
	return KubectlSecret(sec)
}

// NewKubectlSecrets is a helper function for dealing with the latest.KubectlSecret
// to KubectlSecret type definition
func NewKubectlSecrets(secs []latest.KubectlSecret) KubectlSecrets {
	var s KubectlSecrets
	for _, sec := range secs {
		s = append(s, NewKubectlSecret(sec))
	}
	return s
}

// NewKubectlManifest is a helper function for dealing with the latest.KubectlManifest
// to KubectlManifest type definition
func NewKubectlManifest(man latest.KubectlManifest) KubectlManifest {
	return KubectlManifest(man)
}

// NewKubectlManifests is a helper function for dealing with the latest.KubectlManifest
// to KubectlManifest type definition
func NewKubectlManifests(mans []latest.KubectlManifest) KubectlManifests {
	var m KubectlManifests
	for _, man := range mans {
		m = append(m, NewKubectlManifest(man))
	}
	return m
}

func (d KubectlDeployment) GetKubectlSecrets() KubectlSecrets {
	return NewKubectlSecrets(d.Secrets)
}

func (d KubectlDeployment) GetKubectlManifests() KubectlManifests {
	return NewKubectlManifests(d.Manifests)
}

func (m KubectlManifest) Install(fs *pflag.FlagSet) {
	d, err := fs.GetBool("shallow-dry-run")
	Fatal(err)
	if !d {
		checkKubectl()
		// printStatus("info", "Kubectl applying "+strings.Join(m.Paths, ", ")+"...")
	}
	Debug(kubectlCreateNamespace(d, m.Namespace))
	Error(kubectlExecute(d, m.installArgs(fs)))
}

func (m KubectlManifest) Uninstall(fs *pflag.FlagSet) {
	d, err := fs.GetBool("shallow-dry-run")
	Fatal(err)
	if !d {
		checkKubectl()
		// printStatus("info", "Kubectl deleting "+strings.Join(m.Paths, ", ")+"...")
	}
	Warn(kubectlExecute(d, m.uninstallArgs(fs)))
}

func (m KubectlManifest) GetPriority() int {
	return m.Priority
}

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

func (s KubectlSecret) Uninstall(fs *pflag.FlagSet) {
	d, err := fs.GetBool("shallow-dry-run")
	Fatal(err)
	if !d {
		checkKubectl()
	}
	Debug(kubectlDeleteSecret(d, s.uninstallArgs(fs)))
}

func (m KubectlManifest) installArgs(fs *pflag.FlagSet) []string {
	args := []string{"apply", "--namespace", m.Namespace}
	for _, p := range m.Paths {
		args = append(args, "-f", p)
	}
	return args
}

func (m KubectlManifest) uninstallArgs(fs *pflag.FlagSet) []string {
	args := []string{"delete", "--namespace", m.Namespace}
	for _, p := range m.Paths {
		args = append(args, "-f", p)
	}
	return args
}

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
			up = fmt.Sprintf("Please enter your username for the %s container registry", s.Name)
			pp = fmt.Sprintf("Please enter your password for the %s container registry", s.Name)
			uname, pwd, err := CredentialPrompt(up, pp)
			Fatal(err)
			u = uname
			p = []byte(pwd)
		}
		args = append(args, "--docker-username", u, "--docker-password", string(p))
	default:
		Logger.Fatalf("kubectl secret type: %v not supported", s.Type)
	}
	return args
}

func (s KubectlSecret) uninstallArgs(fs *pflag.FlagSet) []string {
	args := []string{"delete", "secret", s.Name}

	if s.Namespace != "" {
		args = append(args, "--namespace", s.Namespace)
	}
	return args
}

func kubectlCreateNamespace(dry bool, n string) error {
	return NewCmd("kubectl").
		WithArgs([]string{"create", "namespace", n}).
		WithDryRun(dry).
		WithNoStdOut().
		Build().
		Execute()
}

func kubectlDeleteSecret(dry bool, args []string) error {
	return NewCmd("kubectl").
		WithArgs(args).
		WithDryRun(dry).
		WithNoStdOut().
		Build().
		Execute()
}

func kubectlExecute(dry bool, args []string) error {
	return NewCmd("kubectl").
		WithArgs(args).
		WithDryRun(dry).
		Build().
		Execute()
}

func checkKubectl() {
	err := exec.Command("kubectl").Run()
	Fatalf(err, "Kubectl does not appear to be installed")
}
