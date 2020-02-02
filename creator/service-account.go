package creator

import (
	"fmt"
	"io/ioutil"

	"github.com/setlog/panik"
	"gopkg.in/yaml.v2"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type ServiceAccount struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name string `yaml:"name"`
	} `yaml:"metadata"`
}

func (sa *ServiceAccount) CreateYaml(fileName, serviceAccountName string) {
	sa.APIVersion = "v1"
	sa.Kind = "ServiceAccount"
	sa.Metadata.Name = serviceAccountName

	out, err := yaml.Marshal(sa)
	panik.OnError(err)

	panik.OnError(createFile(fileName, out))
}

func createFile(fileName string, out []byte) error {
	return ioutil.WriteFile(fileName, out, 0644)
}

func CreateServiceAccount(clientset *kubernetes.Clientset, namespace, accountName string) *v1.ServiceAccount {
	coreV1 := clientset.CoreV1()
	serviceAccount, err := coreV1.ServiceAccounts(namespace).Get(accountName, metav1.GetOptions{})

	if !errors.IsNotFound(err) {
		fmt.Printf("Service account %s already exists in %s.\n", accountName, namespace)
		return serviceAccount
	}

	serviceAccount = &v1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      accountName,
			Namespace: namespace,
		},
	}

	serviceAccount, err = clientset.CoreV1().ServiceAccounts(namespace).Create(serviceAccount)
	panik.OnError(err)
	return serviceAccount
}
