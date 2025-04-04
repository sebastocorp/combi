package utils

import (
	"crypto/md5"
	"encoding/hex"
	"os"
	"regexp"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// ExpandEnv TODO
func ExpandEnv(input []byte) []byte {
	re := regexp.MustCompile(`\${ENV:([A-Za-z_][A-Za-z0-9_]*)}\$`)
	result := re.ReplaceAllFunc(input, func(match []byte) []byte {
		key := re.FindSubmatch(match)[1]
		if value, exists := os.LookupEnv(string(key)); exists {
			return []byte(value)
		}
		return match
	})

	return result
}

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

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
