
```bash
kubectl apply -f elasticsearch
kubectl port-forward es-cluster-0 9200:9200
curl http://localhost:9200/_cluster/health/?pretty
```

```bash
kubectl apply -f kibana
kubectl port-forward <kibana-pod-name> 5601:5601
curl http://localhost:5601/app/kibana
```

```bash
kubectl apply -f fluentd
```

```bash
kubectl delete -f fluentd
kubectl delete -f elasticsearch
kubectl delete -f kibana
```

```bash
kubectl get po
```