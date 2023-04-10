### Setup
`kubectl apply -f redis.yaml`
`docker build -t gif-service`
`kind load docker-image jg-hello-world`
`kubectl apply -f deployment.yaml`

`kubectl port-forward svc/gif-service 8888:8889`
### Troubleshooting
You can then monitor the status of the Deployment and Autoscaler using the `kubectl get deployment` and `kubectl get hpa commands`, respectively.

