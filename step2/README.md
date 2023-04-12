### Setup
1. `cd step2`
2. Run `docker build . -t gif-service`
3. Run `kind load docker-image gif-service`
4. Run `kubectl apply -f redis.yaml`
5. Run `kubectl apply -f deployment.yaml`
6. In a different terminal, run`kubectl port-forward svc/gif-service 8889:8889`
7. Visit `http://localhost:8889` and upload some gifs
8. Change the replicas to `2` in `/step2/deployment.yaml`
9. Run kubectl apply -f deployment.yaml` again.
10. Check `kubectl get pods` and notice how a second pod has been spawned.
11. Run `curl http://localhost:8889/mine-bitcoin/?iterations=10`
12. Run `kubectl logs -l app=gif`


### Autohealing demo

1. `kubectl get pods`
2. Locate the gif-deployment pod(s)
3. Run `kubectl delete pod <gif-deployment-pod-name>`
4. Run `kubectl get pods` again. Notice how k8s autoheals when a pod becomes unhealthy.
5. Run `kubectl get nodes` and note the note configuration.
6. Run`kubectl get pod <gif-deployment-pod-name> -o wide` and note the value in the NODE column
7. Delete one of the worker nodes via `kubectl delete node <node-name>`
8. Note how k8s responds in order to maintain ahe desired state via running `kubectl get pods` and `kubectl get pod <gif-deployment-pod-name> -o wide` again.


### Troubleshooting
You can then monitor the status of the Redis Deployment or the Gif service and Autoscaler using the `kubectl get deployment` and `kubectl get hpa gif-deployment`, respectively.

