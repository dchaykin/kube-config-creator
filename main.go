package main

import (
	"flag"
	"fmt"
	"path/filepath"

	"github.com/dchaykin/kube-config-creator/creator"
	"github.com/dchaykin/kube-config-creator/input"
	"github.com/setlog/panik"
	"k8s.io/client-go/util/homedir"
)

func main() {

	kubeconfig := getKubeConfig()

	flag.Parse()

	defer panik.Handle(func(r interface{}) {
		fmt.Println(r)
	})

	userInput := input.UserInput{}

	fmt.Println("Input config data for a new service account (use letters, numbers, hyphens and underscores only)")

	fmt.Print("Service Account (e.g. john-smith): ")
	panik.OnError(userInput.ReadServiceAccountName())

	fmt.Print("Cluster Role (e.g. pod-reader): ")
	panik.OnError(userInput.ReadRoleName())

	fmt.Print("Namespace (default, if empty): ")
	panik.OnError(userInput.ReadNamespace())

	clientset := creator.GetClientSet(kubeconfig)

	sa := creator.CreateServiceAccount(clientset, userInput.GetNamespace(), userInput.GetServiceAccountName())
	token := creator.GetServiceAccountToken(clientset, sa)
}

func getKubeConfig() *string {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	return kubeconfig
}
