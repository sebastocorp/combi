package k8s

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"combi/api/v1alpha3"
	"combi/internal/config"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type K8sSourceT struct {
	srcConfig  string
	storConfig string

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

func NewK8sSource(srcConf v1alpha3.SourceConfigT, srcpath string) (s *K8sSourceT, err error) {
	s = &K8sSourceT{
		srcConfig:  filepath.Join(srcpath, "sync", srcConf.Kubernetes.Key),
		storConfig: filepath.Join(srcpath, srcConf.Kubernetes.Key),
		kube: kubeT{
			ctx:       context.Background(),
			kind:      srcConf.Kubernetes.Kind,
			namespace: srcConf.Kubernetes.Namespace,
			name:      srcConf.Kubernetes.Name,
			key:       srcConf.Kubernetes.Key,
		},
	}

	s.kube.client, err = newClient("")

	return s, err
}

func (s *K8sSourceT) SyncConfig() (updated bool, err error) {
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

	storBytes, err := os.ReadFile(s.storConfig)
	if err != nil {
		if os.IsNotExist(err) {
			updated = true
			err = os.WriteFile(s.storConfig, srcBytes, 0777)
			if err != nil {
				return updated, err
			}
		}
		return updated, err
	}

	if !reflect.DeepEqual(srcBytes, storBytes) {
		updated = true
	}

	return updated, err
}

func (s *K8sSourceT) GetConfig() (conf []byte, err error) {
	if conf, err = os.ReadFile(s.storConfig); err != nil {
		return conf, err
	}

	conf = config.ExpandEnv(conf)

	return conf, err
}
