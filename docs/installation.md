# Installation Guide

This guide covers different methods to install the NFS Server Controller in your Kubernetes cluster.

## Prerequisites

- Kubernetes cluster v1.20 or later
- kubectl configured to access your cluster
- Cluster admin permissions
- Container runtime that supports privileged containers

## Method 1: Using Release Manifests (Recommended)

This is the simplest method for production deployments.

### Install Latest Release

```bash
kubectl apply -f https://github.com/sharedvolume/nfs-server-controller/releases/latest/download/install.yaml
```

### Install Specific Version

```bash
VERSION=v0.0.11
kubectl apply -f https://github.com/sharedvolume/nfs-server-controller/releases/download/${VERSION}/install.yaml
```

### Verify Installation

```bash
# Check if the controller is running
kubectl get deployment -n nfs-server-controller-system

# Check if CRDs are installed
kubectl get crd nfsservers.sharedvolume.io

# Check controller logs
kubectl logs -n nfs-server-controller-system deployment/nfs-server-controller-manager
```

## Method 2: Using Kustomize

For customized deployments, you can use Kustomize.

### Clone Repository

```bash
git clone https://github.com/sharedvolume/nfs-server-controller.git
cd nfs-server-controller
```

### Default Installation

```bash
make deploy IMG=sharedvolume/nfs-server-controller:latest
```

### Custom Image

```bash
make deploy IMG=your-registry/nfs-server-controller:your-tag
```

### Custom Namespace

Create a `kustomization.yaml`:

```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:
- github.com/sharedvolume/nfs-server-controller/config/default?ref=main

namespace: your-namespace

images:
- name: controller
  newName: sharedvolume/nfs-server-controller
  newTag: v0.0.11
```

Apply:
```bash
kubectl apply -k .
```

## Method 3: Helm Chart

*Note: Helm chart is planned for future releases.*

## Method 4: Development Installation

For development and testing purposes.

### Prerequisites

- Go 1.24+
- Make
- Docker

### Install CRDs Only

```bash
git clone https://github.com/sharedvolume/nfs-server-controller.git
cd nfs-server-controller
make install
```

### Run Controller Locally

```bash
make run
```

This runs the controller outside the cluster, useful for development.

## Configuration

### Resource Limits

The default installation includes resource limits. To customize:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nfs-server-controller-manager
spec:
  template:
    spec:
      containers:
      - name: manager
        resources:
          limits:
            cpu: 500m
            memory: 512Mi
          requests:
            cpu: 100m
            memory: 64Mi
```

### Namespace Scope

By default, the controller watches all namespaces. To limit to specific namespaces, modify the deployment:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nfs-server-controller-manager
spec:
  template:
    spec:
      containers:
      - name: manager
        env:
        - name: WATCH_NAMESPACE
          value: "namespace1,namespace2"
```

### Logging

Configure log level:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nfs-server-controller-manager
spec:
  template:
    spec:
      containers:
      - name: manager
        args:
        - --log-level=info  # debug, info, error
```

## Security Considerations

### Network Policies

If using network policies, ensure NFS traffic is allowed:

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: allow-nfs
spec:
  podSelector:
    matchLabels:
      app: your-nfs-server
  policyTypes:
  - Ingress
  ingress:
  - ports:
    - protocol: TCP
      port: 2049  # NFS
    - protocol: TCP
      port: 20048 # mountd
    - protocol: TCP
      port: 111   # rpcbind
```

### Pod Security Standards

NFS servers require privileged access. Ensure your cluster allows this:

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: nfs-namespace
  labels:
    pod-security.kubernetes.io/enforce: privileged
    pod-security.kubernetes.io/audit: privileged
    pod-security.kubernetes.io/warn: privileged
```

### Service Account

The controller uses a service account with minimal required permissions. Review the RBAC configuration in the installation manifests.

## Troubleshooting

### Controller Not Starting

1. Check if CRDs are installed:
   ```bash
   kubectl get crd nfsservers.sharedvolume.io
   ```

2. Check controller logs:
   ```bash
   kubectl logs -n nfs-server-controller-system deployment/nfs-server-controller-manager
   ```

3. Verify RBAC permissions:
   ```bash
   kubectl auth can-i create nfsservers --as=system:serviceaccount:nfs-server-controller-system:nfs-server-controller-manager
   ```

### Image Pull Issues

If using a private registry:

```bash
kubectl create secret docker-registry regcred \
  --docker-server=your-registry.com \
  --docker-username=your-username \
  --docker-password=your-password \
  --docker-email=your-email
```

Update the deployment to use the secret:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nfs-server-controller-manager
spec:
  template:
    spec:
      imagePullSecrets:
      - name: regcred
```

### Resource Constraints

If pods are pending due to resource constraints:

```bash
kubectl describe pods -n nfs-server-controller-system
```

Adjust resource requests/limits as needed.

## Upgrading

### From Previous Versions

1. Check release notes for breaking changes
2. Backup existing NFS server configurations
3. Apply new manifests:
   ```bash
   kubectl apply -f https://github.com/sharedvolume/nfs-server-controller/releases/download/v0.0.11/install.yaml
   ```

### Rolling Back

To roll back to a previous version:

```bash
kubectl apply -f https://github.com/sharedvolume/nfs-server-controller/releases/download/v0.0.10/install.yaml
```

## Uninstallation

### Remove Controller

```bash
kubectl delete -f https://github.com/sharedvolume/nfs-server-controller/releases/latest/download/install.yaml
```

### Remove CRDs (Optional)

**Warning**: This will delete all NFS server instances.

```bash
kubectl delete crd nfsservers.sharedvolume.io
```

### Using Make

If installed via make:

```bash
make undeploy
make uninstall
```

## Next Steps

After installation:

1. [Create your first NFS server](examples.md)
2. [Review the API reference](api-reference.md)
3. [Set up monitoring](monitoring.md)
4. [Configure backup and disaster recovery](backup.md)
