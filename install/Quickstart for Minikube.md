Quickstart for Percona XtraDB Cluster Operator for Minikube
==========================================================================
Run Minikube to test a single-node Kubernetes cluster in a desktop environment.

Cloud Provider features, such as LoadBalancers, or features that require multiple nodes, such as Advanced Scheduling policies, will not work as expected in a Minikube cluster.


### Prerequisites
You must have [Minikube installed](https://kubernetes.io/docs/tasks/tools/install-minikube/) on your desktop.

### Installing the Percona XtraDB cluster
The `minikube start`command generates a context to communicate with the Minikube cluster. Starting Minikube may also start the virtual environment. The default amount of RAM allocated to the Minikube VM is 2048.


The Percona XtraDB Cluster has three MySQL nodes and a ProxySQL pod. A pod provides an environment for one or more containers in Kubernetes and always runs on a node, which is the worker machine.
You add an option to increase the RAM memory for the cluster on the desktop environment.

```bash
[~]$ minikube start --memory 6144
```
The return statements provide the progress of the minikube creation.

**Note:** *The exact return statements may be different and are based on your desktop environment.*

```
😄  minikube v1.0.0 on darwin (amd64)
🤹  Downloading Kubernetes v1.14.0 images in the background ...
💡  Tip: Use 'minikube start -p <name>' to create a new cluster, or 'minikube delete' to delete this one.
🔄  Restarting existing virtualbox VM for "minikube" ...
⌛  Waiting for SSH access ...
📶  "minikube" IP address is 192.168.99.103
🐳  Configuring Docker as the container runtime ...
🐳  Version of container runtime is 18.06.2-ce
⌛  Waiting for image downloads to complete ...
✨  Preparing Kubernetes environment ...
🚜  Pulling images required by Kubernetes v1.14.0 ...
🔄  Relaunching Kubernetes v1.14.0 using kubeadm ...
⌛  Waiting for pods: apiserver proxy etcd scheduler controller dns
📯  Updating kube-proxy configuration ...
🤔  Verifying component health ......
💗  kubectl is now configured to use "minikube"
🏄  Done! Thank you for using minikube!
```

### Clone and Download the repository

Use the `git clone` command to download the correct branch of the **percona-xtradb-cluster-operator** repository.  
**Note:** *In the git statement, you must specify the correct branch with the -b option.*
```bash
$ git clone -b release-0.3.0 https://github.com/percona/percona-xtradb-cluster-operator
```
After the repository is downloaded, change to the newly-created directory to run the rest of the commands in this document:

    cd percona-xtradb-cluster-operator

### Adjusting the Pod Memory Allocation
To ensure that the operator works on a desktop environment, make the following modifications to settings in the cr.yaml file.

**Note:** *Shell utilities differ based on the Linux/Unix/BSD variant. The GNU sed works with:
`sed -i`
A MacOS environment sed requires an extension:
`sed -i.bak`
The command may require an edit in your shell before running.*

To reduce the CPU usage:
```bash
$ sed -i 's/600m/200m/g' deploy/cr.yaml
```
You can use `grep` to verify the modification:
```bash
$ grep "cpu" deploy/cr.yaml
        cpu: 200m
#        cpu: "1"
        cpu: 200m
#        cpu: 700m
```
Run all of the PXC instances in one node by changing the antiAffinityTopologyKey settings:
```bash
$ sed -i 's/kubernetes\.io\/hostname/none/g' deploy/cr.yaml
```
Verify the updates:
```bash
$ grep -i  "topology" deploy/cr.yaml
      antiAffinityTopologyKey: "none"
      antiAffinityTopologyKey: "none"
```
### Creating the PXC and operator Objects
Use the following commands to create the PXC and operator objects.
1. Add extensions to the Kubernetes core for the PXC cluster and operator.
```bash
$ kubectl apply -f deploy/crd.yaml
```
The return statement confirms the action.
```
customresourcedefinition.apiextensions.k8s.io/perconaxtradbclusters.pxc.percona.com created
customresourcedefinition.apiextensions.k8s.io/perconaxtradbbackups.pxc.percona.com created
```
2. A Kubernetes namespace provides a scope for names.
```bash
$ kubectl create namespace pxc
```
The return statement confirms the action.
```
namespace/pxc created
```
3. Set the current context for the rest of the commands. The namespace option defines the namespace scope.
```bash
$ kubectl config set-context $(kubectl config current-context) --namespace=pxc
```
The return statement confirms the action.
```
Context "minikube" modified.
```
5. Role-based access control (RBAC) is a way to regulate access to the computer or network resources.
```bash
$ kubectl apply -f deploy/rbac.yaml
```
The return statement confirms the action.
```
role.rbac.authorization.k8s.io/percona-xtradb-cluster-operator created
rolebinding.rbac.authorization.k8s.io/default-account-percona-xtradb-cluster-operator created
```
6. An operator provides a method of managing a Kubernetes application, which is both deployed on Kubernetes and managed using the Kubernetes API.
```bash
$ kubectl apply -f deploy/operator.yaml
```
The return statement confirms the action.
```
deployment.apps/percona-xtradb-cluster-operator created
```
7. Store and manage sensitive information, such as passwords or ssh keys, in the secrets.yaml. Managing the application's secret objects with the file provides flexibility.
```bash
$ kubectl apply -f deploy/secrets.yaml
```
The return statement confirms the action.
```
secret/my-cluster-secrets created
```
8. Create the custom resource with the following command:
```bash
$ kubectl apply -f deploy/cr.yaml
```
The return statement confirms the action.
```
perconaxtradbcluster.pxc.percona.com/cluster1 created
```

### Connecting to the Cluster
After a few minutes, to allow for the cluster to be updated, use the following command to list the pods and their status:

```bash
$ kubectl get pods
```
The return statement lists the requested information.
```
NAME                                               READY   STATUS    RESTARTS   AGE
cluster1-proxysql-0                                3/3     Running   0          5m15s
cluster1-pxc-0                                     1/1     Running   0          5m15s
cluster1-pxc-1                                     1/1     Running   0          3m21s
cluster1-pxc-2                                     0/1     Running   0          78s
percona-xtradb-cluster-operator-776b6fd57d-w5dcc   1/1     Running   0          7m25s
```
If needed to troubleshoot a pod, you can use a `kubectl describe` command for the pod details.  For example, the command returns information for a selected pod:
```bash
kubectl describe pod cluster1-pxc-0
```
The resource's detailed description is returned. Review the description for `Warning` information and then correct the configuration. An example of a Warning is as follows:
```
Warning  FailedScheduling  68s (x4 over 2m22s)  default-scheduler  0/1 nodes are available: 1 node(s) didn't match pod affinity/anti-affinity, 1 node(s) didn't satisfy existing pods anti-affinity rules.
```

To connect to the cluster, you deploy a client shell access.
```bash
$ kubectl run -i --rm --tty percona-client --image=percona:5.7 --restart=Never -- bash -il
```
The return statement returns the `bash` shell prompt.
```
If you don't see a command prompt, try pressing enter.

bash-4.2$
```
With the shell prompt, connect to the MySQL server. You provide the MySQL host name, user name, and password.

```bash
bash-4.2$ mysql -h cluster1-proxysql -uroot -proot_password
```
The return statements provides the `mysql` prompt.
```
mysql: [Warning] Using a password on the command line interface can be insecure.
Welcome to the MySQL monitor.  Commands end with ; or \g.
Your MySQL connection id is 113
Server version: 5.5.30 (ProxySQL)

Copyright (c) 2009-2019 Percona LLC and/or its affiliates
Copyright (c) 2000, 2019, Oracle and/or its affiliates. All rights reserved.

Oracle is a registered trademark of Oracle Corporation and/or its
affiliates. Other names may be trademarks of their respective
owners.

Type 'help;' or '\h' for help. Type '\c' to clear the current input statement.
```
In MySQL, display the values of selected MySQL system variables.
```
mysql> SHOW VARIABLES like "max_connections";
```
The return statement displays the matching variable.
```
+-----------------+-------+
| Variable_name   | Value |
+-----------------+-------+
| max_connections | 151   |
+-----------------+-------+
1 row in set (0.00 sec)

mysql>
```
### Troubleshooting
You can modify your Minikube start process by adding options. To find which options are available, run the following command:
```bash
[~]$ minikube start --help
```
The phrases in the install/secrets.yaml can be decoded with the following:
```bash
$ echo -n `phrase` | base64 -D
```
If needed, use the `kubectl describe` command for the pod details.  For example, this command returns information for the selected pod:
```bash
kubectl describe pod cluster1-pxc-0
```
Review the detailed information for `Warning` statements and then correct the configuration. An example of a warning is as follows:
```
Warning  FailedScheduling  68s (x4 over 2m22s)  default-scheduler  0/1 nodes are available: 1 node(s) didn't match pod affinity/anti-affinity, 1 node(s) didn't satisfy existing pods anti-affinity rules.
```

### Clean Up
You can clean up the minikube cluster with the following command:

```bash
minikube stop
```
To delete the resource, use the following command:
```bash
minikube delete
```