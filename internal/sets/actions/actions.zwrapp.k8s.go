package actions

import (
	"bytes"
	"combi/internal/sets/credentials"
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/remotecommand"
)

type K8sActionT struct {
	credRef *credentials.KubeT
	on      string

	namespace string
	pod       string
	container string
	cmd       []string
}

type OptionsK8sT struct {
	Namespace string
	Pod       string
	Container string
}

func newK8sAction(ops OptionsT) (a *K8sActionT, err error) {
	a = &K8sActionT{
		on: ops.On,

		namespace: ops.K8s.Namespace,
		pod:       ops.K8s.Pod,
		container: ops.K8s.Container,
	}

	switch ops.CredRef.(type) {
	case *credentials.KubeT:
		a.credRef = ops.CredRef.(*credentials.KubeT)
	default:
		err = fmt.Errorf("wrong credential type in '%s' action, must be K8S", ops.Name)
		return a, err
	}

	return a, err
}

func (a *K8sActionT) getOn() string {
	return a.on
}

func (a *K8sActionT) exec() (r ActionResultT, err error) {
	req := a.credRef.Cli.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(a.pod).
		Namespace(a.namespace).
		SubResource("exec")

	scheme := runtime.NewScheme()
	if err := corev1.AddToScheme(scheme); err != nil {
		return r, err
	}

	parameterCodec := runtime.NewParameterCodec(scheme)
	req = req.VersionedParams(&corev1.PodExecOptions{
		Command:   a.cmd,
		Container: a.container,
		Stdin:     false,
		Stdout:    true,
		Stderr:    true,
		TTY:       false,
	}, parameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(a.credRef.Cfg, "POST", req.URL())
	if err != nil {
		return r, err
	}

	var stdout, stderr bytes.Buffer
	err = exec.StreamWithContext(context.Background(), remotecommand.StreamOptions{
		Stdin:  nil,
		Stdout: &stdout,
		Stderr: &stderr,
		Tty:    false,
	})
	r.Stdout = stdout.String()
	r.Stderr = stderr.String()

	return r, err
}
