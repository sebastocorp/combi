package sources

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"combi/internal/sets/credentials"
	"combi/internal/utils"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type K8sSourceT struct {
	name    string
	encType string
	workDir string
	credRef *credentials.KubeT

	resKind      string
	resNamespace string
	resName      string
	resKey       string
}

type OptionsK8sT struct {
	Kind      string
	Namespace string
	Name      string
	Key       string
}

func NewK8sSource(ops OptionsT) (s *K8sSourceT, err error) {
	s = &K8sSourceT{
		name:    ops.Name,
		encType: ops.EncType,
		workDir: ops.WorkDir,
		credRef: ops.CredRef.(*credentials.KubeT),

		resKind:      ops.K8s.Kind,
		resNamespace: ops.K8s.Namespace,
		resName:      ops.K8s.Name,
		resKey:       ops.K8s.Key,
	}

	switch ops.CredRef.(type) {
	case *credentials.KubeT:
		s.credRef = ops.CredRef.(*credentials.KubeT)
	default:
		err = fmt.Errorf("wrong credential type in '%s' source, must be K8S", ops.Name)
		return s, err
	}

	return s, err
}

func (s *K8sSourceT) getName() string {
	return s.name
}

func (s *K8sSourceT) getData() (srcd SourceDataT, err error) {
	srcd.Name = s.name
	srcd.SrcType = TypeK8S
	srcd.EncType = s.encType

	storConfig := filepath.Join(s.workDir, s.resKey)
	if srcd.Data, err = os.ReadFile(storConfig); err != nil {
		return srcd, err
	}
	srcd.Data = utils.ExpandEnv(srcd.Data)

	return srcd, err
}

func (s *K8sSourceT) sync() (updated bool, err error) {
	srcBytes := []byte{}
	switch s.resKind {
	case "ConfigMap":
		{
			res, err := s.credRef.Cli.CoreV1().ConfigMaps(s.resNamespace).Get(s.credRef.Ctx, s.resName, v1.GetOptions{})
			if err != nil {
				return updated, err
			}

			configStr, ok := res.Data[s.resKey]
			if !ok {
				err = fmt.Errorf("key '%s' does not exist in '%s' ConfigMap source", s.resKey, s.resName)
				return updated, err
			}
			srcBytes = []byte(configStr)
		}
	case "Secret":
		{
			res, err := s.credRef.Cli.CoreV1().Secrets(s.resNamespace).Get(s.credRef.Ctx, s.resName, v1.GetOptions{})
			if err != nil {
				return updated, err
			}

			configStr, ok := res.StringData[s.resKey]
			if !ok {
				err = fmt.Errorf("key '%s' does not exist in '%s' Secret source", s.resKey, s.resName)
				return updated, err
			}
			srcBytes = []byte(configStr)
		}
	}

	storConfig := filepath.Join(s.workDir, s.resKey)
	storBytes, err := os.ReadFile(storConfig)
	if err != nil {
		if os.IsNotExist(err) {
			updated = true
			err = os.WriteFile(storConfig, srcBytes, utils.FileModePerm)
			if err != nil {
				return updated, err
			}
		}
		return updated, err
	}

	if !reflect.DeepEqual(storBytes, srcBytes) {
		updated = true
		err = os.WriteFile(storConfig, srcBytes, utils.FileModePerm)
		if err != nil {
			return updated, err
		}
	}

	return updated, err
}
