## Requirements

- Docker
- Kind
- Go 1.15 or higher
- Python 3.8 or higher (optional)


## Setup
You will need install [Docker](https://www.docker.com/>) and [Kind](https://kind.sigs.k8s.io) to do this workshop. Kind is a tool for running local Kubernetes clusters using Docker container “nodes”.
kind was primarily designed for testing Kubernetes itself, but may be used for local development or CI. Kind will require you to also install Golang (Go) in order to run in.

### MacOS Setup
Install and setup Docker Desktop if you have not already - <https://www.docker.com/>
Install Go - <https://go.dev/>
You may need to add the Go bin path to your profile of choice. For example I added the below line to my .zshrc file in `/Users/jhutchins:
export PATH=“/Users/jhutchins/go/bin:$PATH”`
Now open a terminal and run: go install sigs.k8s.io/kind@v0.17.0
If all was done correctly, you should now be able to run kind create cluster
### Windows Setup
Install and setup Docker Desktop if you have not already - <https://www.docker.com/>
Install Go - <https://go.dev/>
Open command prompt and run go install sigs.k8s.io/kind@v0.17.0
If all was done correctly, you should now be able to run kind create cluster

### Configure Kind
`kind create cluster --config kind-config.yaml`

This command tells kind to create a k8s cluster inside Docker that concists of single control plane and 3 worker.

### Create the deployment and service

Make sure Kind is running by executing `kubectl get nodes` at your command prompt and verifying that `kind-control-plane` and at least one kind node appear with a status of `Ready`


### Build and push the Docker image
You can deploy the service to Kubernetes using the following commands:

`docker build -t jg-hello-world:latest .`

`kind load docker-image jg-hello-world`


### Create the deployment and service

`kubectl apply -f deployment.yaml`

### Use the following commands to assess the status of your pod
`kubectl get pods`
`kubectl logs -f <pod name>`

### Expose the load balancer service

`kubectl port-forward svc/jg-helloworld-deploy 8888:8888`

## Troubleshooting

`kubectl apply -f deployment.yaml`

`kubectl delete -f deployment.yaml`

`kubectl describe pods` 

This command will return a basic description of each of your Pods, including their state. In the output, you’ll also be able to see if you have reached CPU, memory, or network limits. This is one of the most likely reasons for a pod remaining in the “pending” state.

`kubectl describe svc jg-helloworld-deploy`

The `kubectl describe svc jg-helloworld-deploy` command will display detailed information about the Kubernetes Service object named jg-helloworld-deploy. This includes information such as the IP address of the Service, the ports it's listening on, the endpoints it's routing traffic to, and more.

`kubectl get pods`

Provides a description of all pods in the k8s cluster.

`kubectl get svc`

Describes description of all services in the k8s cluster

### Testing
`curl localhost:8888`