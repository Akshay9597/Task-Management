deploy:
	kubectl apply -f redis-deployment.yaml
	kubectl apply -f redis-service.yaml
	sleep 10
	kubectl port-forward svc/redis 6379:6379

delete:
	kubectl delete -f redis-deployment.yaml
	kubectl delete -f redis-service.yaml