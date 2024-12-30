package kube

import (
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func newClient(connectionMode string) (client *kubernetes.Clientset, err error) {
	// Create configuration to connect from inside the cluster using Kubernetes mechanisms
	config, err := rest.InClusterConfig()

	// Create configuration to connect from outside the cluster, using kubectl
	if connectionMode == "kubectl" {
		if home := homedir.HomeDir(); home != "" {
			config, err = clientcmd.BuildConfigFromFlags("", filepath.Join(home, ".kube", "config"))
		}
	}

	// Check configuration errors in both cases
	if err != nil {
		return client, err
	}

	// Construct the client
	client, err = kubernetes.NewForConfig(config)
	return client, err
}
