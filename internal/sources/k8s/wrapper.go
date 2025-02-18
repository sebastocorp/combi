package k8s

import (
	"combi/api/v1alpha4"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func newClient(ctx v1alpha4.SourceK8sContextConfigT) (client *kubernetes.Clientset, err error) {
	var config *rest.Config
	if ctx.InCluster {
		config, err = rest.InClusterConfig()
		if err != nil {
			return client, err
		}
	} else {
		config, err = clientcmd.BuildConfigFromFlags(ctx.MasterUrl, ctx.ConfigFilepath)
		if err != nil {
			return client, err
		}
	}

	// Construct the client
	client, err = kubernetes.NewForConfig(config)
	return client, err
}
