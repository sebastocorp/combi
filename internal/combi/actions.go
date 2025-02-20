package combi

import (
	"bytes"
	"context"
	"os/exec"

	"combi/api/v1alpha4"
	"combi/internal/utils"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

type ActionT struct {
	Name string `json:"name"`
	On   string `json:"on"`

	cmd []string   `json:"-"`
	k8s ActionK8sT `json:"-"`
}

type ActionK8sT struct {
	cfg    *rest.Config
	client *kubernetes.Clientset

	namespace string
	pod       string
	container string

	cmd []string
}

func NewAction(action v1alpha4.ActionConfigT) (a ActionT, err error) {
	a = ActionT{
		Name: action.Name,
		On:   action.On,

		cmd: action.Command,
		k8s: ActionK8sT{
			namespace: action.K8s.Namespace,
			pod:       action.K8s.Pod,
			container: action.K8s.Container,
			cmd:       action.K8s.Command,
		},
	}

	a.k8s.cfg, err = utils.GetK8sConfig(action.K8s.Context.InCluster, action.K8s.Context.ConfigFilepath, action.K8s.Context.MasterUrl)
	if err != nil {
		return a, err
	}

	a.k8s.client, err = utils.GetK8sClient(a.k8s.cfg)

	return a, err
}

func (a *ActionT) Exec() (stdoutBytes []byte, stderrBytes []byte, err error) {
	if len(a.cmd) > 0 {
		var stdout, stderr bytes.Buffer
		cmd := exec.Command(a.cmd[0], a.cmd[1:]...)
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		err = cmd.Run()
		if err != nil {
			return stdoutBytes, stderrBytes, err
		}

		stdoutBytes = stdout.Bytes()
		stderrBytes = stderr.Bytes()
	}

	if len(a.k8s.cmd) > 0 {
		stdoutBytes, stderrBytes, err = a.execK8sCommand()
	}

	return stdoutBytes, stderrBytes, err
}

func (a *ActionT) execK8sCommand() (stdoutBytes []byte, stderrBytes []byte, err error) {
	req := a.k8s.client.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(a.k8s.pod).
		Namespace(a.k8s.namespace).
		SubResource("exec")

	scheme := runtime.NewScheme()
	if err := corev1.AddToScheme(scheme); err != nil {
		return stdoutBytes, stderrBytes, err
	}

	parameterCodec := runtime.NewParameterCodec(scheme)
	req = req.VersionedParams(&corev1.PodExecOptions{
		Command:   a.k8s.cmd,
		Container: a.k8s.container,
		Stdin:     false,
		Stdout:    true,
		Stderr:    true,
		TTY:       false,
	}, parameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(a.k8s.cfg, "POST", req.URL())
	if err != nil {
		return stdoutBytes, stderrBytes, err
	}

	var stdout, stderr bytes.Buffer
	err = exec.StreamWithContext(context.Background(), remotecommand.StreamOptions{
		Stdin:  nil,
		Stdout: &stdout,
		Stderr: &stderr,
		Tty:    false,
	})
	if err != nil {
		return stdoutBytes, stderrBytes, err
	}
	stdoutBytes = stdout.Bytes()
	stderrBytes = stderr.Bytes()

	return stdoutBytes, stderrBytes, err
}
