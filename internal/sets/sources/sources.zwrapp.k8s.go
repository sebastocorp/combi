package sources

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"combi/internal/utils"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type K8sSourceT struct {
	name    string
	tmpPath string

	kube kubeT
}

type kubeT struct {
	ctx       context.Context
	client    *kubernetes.Clientset
	kind      string
	namespace string
	name      string
	key       string
}

type OptionsK8sT struct {
	InCluster      bool
	ConfigFilepath string
	MasterUrl      string
	Kind           string
	Namespace      string
	Name           string
	Key            string
}

func NewK8sSource(ops OptionsT) (s *K8sSourceT, err error) {
	s = &K8sSourceT{
		name:    ops.Name,
		tmpPath: ops.Path,

		kube: kubeT{
			ctx:       context.Background(),
			kind:      ops.K8s.Kind,
			namespace: ops.K8s.Namespace,
			name:      ops.K8s.Name,
			key:       ops.K8s.Key,
		},
	}

	var config *rest.Config
	if ops.K8s.InCluster {
		config, err = rest.InClusterConfig()
		if err != nil {
			return s, err
		}
	} else {
		config, err = clientcmd.BuildConfigFromFlags(ops.K8s.MasterUrl, ops.K8s.ConfigFilepath)
		if err != nil {
			return s, err
		}
	}

	s.kube.client, err = kubernetes.NewForConfig(config)

	return s, err
}

func (s *K8sSourceT) Name() string {
	return s.name
}

func (s *K8sSourceT) sync() (updated bool, err error) {
	srcBytes := []byte{}
	switch s.kube.kind {
	case "ConfigMap":
		{
			res, err := s.kube.client.CoreV1().ConfigMaps(s.kube.namespace).Get(s.kube.ctx, s.kube.name, v1.GetOptions{})
			if err != nil {
				return updated, err
			}

			configStr, ok := res.Data[s.kube.key]
			if !ok {
				err = fmt.Errorf("key '%s' does not exist in '%s' ConfigMap source", s.kube.key, s.kube.name)
				return updated, err
			}
			srcBytes = []byte(configStr)
		}
	case "Secret":
		{
			res, err := s.kube.client.CoreV1().Secrets(s.kube.namespace).Get(s.kube.ctx, s.kube.name, v1.GetOptions{})
			if err != nil {
				return updated, err
			}

			configStr, ok := res.StringData[s.kube.key]
			if !ok {
				err = fmt.Errorf("key '%s' does not exist in '%s' Secret source", s.kube.key, s.kube.name)
				return updated, err
			}
			srcBytes = []byte(configStr)
		}
	}

	storConfig := filepath.Join(s.tmpPath, s.kube.key)
	storBytes, err := os.ReadFile(storConfig)
	if err != nil {
		if os.IsNotExist(err) {
			updated = true
			err = os.WriteFile(storConfig, srcBytes, 0777)
			if err != nil {
				return updated, err
			}
		}
		return updated, err
	}

	if !reflect.DeepEqual(srcBytes, storBytes) {
		updated = true
		err = os.WriteFile(storConfig, srcBytes, 0777)
		if err != nil {
			return updated, err
		}
	}

	return updated, err
}

func (s *K8sSourceT) get() (conf []byte, err error) {
	storConfig := filepath.Join(s.tmpPath, s.kube.key)
	if conf, err = os.ReadFile(storConfig); err != nil {
		return conf, err
	}

	conf = utils.ExpandEnv(conf)

	return conf, err
}
