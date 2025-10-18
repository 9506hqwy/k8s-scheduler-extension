# Kubernetes Scheduler Extension

## Sample Scheduler Extender

Build Container image.

```sh
buildah bud -t <Scheduler Image Path> -f build/index-extender/Dockerfile .
```

Push container image.

```sh
podman push <Scheduler Image Path>
```

Create namespace.

```sh
kubectl create namespace sched-system
```

Deploy scheduler to target namespace.

```sh
cat | kubectl apply -f - <<EOF
apiVersion: v1
kind: ServiceAccount
metadata:
  name: index-scheduler
  namespace: sched-system
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: index-scheduler-kube-scheduler
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: system:kube-scheduler
subjects:
- kind: ServiceAccount
  name: index-scheduler
  namespace: sched-system
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: index-scheduler-volume-scheduler
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: system:volume-scheduler
subjects:
- kind: ServiceAccount
  name: index-scheduler
  namespace: sched-system
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: index-scheduler-config-reader
  namespace: kube-system
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: extension-apiserver-authentication-reader
subjects:
- kind: ServiceAccount
  name: index-scheduler
  namespace: sched-system
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: scheduler-config
  namespace: sched-system
data:
  config: |
    apiVersion: kubescheduler.config.k8s.io/v1
    kind: KubeSchedulerConfiguration
    leaderElection:
      leaderElect: false
    profiles:
    - schedulerName: index-scheduler
    extenders:
    - urlPrefix: http://127.0.0.1:10261/api/scheduler
      filterVerb: filter
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: index-scheduler
  namespace: sched-system
spec:
  selector:
    matchLabels:
      app: index-scheduler
  template:
    metadata:
      labels:
        app: index-scheduler
    spec:
      serviceAccountName: index-scheduler
      containers:
      - name: default-scheduler
        image: registry.k8s.io/kube-scheduler:v1.33.0
        command:
        - kube-scheduler
        - --config=/opt/config.yaml
        - --secure-port=10260
        ports:
        - containerPort: 10260
          protocol: TCP
        volumeMounts:
        - name: scheduler-config
          mountPath: /opt
      - name: index-extender
        image: <Scheduler Image Path>
        ports:
        - containerPort: 10261
          protocol: TCP
      volumes:
      - name: scheduler-config
        configMap:
          name: scheduler-config
          items:
          - key: config
            path: config.yaml
            mode: 0444
EOF
```

Deploy pod with specified scheduler.

```sh
cat | kubectl apply -f - <<EOF
apiVersion: v1
kind: Pod
metadata:
  name: sample-01
spec:
  schedulerName: index-scheduler
  containers:
  - name: sample-01
    image: nginx
EOF
```

## Sample Scheduler Plugin

Build Container image.

```sh
buildah bud -t <Scheduler Image Path> -f build/index-scheduler/Dockerfile .
```

Push container image.

```sh
podman push <Scheduler Image Path>
```

Create namespace.

```sh
kubectl create namespace sched-system
```

Deploy scheduler to target namespace.

```sh
cat | kubectl apply -f - <<EOF
apiVersion: v1
kind: ServiceAccount
metadata:
  name: index-scheduler
  namespace: sched-system
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: index-scheduler-kube-scheduler
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: system:kube-scheduler
subjects:
- kind: ServiceAccount
  name: index-scheduler
  namespace: sched-system
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: index-scheduler-volume-scheduler
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: system:volume-scheduler
subjects:
- kind: ServiceAccount
  name: index-scheduler
  namespace: sched-system
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: index-scheduler-config-reader
  namespace: kube-system
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: extension-apiserver-authentication-reader
subjects:
- kind: ServiceAccount
  name: index-scheduler
  namespace: sched-system
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: scheduler-config
  namespace: sched-system
data:
  config: |
    apiVersion: kubescheduler.config.k8s.io/v1
    kind: KubeSchedulerConfiguration
    leaderElection:
      leaderElect: false
    profiles:
    - schedulerName: index-scheduler
      plugins:
        filter:
          enabled:
          - name: IndexScheduling
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: index-scheduler
  namespace: sched-system
spec:
  selector:
    matchLabels:
      app: index-scheduler
  template:
    metadata:
      labels:
        app: index-scheduler
    spec:
      serviceAccountName: index-scheduler
      containers:
      - name: index-scheduler
        image: <Scheduler Image Path>
        command:
        - /index-scheduler
        - --config=/opt/config.yaml
        ports:
        - containerPort: 10259
          protocol: TCP
        volumeMounts:
        - name: scheduler-config
          mountPath: /opt
      volumes:
      - name: scheduler-config
        configMap:
          name: scheduler-config
          items:
          - key: config
            path: config.yaml
            mode: 0444
EOF
```

Deploy pod with specified scheduler.

```sh
cat | kubectl apply -f - <<EOF
apiVersion: v1
kind: Pod
metadata:
  name: sample-01
spec:
  schedulerName: index-scheduler
  containers:
  - name: sample-01
    image: nginx
EOF
```
