package utils

import (
	"crypto/md5"
	"encoding/hex"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func GenHashString(args ...string) (h string, err error) {
	str := ""
	for _, sv := range args {
		str += sv
	}

	md5Hash := md5.New()
	_, err = md5Hash.Write([]byte(str))
	if err != nil {
		return h, err
	}

	h = hex.EncodeToString(md5Hash.Sum(nil))

	return h, err
}

// GetK8sConfig TODO
func GetK8sConfig(inCluster bool, configFilepath, masterUrl string) (cfg *rest.Config, err error) {
	if inCluster {
		cfg, err = rest.InClusterConfig()
	} else {
		cfg, err = clientcmd.BuildConfigFromFlags(masterUrl, configFilepath)
	}

	return cfg, err
}

// GetK8sClient TODO
func GetK8sClient(cfg *rest.Config) (*kubernetes.Clientset, error) {
	return kubernetes.NewForConfig(cfg)
}
