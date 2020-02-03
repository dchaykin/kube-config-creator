package creator

import (
	"fmt"

	"github.com/setlog/panik"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	clientApi "k8s.io/client-go/tools/clientcmd/api"
)

type ConfigCreator struct {
	clientset      *kubernetes.Clientset
	kubeconfigPath string
	serviceAccount *v1.ServiceAccount
}

func (cc *ConfigCreator) Login(kubeconfig *string) {
	cc.kubeconfigPath = *kubeconfig
	cc.clientset = getClientSet(kubeconfig)
}

func (cc *ConfigCreator) CreateServiceAccount(namespace, accountName string) {
	coreV1 := cc.clientset.CoreV1()
	var err error
	cc.serviceAccount, err = coreV1.ServiceAccounts(namespace).Get(accountName, metav1.GetOptions{})

	if !errors.IsNotFound(err) {
		fmt.Printf("Service account %s already exists in %s.\n", accountName, namespace)
		return
	}

	serviceAccount := &v1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      accountName,
			Namespace: namespace,
		},
	}

	cc.serviceAccount, err = coreV1.ServiceAccounts(namespace).Create(serviceAccount)
	panik.OnError(err)
}

func (cc *ConfigCreator) CreateNew() {
	nc, err := clientcmd.LoadFromFile(cc.kubeconfigPath)
	panik.OnError(err)

	if len(nc.Clusters) == 0 {
		panik.Panicf("No clusters found in %s.", cc.kubeconfigPath)
	}

	conf := clientApi.NewConfig()
	clientApi.MinifyConfig(conf)
	name, cluster := createClusterFromFirst(nc.Clusters)
	conf.Clusters[name] = cluster

	userName := cc.serviceAccount.GetName()

	authInfo := clientApi.NewAuthInfo()
	authInfo.Token = cc.GetSecretToken()
	conf.AuthInfos[userName] = authInfo

	contextName := name + "-context"
	context := clientApi.NewContext()
	context.Cluster = name
	context.AuthInfo = userName
	if cc.serviceAccount.GetNamespace() != "default" {
		context.Namespace = cc.serviceAccount.GetNamespace()
	}

	conf.Contexts[contextName] = context
	conf.CurrentContext = contextName

	fileName := fmt.Sprintf("kube-config-%s.yaml", userName)
	panik.OnError(clientcmd.WriteToFile(*conf, fileName))

	fmt.Printf("%s has been created. Copy this file into your kube-config path, e.g. $HOME/.kube/config\n", fileName)
}

func (cc *ConfigCreator) GetSecretToken() string {
	return getServiceAccountToken(cc.clientset, cc.serviceAccount)
}

func createClusterFromFirst(clusters map[string]*clientApi.Cluster) (string, *clientApi.Cluster) {
	for name, info := range clusters {
		cluster := clientApi.NewCluster()
		cluster.Server = info.Server
		cluster.CertificateAuthorityData = info.CertificateAuthorityData
		return name, cluster
	}
	return "", nil
}

func getClientSet(kubeconfig *string) *kubernetes.Clientset {
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	panik.OnError(err)

	clientset, err := kubernetes.NewForConfig(config)
	panik.OnError(err)

	return clientset
}

func getServiceAccountToken(clientset *kubernetes.Clientset, serviceAccount *v1.ServiceAccount) string {
	for _, secret := range serviceAccount.Secrets {
		s, err := clientset.CoreV1().Secrets(getNamespace(secret.Namespace)).Get(secret.Name, metav1.GetOptions{})
		panik.OnError(err)

		if data, ok := s.Data["token"]; ok {
			return string(data)
		}
	}

	panik.Panicf("Token not found for the service account %s:%s", serviceAccount.Namespace, serviceAccount.Name)
	return ""
}

func getNamespace(name string) string {
	if name == "" {
		return "default"
	}
	return name
}
