package deploy

import (
	"fmt"
	"os"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/helm/pkg/helm"
	"k8s.io/helm/pkg/helm/portforwarder"
	"k8s.io/helm/pkg/kube"
)

func NewClient() (*helm.Client, error) {
	kubeconfig := os.Getenv("KUBECONFIG")
	if len(kubeconfig) == 0 {
		kubeconfig = "~/.kube/config"
	}

	config, client, err := getKubeClient("", kubeconfig)
	if err != nil {
		return nil, err
	}
	tillerTunnel, err := portforwarder.New("kube-system", client, config)
	if err != nil {
		return nil, err
	}
	tillerHost := fmt.Sprintf("127.0.0.1:%d", tillerTunnel.Local)
	options := []helm.Option{helm.Host(tillerHost), helm.ConnectTimeout(300)}
	return helm.NewClient(options...), nil
}

func configForContext(context string, kubeconfig string) (*rest.Config, error) {
	config, err := kube.GetConfig(context, kubeconfig).ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("could not get Kubernetes config for context %q: %s", context, err)
	}
	return config, nil
}

func getKubeClient(context string, kubeconfig string) (*rest.Config, kubernetes.Interface, error) {
	config, err := configForContext(context, kubeconfig)
	if err != nil {
		return nil, nil, err
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, fmt.Errorf("could not get Kubernetes client: %s", err)
	}
	return config, client, nil
}
