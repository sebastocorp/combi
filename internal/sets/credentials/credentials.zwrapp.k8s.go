package credentials

import (
	"context"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type KubeT struct {
	Ctx context.Context
	Cfg *rest.Config
	Cli *kubernetes.Clientset
}

type OptionsKubeT struct {
	InCluster      bool
	KubeconfigPath string
	MasterUrl      string
}

func NewKube(ops OptionsKubeT) (kc *KubeT, err error) {
	kc = &KubeT{
		Ctx: context.Background(),
	}

	if ops.InCluster {
		kc.Cfg, err = rest.InClusterConfig()
		if err != nil {
			return kc, err
		}
	} else {
		kc.Cfg, err = clientcmd.BuildConfigFromFlags(ops.MasterUrl, ops.KubeconfigPath)
		if err != nil {
			return kc, err
		}
	}

	kc.Cli, err = kubernetes.NewForConfig(kc.Cfg)

	return kc, err
}
