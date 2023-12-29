# 使用kustomize发布到k8s集群

1 创建ingress
按需要更新ingress.yaml，确认其实的子域名和端口与实际情况匹配

```shell
k8prod create -f ingress.yaml
```

2 router 53中配置域名的A转发纪录

## 部署到k8s

preview
```
kubectl kustomize overlays/dev
```

更新
```
kubectl apply -k overlays/dev
```

