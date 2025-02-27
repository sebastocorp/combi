package actionset

import (
	"bytes"
	"context"
	"os/exec"
	"reflect"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

type actionT struct {
	Name string
	On   string

	k8s actionK8sT `json:"-"`
	cmd []string   `json:"-"`
}

type actionK8sT struct {
	cfg    *rest.Config
	client *kubernetes.Clientset

	namespace string
	pod       string
	container string
}

type ActionResultT struct {
	Name   string `json:"name"`
	On     string `json:"on"`
	Cmd    string `json:"cmd"`
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
}

func (a *actionT) exec() (r ActionResultT, err error) {
	r.Name = a.Name
	r.On = a.On
	r.Cmd = strings.Join(a.cmd, " ")

	if reflect.TypeOf(a.k8s) == reflect.TypeOf(actionK8sT{}) {
		var stdout, stderr bytes.Buffer
		cmd := exec.Command(a.cmd[0], a.cmd[1:]...)
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		err = cmd.Run()
		r.Stdout = stdout.String()
		r.Stderr = stderr.String()
	} else {
		r.Stdout, r.Stderr, err = a.execK8sCommand()
	}

	return r, err
}

func (a *actionT) execK8sCommand() (stdoutStr, stderrStr string, err error) {
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
		Command:   a.cmd,
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
