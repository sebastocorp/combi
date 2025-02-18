package k8s

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"combi/api/v1alpha4"
	"combi/internal/config"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type K8sSourceT struct {
	name       string
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

func NewK8sSource(srcConf v1alpha4.SourceConfigT, srcpath string) (s *K8sSourceT, err error) {
	s = &K8sSourceT{
		name:       srcConf.Name,
		srcConfig:  filepath.Join(srcpath, "sync", srcConf.K8s.Key),
		storConfig: filepath.Join(srcpath, srcConf.K8s.Key),

		kube: kubeT{
			ctx:       context.Background(),
			kind:      srcConf.K8s.Kind,
			namespace: srcConf.K8s.Namespace,
			name:      srcConf.K8s.Name,
			key:       srcConf.K8s.Key,
		},
	}

	s.kube.client, err = newClient(srcConf.K8s.Context)

	return s, err
}

func (s *K8sSourceT) GetName() string {
	return s.name
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
		err = os.WriteFile(s.storConfig, srcBytes, 0777)
		if err != nil {
			return updated, err
		}
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
