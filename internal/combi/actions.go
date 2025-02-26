package combi

import (
	"bytes"
	"context"
	"os/exec"
	"strings"

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

type ActionResultT struct {
	Cmd    string `json:"cmd"`
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
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

func (a *ActionT) Exec() (result ActionResultT, err error) {
	if len(a.cmd) > 0 {
		result.Cmd = strings.Join(a.cmd, " ")
		var stdout, stderr bytes.Buffer
		cmd := exec.Command(a.cmd[0], a.cmd[1:]...)
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		err = cmd.Run()
		if err != nil {
			return result, err
		}

		result.Stdout = stdout.String()
		result.Stderr = stderr.String()
	}

	if len(a.k8s.cmd) > 0 {
		result.Cmd = strings.Join(a.k8s.cmd, " ")
		result.Stdout, result.Stderr, err = a.execK8sCommand()
	}

	return result, err
}

func (a *ActionT) execK8sCommand() (stdoutStr, stderrStr string, err error) {
	req := a.k8s.client.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(a.k8s.pod).
		Namespace(a.k8s.namespace).
		SubResource("exec")

	scheme := runtime.NewScheme()
	if err := corev1.AddToScheme(scheme); err != nil {
		return stdoutStr, stderrStr, err
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
		return stdoutStr, stderrStr, err
	}

	var stdout, stderr bytes.Buffer
	err = exec.StreamWithContext(context.Background(), remotecommand.StreamOptions{
		Stdin:  nil,
		Stdout: &stdout,
		Stderr: &stderr,
		Tty:    false,
	})
	stdoutStr = stdout.String()
	stderrStr = stderr.String()

	return stdoutStr, stderrStr, err
}
