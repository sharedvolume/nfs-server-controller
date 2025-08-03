# NFS Server Controller

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/sharedvolume/nfs-server-controller)](https://goreportcard.com/report/github.com/sharedvolume/nfs-server-controller)
[![Docker Pulls](https://img.shields.io/docker/pulls/sharedvolume/nfs-server-controller)](https://hub.docker.com/r/sharedvolume/nfs-server-controller)
[![Release](https://img.shields.io/github/release/sharedvolume/nfs-server-controller.svg)](https://github.com/sharedvolume/nfs-server-controller/releases)

A Kubernetes operator that manages NFS servers as custom resources, providing dynamic provisioning and lifecycle management of NFS services within your cluster.

## Features

- **Custom Resource Definition (CRD)**: Define NFS servers declaratively using Kubernetes resources
- **Dynamic Provisioning**: Automatically provision NFS servers with persistent storage
- **Lifecycle Management**: Handle creation, updates, and deletion of NFS server instances
- **Storage Flexibility**: Support for both StorageClass-based and pre-existing PersistentVolume storage
- **High Availability**: Configurable replica count for NFS server instances
- **Service Discovery**: Automatic service creation for NFS server connectivity
- **Status Monitoring**: Real-time status updates and health checks

## Quick Start

### Prerequisites

- Kubernetes cluster (v1.20+)
- kubectl configured to access your cluster
- Cluster admin permissions

### Installation

1. **Install the CRDs and operator:**
   ```bash
   kubectl apply -f https://github.com/sharedvolume/nfs-server-controller/releases/latest/download/install.yaml
   ```

2. **Verify the installation:**
   ```bash
   kubectl get deployment -n nfs-server-controller-system
   kubectl get crd nfsservers.sharedvolume.io
   ```

### Creating an NFS Server

Create an NFS server using a StorageClass:

```yaml
apiVersion: sharedvolume.io/v1alpha1
kind: NfsServer
metadata:
  name: my-nfs-server
  namespace: default
spec:
  storage:
    capacity: "10Gi"
    storageClassName: "fast-ssd"
  replicas: 2
  path: "/shared"
```

Apply the configuration:
```bash
kubectl apply -f nfs-server.yaml
```

### Using the NFS Server

Once the NFS server is running, you can mount it in your pods:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: nfs-client
spec:
  containers:
  - name: app
    image: nginx
    volumeMounts:
    - name: nfs-storage
      mountPath: /data
  volumes:
  - name: nfs-storage
    nfs:
      server: my-nfs-server.default.svc.cluster.local
      path: /shared
```

## Configuration

### NfsServer Spec

| Field | Type | Description | Required |
|-------|------|-------------|----------|
| `storage.capacity` | string | Storage capacity (e.g., "10Gi") | Yes |
| `storage.storageClassName` | string | StorageClass name for dynamic provisioning | No* |
| `storage.persistentVolume` | string | Pre-existing PersistentVolume name | No* |
| `replicas` | int32 | Number of NFS server replicas (default: 2) | No |
| `path` | string | NFS export path (default: "/nfs") | No |
| `image` | string | NFS server image (default: auto-detected) | No |

*Either `storageClassName` or `persistentVolume` must be specified, but not both.

### Examples

#### Using a specific PersistentVolume:
```yaml
apiVersion: sharedvolume.io/v1alpha1
kind: NfsServer
metadata:
  name: nfs-with-pv
spec:
  storage:
    capacity: "50Gi"
    persistentVolume: "my-existing-pv"
  replicas: 1
```

#### Custom NFS image and path:
```yaml
apiVersion: sharedvolume.io/v1alpha1
kind: NfsServer
metadata:
  name: custom-nfs
spec:
  storage:
    capacity: "20Gi"
    storageClassName: "standard"
  image: "sharedvolume/nfs-server:custom"
  path: "/exports"
  replicas: 3
```

## Development

### Prerequisites

- Go 1.24+
- Docker
- kubectl
- Kind (for local testing)

### Building from Source

1. **Clone the repository:**
   ```bash
   git clone https://github.com/sharedvolume/nfs-server-controller.git
   cd nfs-server-controller
   ```

2. **Build the manager:**
   ```bash
   make build
   ```

3. **Run tests:**
   ```bash
   make test
   ```

4. **Build Docker image:**
   ```bash
   make docker-build IMG=nfs-server-controller:dev
   ```

### Local Development

1. **Install CRDs:**
   ```bash
   make install
   ```

2. **Run the controller locally:**
   ```bash
   make run
   ```

3. **Run tests with Kind:**
   ```bash
   make test-e2e
   ```

### Contributing

We welcome contributions! Please see our [Contributing Guidelines](CONTRIBUTING.md) for details on:

- Code of conduct
- Development setup
- Pull request process
- Testing requirements

For detailed development information, including how this project was built with Kubebuilder, see our [Development Guide](docs/development.md).

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass
6. Submit a pull request

## Architecture

The NFS Server Controller consists of:

- **Custom Resource Definition (CRD)**: Defines the `NfsServer` resource schema
- **Controller**: Watches for `NfsServer` resources and manages their lifecycle
- **Reconciler**: Ensures the desired state matches the actual state by creating/updating:
  - PersistentVolumeClaims for storage
  - ReplicaSets for NFS server pods
  - Services for network access

## Monitoring

The controller provides the following status information:

```bash
kubectl get nfsservers
NAME            READY   ADDRESS                              CAPACITY
my-nfs-server   true    my-nfs-server.default.svc.cluster.local   10Gi
```

For detailed status:
```bash
kubectl describe nfsserver my-nfs-server
```

## Troubleshooting

### Common Issues

1. **NFS Server not ready:**
   - Check PVC status: `kubectl get pvc`
   - Verify storage class exists: `kubectl get storageclass`
   - Check pod logs: `kubectl logs -l app=my-nfs-server`

2. **Mount issues from clients:**
   - Ensure NFS client utilities are installed in client pods
   - Verify network policies allow NFS traffic
   - Check service endpoints: `kubectl get endpoints my-nfs-server`

3. **Permission issues:**
   - Verify the controller has proper RBAC permissions
   - Check if security policies allow privileged containers

### Logs

View controller logs:
```bash
kubectl logs -n nfs-server-controller-system deployment/nfs-server-controller-manager
```

## Security Considerations

- NFS server pods run with privileged security context (required for NFS functionality)
- Ensure proper network policies to restrict NFS access
- Consider using storage encryption for sensitive data
- Regularly update the NFS server image for security patches

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

- **Issues**: Report bugs and feature requests on [GitHub Issues](https://github.com/sharedvolume/nfs-server-controller/issues)
- **Discussions**: Join community discussions on [GitHub Discussions](https://github.com/sharedvolume/nfs-server-controller/discussions)
- **Security**: Report security vulnerabilities privately to bilgehan.nal@gmail.com

## Acknowledgments

Built with [Kubebuilder](https://github.com/kubernetes-sigs/kubebuilder) and inspired by the Kubernetes community's best practices for operators.