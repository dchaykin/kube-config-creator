package creator

import (
	"github.com/setlog/panik"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func GetClientSet(kubeconfig *string) *kubernetes.Clientset {
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	panik.OnError(err)

	clientset, err := kubernetes.NewForConfig(config)
	panik.OnError(err)

	return clientset
}
