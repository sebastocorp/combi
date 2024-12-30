package kube

import (
	"context"
	"fmt"
	"reflect"

	"combi/api/v1alpha2"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type KubeSourceT struct {
	context   context.Context
	client    *kubernetes.Clientset
	kind      string
	namespace string
	name      string
	key       string

	StoredConfig []byte
}

func (s *KubeSourceT) Init(source v1alpha2.SourceT) (err error) {
	if source.Kubernetes.Kind != "ConfigMap" && source.Kubernetes.Kind != "Secret" {
		err = fmt.Errorf("unsuported kind '%s' in kubernetes source", source.Kubernetes.Kind)
		return err
	}
	s.kind = source.Kubernetes.Kind
	s.namespace = source.Kubernetes.Namespace
	s.name = source.Kubernetes.Name
	s.key = source.Kubernetes.Key

	s.client, err = newClient("")
	if err != nil {
		return err
	}

	s.context = context.Background()

	return err
}

func (s *KubeSourceT) GetConfig() (config []byte, updated bool, err error) {
	configStr := ""
	ok := false
	if s.kind == "ConfigMap" {
		res, err := s.client.CoreV1().ConfigMaps(s.namespace).Get(s.context, s.name, v1.GetOptions{})
		if err != nil {
			return config, updated, err
		}

		if configStr, ok = res.Data[s.key]; !ok {
			err = fmt.Errorf("key '%s' does not exist in '%s' ConfigMap source", s.key, s.name)
			return config, updated, err
		}
	} else {
		res, err := s.client.CoreV1().Secrets(s.namespace).Get(s.context, s.name, v1.GetOptions{})
		if err != nil {
			return config, updated, err
		}

		if configStr, ok = res.StringData[s.key]; !ok {
			err = fmt.Errorf("key '%s' does not exist in '%s' Secret source", s.key, s.name)
			return config, updated, err
		}
	}
	config = []byte(configStr)

	if !reflect.DeepEqual(s.StoredConfig, config) {
		updated = true
	}

	return config, updated, err
}
