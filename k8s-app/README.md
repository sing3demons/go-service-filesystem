
```bash
kubectl apply -f .\k8s.yml
```

```bash
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.9.5/deploy/static/provider/aws/deploy.yaml
kubectl apply -f .\ingress-nginx.yaml
```

```bash
kubectl get po -n go-filesystem
kubectl -n go-filesystem get svc
```