# sample-controller
Kubernetes Controller to rolling restart a deployment pods on a change in configmap

## Description
This controller performs rolling restart of deployment pods which is inspired by (https://github.com/k8spatterns/examples/tree/main/advanced/Controller/expose-controller)[k8spatterns] with limited functionality of making changes in data of configmap

This works in default namespace for webapp-config configmap and webapp deployment. 

## Getting Started

### Prerequisites
- go version v1.22.0+
- docker version 17.03+.
- kubectl version v1.11.3+.
- Access to a Kubernetes v1.11.3+ cluster.

### To Run/Test On your Local Machine

Start the Minikube

```sh
minikube start
```
To deploy sample app so that we can run the operator and test how it is restarting the deployment on making change in configmap. 

```sh
kubectl apply -f ./config/samples/webapp.yaml

# to get the endpoint to access the webapp
minikube service list
sample output: 

|-----------|---------|-------------|---------------------------|
| NAMESPACE | NAME    | TARGET PORT |            URL            |
|-----------|---------|-------------|---------------------------|
| default   | webapp  | http/8080   | http://192.168.49.2:31803 |
|-----------|---------|-------------|---------------------------|

# curl URL
curl http://192.168.49.2:31803
output: sample data for configmap!

```
Run the operator on your local machine

```sh
make run
```

sample output:
```sh
/home/devops/projects/git/sample-controller/bin/controller-gen rbac:roleName=manager-role crd webhook paths="./..." output:crd:artifacts:config=config/crd/bases
/home/devops/projects/git/sample-controller/bin/controller-gen object:headerFile="hack/boilerplate.go.txt" paths="./..."
go fmt ./...
internal/controller/configmap_controller.go
go vet ./...
go run ./cmd/main.go
2024-07-30T17:18:24+05:30	INFO	setup	starting manager
2024-07-30T17:18:24+05:30	INFO	starting server	{"name": "health probe", "addr": "[::]:8081"}
2024-07-30T17:18:24+05:30	INFO	Starting EventSource	{"controller": "configmap", "controllerGroup": "", "controllerKind": "ConfigMap", "source": "kind source: *v1.ConfigMap"}
2024-07-30T17:18:24+05:30	INFO	Starting Controller	{"controller": "configmap", "controllerGroup": "", "controllerKind": "ConfigMap"}
2024-07-30T17:18:24+05:30	INFO	Starting workers	{"controller": "configmap", "controllerGroup": "", "controllerKind": "ConfigMap", "worker count": 1}

```

Now edit the data in configmap and save the changes

```sh
kubectl edit cm webapp-config
```

After making the changes in configmap, you will see the pods of webapp deployment getting restarted.

```sh
kubectl get pods -n default 
```

