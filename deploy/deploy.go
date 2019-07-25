package deploy

import (
	"gopkg.in/yaml.v2"
	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/helm"
	"k8s.io/helm/pkg/kube"
	"k8s.io/helm/pkg/proto/hapi/release"
)

type Deploy struct {
	client      *helm.Client
	kubeContext string
	kubeConfig  string
}

func New(kubeContext, kubeConfig string) (*Deploy, error) {
	cli, err := NewClient(kubeContext, kubeConfig)
	if err != nil {
		return nil, err
	}
	c := &Deploy{
		client: cli,
	}
	return c, nil
}

func (d *Deploy) Version() (string, error) {
	version, err := d.client.GetVersion()
	if err != nil {
		return "", err
	}
	return version.Version.GetSemVer(), nil
}

func (d *Deploy) NewRelease(chartPath string) (*release.Release, error) {
	chartRequested, err := chartutil.Load(chartPath)
	if err != nil {
		return nil, err
	}

	namespace, _, err := kube.GetConfig(d.kubeContext, d.kubeConfig).Namespace()
	if err != nil {
		return nil, err
	}

	values, err := d.overrideValues()
	if err != nil {
		return nil, err
	}

	res, err := d.client.InstallReleaseFromChart(
		chartRequested,
		namespace,
		helm.ValueOverrides(values),
		helm.ReleaseName("akira"),
		helm.InstallDryRun(true),
		helm.InstallReuseName(false),
		helm.InstallDisableHooks(false),
		helm.InstallDisableCRDHook(false),
		helm.InstallSubNotes(false),
		helm.InstallTimeout(300),
		helm.InstallWait(true),
		helm.InstallDescription(""),
	)
	if err != nil {
		return nil, err
	}
	release := res.GetRelease()
	if release == nil {
		return nil, nil
	}
	return release, nil
}

func (d *Deploy) overrideValues() ([]byte, error) {
	base := map[string]interface{}{}

	return yaml.Marshal(base)
}
