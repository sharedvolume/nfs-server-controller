# NFS Client Examples

This directory contains examples of how to use the NFS servers created by the NFS Server Controller.

## Basic NFS Client Pod

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: nfs-client-basic
  namespace: default
spec:
  containers:
  - name: nfs-client
    image: nginx:alpine
    volumeMounts:
    - name: nfs-storage
      mountPath: /data
    command: ["/bin/sh"]
    args: ["-c", "while true; do echo $(date) >> /data/test.log; sleep 60; done"]
  volumes:
  - name: nfs-storage
    nfs:
      server: basic-nfs-server.default.svc.cluster.local
      path: /shared
```

## Deployment with NFS Storage

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: web-app-with-nfs
  namespace: default
spec:
  replicas: 3
  selector:
    matchLabels:
      app: web-app
  template:
    metadata:
      labels:
        app: web-app
    spec:
      containers:
      - name: web-app
        image: nginx:alpine
        ports:
        - containerPort: 80
        volumeMounts:
        - name: web-content
          mountPath: /usr/share/nginx/html
        - name: logs
          mountPath: /var/log/nginx
      volumes:
      - name: web-content
        nfs:
          server: basic-nfs-server.default.svc.cluster.local
          path: /shared/web-content
      - name: logs
        nfs:
          server: basic-nfs-server.default.svc.cluster.local
          path: /shared/logs
```

## PersistentVolume and PersistentVolumeClaim

For a more Kubernetes-native approach, you can create PV and PVC:

```yaml
# PersistentVolume pointing to NFS server
apiVersion: v1
kind: PersistentVolume
metadata:
  name: nfs-pv
spec:
  capacity:
    storage: 10Gi
  accessModes:
    - ReadWriteMany
  nfs:
    server: basic-nfs-server.default.svc.cluster.local
    path: /shared
  persistentVolumeReclaimPolicy: Retain

---

# PersistentVolumeClaim
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: nfs-pvc
  namespace: default
spec:
  accessModes:
    - ReadWriteMany
  resources:
    requests:
      storage: 10Gi
  volumeName: nfs-pv

---

# Pod using PVC
apiVersion: v1
kind: Pod
metadata:
  name: nfs-client-pvc
  namespace: default
spec:
  containers:
  - name: client
    image: alpine:latest
    command: ["/bin/sh"]
    args: ["-c", "while true; do echo $(date) >> /data/timestamp.log; sleep 30; done"]
    volumeMounts:
    - name: nfs-storage
      mountPath: /data
  volumes:
  - name: nfs-storage
    persistentVolumeClaim:
      claimName: nfs-pvc
```

## StatefulSet with NFS

```yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: database-with-nfs
  namespace: default
spec:
  serviceName: database
  replicas: 2
  selector:
    matchLabels:
      app: database
  template:
    metadata:
      labels:
        app: database
    spec:
      containers:
      - name: database
        image: postgres:13
        env:
        - name: POSTGRES_DB
          value: myapp
        - name: POSTGRES_USER
          value: user
        - name: POSTGRES_PASSWORD
          value: password
        - name: PGDATA
          value: /var/lib/postgresql/data/pgdata
        volumeMounts:
        - name: data
          mountPath: /var/lib/postgresql/data
        - name: backup
          mountPath: /backup
      volumes:
      - name: data
        nfs:
          server: high-performance-nfs.production.svc.cluster.local
          path: /fast-storage/database
      - name: backup
        nfs:
          server: basic-nfs-server.default.svc.cluster.local
          path: /shared/database-backup
```

## Job with NFS for Batch Processing

```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: data-processing-job
  namespace: default
spec:
  template:
    spec:
      containers:
      - name: processor
        image: python:3.9
        command: ["python"]
        args: ["-c", "import os; print('Processing files:', os.listdir('/input')); open('/output/result.txt', 'w').write('Processing complete')"]
        volumeMounts:
        - name: input-data
          mountPath: /input
          readOnly: true
        - name: output-data
          mountPath: /output
      volumes:
      - name: input-data
        nfs:
          server: basic-nfs-server.default.svc.cluster.local
          path: /shared/input
      - name: output-data
        nfs:
          server: basic-nfs-server.default.svc.cluster.local
          path: /shared/output
      restartPolicy: Never
  backoffLimit: 3
```

## Multi-Container Pod Sharing Data

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: multi-container-nfs
  namespace: default
spec:
  containers:
  - name: writer
    image: alpine:latest
    command: ["/bin/sh"]
    args: ["-c", "while true; do echo 'Writer: '$(date) >> /shared/data.log; sleep 10; done"]
    volumeMounts:
    - name: shared-storage
      mountPath: /shared
  - name: reader
    image: alpine:latest
    command: ["/bin/sh"]
    args: ["-c", "while true; do echo 'Last 5 lines:'; tail -5 /shared/data.log; sleep 30; done"]
    volumeMounts:
    - name: shared-storage
      mountPath: /shared
      readOnly: true
  volumes:
  - name: shared-storage
    nfs:
      server: basic-nfs-server.default.svc.cluster.local
      path: /shared
```

## ConfigMap and Secret Storage on NFS

```yaml
# Store config files on NFS
apiVersion: v1
kind: Pod
metadata:
  name: app-with-nfs-config
  namespace: default
spec:
  containers:
  - name: app
    image: nginx:alpine
    volumeMounts:
    - name: app-config
      mountPath: /etc/app
    - name: nginx-config
      mountPath: /etc/nginx/conf.d
      readOnly: true
  volumes:
  - name: app-config
    nfs:
      server: basic-nfs-server.default.svc.cluster.local
      path: /shared/config
  - name: nginx-config
    nfs:
      server: basic-nfs-server.default.svc.cluster.local
      path: /shared/nginx-config
```

## Troubleshooting NFS Client Issues

### Common Issues and Solutions

1. **Mount Permission Denied**
   ```bash
   # Check if the NFS server is running
   kubectl get nfsserver basic-nfs-server
   
   # Check NFS server pod logs
   kubectl logs -l app=basic-nfs-server
   
   # Test NFS connectivity from a debug pod
   kubectl run nfs-debug --image=alpine:latest --rm -it -- /bin/sh
   # Inside the pod:
   apk add nfs-utils
   showmount -e basic-nfs-server.default.svc.cluster.local
   ```

2. **Stale File Handle**
   ```bash
   # Restart the pod using NFS
   kubectl delete pod nfs-client-basic
   kubectl apply -f nfs-client-basic.yaml
   ```

3. **Performance Issues**
   - Use the high-performance NFS server for I/O intensive workloads
   - Consider using local storage for temporary files
   - Implement proper caching strategies

### Debug Pod for NFS Testing

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: nfs-debug
  namespace: default
spec:
  containers:
  - name: debug
    image: alpine:latest
    command: ["/bin/sh"]
    args: ["-c", "apk add nfs-utils curl && sleep 3600"]
    volumeMounts:
    - name: nfs-test
      mountPath: /nfs-test
  volumes:
  - name: nfs-test
    nfs:
      server: basic-nfs-server.default.svc.cluster.local
      path: /shared
```

Use this pod to test NFS connectivity and perform debugging:

```bash
kubectl exec -it nfs-debug -- /bin/sh

# Test NFS mount
ls -la /nfs-test/

# Test write access
echo "test" > /nfs-test/test.txt

# Test NFS server status
showmount -e basic-nfs-server.default.svc.cluster.local

# Check mount options
mount | grep nfs
```
