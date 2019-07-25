package deploy

import "k8s.io/helm/pkg/helm"

type Deploy struct {
	client *helm.Client
}

func New() (*Deploy, error) {
	cli, err := NewClient()
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
