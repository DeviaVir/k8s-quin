# k8s-quin

Pronunciation: IPA(key): /kwɪn/, [kʰw̥ɪn].

Meaning: What or Which.

## Usage

```
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: network-exporter
  namespace: metrics
  labels:
    k8s-app: network-exporter
    kubernetes.io/cluster-service: "true"
    addonmanager.kubernetes.io/mode: Reconcile
    version: v0.0.2
spec:
  selector:
    matchLabels:
      k8s-app: network-exporter
      version: v0.0.2
  updateStrategy:
    type: RollingUpdate
  template:
    metadata:
      labels:
        k8s-app: network-exporter
        version: v0.0.2
      annotations:
         prometheus.io/scrape: "true"
         prometheus.io/port: "9666"
    spec:
      tolerations:
        - key: "key"
          operator: "Exists"
          effect: "NoSchedule"
        - key: "role"
          operator: "Exists"
          effect: "NoSchedule"
      containers:
        - name: network-exporter
          image: "deviavir/k8s-quin:0.0.2"
          imagePullPolicy: "IfNotPresent"
          ports:
            - name: metrics
              containerPort: 9666
              hostPort: 9666
          resources:
            limits:
              cpu: 50m
              memory: 100Mi
            requests:
              cpu: 10m
              memory: 50Mi
      hostNetwork: true
```

## Building/Deploying

A simple docker build in this repo does the trick.
