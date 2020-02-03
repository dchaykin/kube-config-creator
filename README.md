#### kube-config-creator
Go app to create .kube/config for a service account

##### Prerequisite
You must have an admin access to the kubernetes cluster in order to be able to create a new kube-config file.

##### How to use kube-config-creator
```kube-config-creator --kubeconfig=[PATH_TO_ADMIN_CONFIG_FILE]```
Input:
* Service Account (mandatory): name of the service account to be created in the kubernetes cluster. If such service account already exists, it will be used without changing anything.
* Namespace (optional): namespace that the service account is going to be created in. If empty, `default` will be used.
Output:
* File `kube-config-<SERVICE_ACCOUNT>.yaml`. Copy this file into your ${HOME}/.kube/config to activate this account.

##### Note
This app activates a service account on a client computer, but does not create any roles or rolebindings in the cluster. You as administrator still have to configure the access to the cluster resources by creating RBAC-policies.