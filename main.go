package main

import (
	"fmt"
	"log"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/helm/pkg/kube"
)

func main() {
	// config, client, err := getKubeClient("", "$HOME/.kube/config-internal-stg")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.Print(config)
	// log.Print(client)
	log.Print("echo")
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
