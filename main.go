package main

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/helm/pkg/kube"
)

func main() {
	kubeconfig := os.Getenv("KUBECONFIG")
	if len(kubeconfig) == 0 {
		kubeconfig = "~/.kube/config"
	}

	config, client, err := getKubeClient("", kubeconfig)
	if err != nil {
		log.Fatal(err)
	}
	log.Info(config)
	log.Info(client)
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
