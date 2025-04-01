package credentials

import (
	"context"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type KubeT struct {
	ctx context.Context
	cli *kubernetes.Clientset
}

type OptionsKubeT struct {
	InCluster      bool
	ConfigFilepath string
	MasterUrl      string
}

func NewKube(ops OptionsKubeT) (kc *KubeT, err error) {
	kc = &KubeT{
		ctx: context.Background(),
	}

	var config *rest.Config
	if ops.InCluster {
		config, err = rest.InClusterConfig()
		if err != nil {
			return kc, err
		}
	} else {
		config, err = clientcmd.BuildConfigFromFlags(ops.MasterUrl, ops.ConfigFilepath)
		if err != nil {
			return kc, err
		}
	}

	kc.cli, err = kubernetes.NewForConfig(config)

	return kc, err
}
