package main

import (
	"context"

	"k8s.io/apimachinery/pkg/api/errors"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	microk8scmd = map[string]string{
		"config":     "microk8s.config",
		"inspect":    "sudo microk8s.inspect",
		"iptable":    "sudo iptables -P FORWARD ACCEPT",
		"docker":     "microk8s.docker version",
		"get pods":   "microk8s.kubectl get pods --all-namespaces",
		"deployment": "microk8s.kubectl get all",
	}
)

func main() {
	kubeconfig := ""
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	_, err = clientset.CoreV1().Pods("").List(context.Background(), metav1.ListOptions{})
	if k8serrors.IsNotFound(err) {

	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		println(statusError.Error())
	} else if err != nil {
		panic(err.Error())
	}
}
