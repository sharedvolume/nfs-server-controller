# NFS Server Controller

**Repository**: [https://github.com/sharedvolume/nfs-server-controller](https://github.com/sharedvolume/nfs-server-controller)  
**Releases**: [https://github.com/sharedvolume/nfs-server-controller/releases](https://github.com/sharedvolume/nfs-server-controller/releases)

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Report Card](https://goreportcard.com/badge/github.com/sharedvolume/nfs-server-controller)](https://goreportcard.com/report/github.com/sharedvolume/nfs-server-controller)
[![Docker Pulls](https://img.shields.io/docker/pulls/sharedvolume/nfs-server-controller)](https://hub.docker.com/r/sharedvolume/nfs-server-controller)
[![Release](https://img.shields.io/github/release/sharedvolume/nfs-server-controller.svg)](https://github.com/sharedvolume/nfs-server-controller/releases)

A professional-grade Kubernetes operator that manages NFS servers as custom resources, providing dynamic provisioning and comprehensive lifecycle management of NFS services within your cluster. Part of the SharedVolume ecosystem for enterprise storage solutions.

---

## 📋 Table of Contents

- [Overview](#overview)
- [Enterprise Features](#enterprise-features)
- [Quick Start](#quick-start)
- [Configuration](#configuration)
- [Use Cases](#use-cases)
- [Architecture](#architecture)
- [Development](#development)
- [Roadmap](#roadmap)
- [Community](#community)
- [License](#license)

---

## Overview

**NFS Server Controller** is a professional-grade Kubernetes operator that provides enterprise-ready NFS server management through custom resources. It enables organizations to deploy and manage NFS servers declaratively within their Kubernetes clusters, offering a cloud-native approach to shared storage provisioning designed for production environments at scale.

### Problem Statement

Traditional NFS server deployment in Kubernetes environments involves:
- Manual pod and service configuration
- Complex storage management
- Inconsistent deployment patterns
- Limited automation and lifecycle management
- Difficulty in scaling and high availability setup

### Solution

The NFS Server Controller addresses these challenges by:
- **Declarative Management**: Define NFS servers as Kubernetes custom resources
- **Automated Provisioning**: Automatic creation of storage, pods, and services
- **Lifecycle Management**: Handle creation, updates, scaling, and deletion
- **Storage Flexibility**: Support for both dynamic provisioning and existing volumes
- **High Availability**: Built-in support for multiple replicas
- **Kubernetes Native**: Follows Kubernetes patterns and best practices

## Enterprise Features

- **Production-Ready Architecture**: Built with enterprise-grade reliability and scalability in mind
- **Custom Resource Definition (CRD)**: Define NFS servers declaratively using Kubernetes resources
- **Dynamic Provisioning**: Automatically provision NFS servers with persistent storage
- **Advanced Lifecycle Management**: Handle creation, updates, scaling, and deletion of NFS server instances
- **Storage Flexibility**: Support for both StorageClass-based and pre-existing PersistentVolume storage
- **High Availability**: Configurable replica count for NFS server instances with built-in redundancy
- **Service Discovery**: Automatic service creation for seamless NFS server connectivity
- **Comprehensive Monitoring**: Real-time status updates, health checks, and observability
- **Security-First Design**: Built with Kubernetes security best practices and RBAC integration
- **Multi-Platform Support**: Container images for multiple architectures (amd64, arm64)

### Key Capabilities

#### 🎯 **Declarative Configuration**
Define NFS servers using familiar Kubernetes YAML manifests with simple, intuitive specifications.

#### 🔄 **Automated Lifecycle Management**
Complete automation of NFS server deployment, scaling, updates, and cleanup operations.

#### 💾 **Flexible Storage Options**
Support for both StorageClass-based dynamic provisioning and pre-existing PersistentVolume binding.

#### 🏗️ **High Availability**
Configurable replica count for redundancy and improved availability of NFS services.

#### 🔍 **Status Monitoring**
Real-time status updates, health checks, and comprehensive observability features.

#### 🛡️ **Security First**
Security-focused design with proper RBAC, network policies, and container security practices.

## Quick Start

### Prerequisites

- **Kubernetes cluster** (v1.20+)
- **kubectl** configured to access your cluster
- **Cluster admin permissions**

> **📝 Note**: cert-manager is only required if you plan to enable webhook validation or secure metrics with TLS certificates. For basic NFS Server Controller functionality, cert-manager is not needed.

### Installation

1. **Install the CRDs and operator:**
   ```bash
   kubectl apply -f https://github.com/sharedvolume/nfs-server-controller/releases/latest/download/install.yaml
   ```

   Or using kustomize:
   ```bash
   kubectl apply -k config/default
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

## Development

### Prerequisites

- **Go 1.24+**: The project requires Go 1.24 or later
- **Docker**: For building container images
- **kubectl**: For interacting with Kubernetes
- **Kind**: For local development and testing
- **Make**: For build automation

### Building from Source

1. **Clone the repository:**
   ```bash
   git clone https://github.com/sharedvolume/nfs-server-controller.git
   cd nfs-server-controller
   ```

2. **Install dependencies:**
   ```bash
   go mod download
   ```

3. **Install Kubebuilder tools:**
   ```bash
   make controller-gen kustomize
   ```

4. **Build the manager:**
   ```bash
   make build
   ```

5. **Run tests:**
   ```bash
   make test
   ```

6. **Build Docker image:**
   ```bash
   make docker-build IMG=nfs-server-controller:dev
   ```

### Architecture Overview

The NFS Server Controller is built using the [Kubebuilder](https://kubebuilder.io/) framework, providing a comprehensive toolkit for building Kubernetes operators and controllers in Go.

#### Project Structure

```
nfs-server-controller/
├── api/                    # API definitions (CRDs)
│   └── v1alpha1/
├── cmd/                    # Main application
├── config/                 # Kubernetes manifests
│   ├── crd/               # Custom Resource Definitions
│   ├── rbac/              # RBAC permissions
│   └── samples/           # Sample configurations
├── internal/               # Internal packages
│   └── controller/        # Controller logic
└── test/                  # Test suites
```

#### Technology Stack

- **Language**: Go 1.24+
- **Framework**: Kubebuilder/controller-runtime
- **Container**: Distroless base images
- **Storage**: Kubernetes PersistentVolumes
- **Networking**: Kubernetes Services
- **Security**: RBAC, Pod Security Standards

### Local Development

1. **Install CRDs:**
   ```bash
   make install
   ```

2. **Run the controller locally:**
   ```bash
   make run
   ```

3. **Apply sample configurations:**
   ```bash
   kubectl apply -f config/samples/
   ```

4. **Run tests with Kind:**
   ```bash
   make test-e2e
   ```

### API Reference

#### NfsServer Spec

| Field | Type | Description | Required | Default |
|-------|------|-------------|----------|---------|
| `storage.capacity` | string | Storage capacity (e.g., "10Gi") | Yes | - |
| `storage.storageClassName` | string | StorageClass name for dynamic provisioning | No* | - |
| `storage.persistentVolume` | string | Pre-existing PersistentVolume name | No* | - |
| `replicas` | *int32 | Number of NFS server replicas | No | 2 |
| `path` | string | NFS export path | No | "/nfs" |
| `image` | string | Container image for NFS server | No | Auto-detected |

*Either `storageClassName` or `persistentVolume` must be specified, but not both.

### Testing Strategy

#### Unit Tests
```bash
# Run unit tests
make test

# Run with coverage
go test ./... -coverprofile cover.out
go tool cover -html=cover.out
```

#### Integration Tests
```bash
# Run end-to-end tests
make test-e2e
```

### Deployment

```bash
# Deploy to cluster
make deploy IMG=sharedvolume/nfs-server-controller:latest

# Undeploy
make undeploy
```

### Troubleshooting

#### Common Issues

**Controller Not Starting:**
```bash
# Check deployment status
kubectl describe deployment -n nfs-server-controller-system nfs-server-controller-manager

# Check controller logs
kubectl logs -n nfs-server-controller-system deployment/nfs-server-controller-manager -f
```

**NFS Server Not Ready:**
```bash
# Check NFS server status
kubectl describe nfsserver <name>

# Check related resources
kubectl get pods,svc,pvc -l app=<nfs-server-name>
```

**Mount Issues:**
```bash
# Test connectivity
kubectl run nfs-test --image=alpine:latest --rm -it -- /bin/sh
# Inside pod: apk add nfs-utils && showmount -e <service-name>.<namespace>.svc.cluster.local
```

#### Debug Mode
```bash
# Enable debug logging
kubectl patch deployment -n nfs-server-controller-system nfs-server-controller-manager \
  -p '{"spec":{"template":{"spec":{"containers":[{"name":"manager","args":["--log-level=debug"]}]}}}}'
```

## Community

### Contributing

We welcome contributions from the community! Whether it's:

- 🐛 **Bug reports and fixes**
- ✨ **New features and enhancements**
- 📚 **Documentation improvements**
- 🧪 **Testing and quality assurance**
- 💡 **Ideas and suggestions**

#### How to Contribute

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass
6. Submit a pull request

#### Contributing Guidelines

See our [Contributing Guidelines](CONTRIBUTING.md) for details on:

- Code of conduct
- Development setup
- Pull request process
- Testing requirements

### Support and Community

- **GitHub**: [sharedvolume/nfs-server-controller](https://github.com/sharedvolume/nfs-server-controller)
- **Issues**: Report bugs and feature requests on [GitHub Issues](https://github.com/sharedvolume/nfs-server-controller/issues)
- **Discussions**: Join community discussions on [GitHub Discussions](https://github.com/sharedvolume/nfs-server-controller/discussions)
- **Security**: Report security vulnerabilities privately to bilgehan.nal@gmail.com

### Acknowledgments

Built with ❤️ using open source technologies and inspired by the Kubernetes community's best practices for operators and controllers.

## Use Cases

### **Development Teams**
- Shared development environments
- Code repositories and build artifacts
- Temporary storage for CI/CD pipelines
- Cross-team collaboration spaces

### **Data Analytics**
- Shared datasets for ML/AI workloads
- Data lakes and warehouses
- ETL pipeline intermediate storage
- Research data sharing

### **Enterprise Applications**
- Legacy application integration
- Shared configuration and templates
- Backup and archive storage
- Multi-tenant shared storage

### **Enterprise Production Environments**
- Mission-critical shared storage for business applications
- High-availability storage for enterprise workloads
- Compliance-ready storage solutions with audit trails
- Multi-tenant environments with proper isolation

### **DevOps and Platform Engineering**
- Infrastructure as Code storage templates
- Centralized monitoring and observability data storage
- Configuration management at enterprise scale
- Business continuity and disaster recovery scenarios

## Architecture

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   kubectl/API  │    │  NFS Controller  │    │   NFS Server    │
│                 │───▶│                  │───▶│     Pods        │
│ apply nfs.yaml  │    │  Reconcile Loop  │    │                 │
└─────────────────┘    └──────────────────┘    └─────────────────┘
                                │                         │
                                ▼                         ▼
                       ┌──────────────────┐    ┌─────────────────┐
                       │       PVC        │    │    Service      │
                       │   (Storage)      │    │  (Discovery)    │
                       └──────────────────┘    └─────────────────┘
```

The NFS Server Controller consists of:

- **Custom Resource Definition (CRD)**: Defines the `NfsServer` resource schema
- **Controller**: Watches for `NfsServer` resources and manages their lifecycle
- **Reconciler**: Ensures the desired state matches the actual state by creating/updating:
  - PersistentVolumeClaims for storage
  - ReplicaSets for NFS server pods
  - Services for network access

### Technology Stack

- **Language**: Go 1.24+
- **Framework**: Kubebuilder/controller-runtime
- **Container**: Distroless base images
- **Storage**: Kubernetes PersistentVolumes
- **Networking**: Kubernetes Services
- **Security**: RBAC, Pod Security Standards

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

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## About SharedVolume

The NFS Server Controller is part of the SharedVolume ecosystem, a comprehensive suite of storage solutions for Kubernetes environments. SharedVolume provides enterprise-grade storage orchestration tools designed for production workloads at scale.

## Support

- **Issues**: Report bugs and feature requests on [GitHub Issues](https://github.com/sharedvolume/nfs-server-controller/issues)
- **Discussions**: Join community discussions on [GitHub Discussions](https://github.com/sharedvolume/nfs-server-controller/discussions)
- **Security**: Report security vulnerabilities privately to bilgehan.nal@gmail.com

## Acknowledgments

Built with [Kubebuilder](https://github.com/kubernetes-sigs/kubebuilder) and inspired by the Kubernetes community's best practices for operators.

---

<div align="center">

**⭐ Star this repository if it helped you!**

[![GitHub stars](https://img.shields.io/github/stars/sharedvolume/nfs-server-controller?style=social)](https://github.com/sharedvolume/nfs-server-controller/stargazers)
[![GitHub forks](https://img.shields.io/github/forks/sharedvolume/nfs-server-controller?style=social)](https://github.com/sharedvolume/nfs-server-controller/network/members)
[![GitHub watchers](https://img.shields.io/github/watchers/sharedvolume/nfs-server-controller?style=social)](https://github.com/sharedvolume/nfs-server-controller/watchers)

[Website](https://sharedvolume.github.io) • [Docker Hub](https://hub.docker.com/r/sharedvolume/nfs-server-controller) • [Contributing](CONTRIBUTING.md)

</div>