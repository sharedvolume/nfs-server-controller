# Development Guide

This guide covers the development setup and practices for the NFS Server Controller project.

## Architecture Overview

The NFS Server Controller is built using the [Kubebuilder](https://kubebuilder.io/) framework, which provides a comprehensive toolkit for building Kubernetes operators and controllers in Go.

### Kubebuilder Foundation

This project was bootstrapped with Kubebuilder v4, which provides:

- **Controller Runtime**: Built on top of [controller-runtime](https://github.com/kubernetes-sigs/controller-runtime)
- **Code Generation**: Automatic generation of client code, deepcopy methods, and CRDs
- **Testing Framework**: Integration testing with Ginkgo and Gomega
- **Makefile**: Comprehensive build, test, and deployment targets
- **Scaffolding**: Well-structured project layout following Kubernetes best practices

### Project Structure

```
nfs-server-controller/
├── api/                    # API definitions (CRDs)
│   └── v1alpha1/
│       ├── groupversion_info.go
│       ├── nfsserver_types.go
│       └── zz_generated.deepcopy.go
├── cmd/                    # Main application
│   └── main.go
├── config/                 # Kubernetes manifests
│   ├── crd/               # Custom Resource Definitions
│   ├── default/           # Default Kustomize configuration
│   ├── manager/           # Controller manager deployment
│   ├── rbac/              # RBAC permissions
│   └── samples/           # Sample configurations
├── internal/               # Internal packages
│   └── controller/        # Controller logic
├── test/                  # Test suites
└── docs/                  # Documentation
```

## Development Setup

### Prerequisites

- **Go 1.24+**: The project requires Go 1.24 or later
- **Docker**: For building container images
- **kubectl**: For interacting with Kubernetes
- **Kind**: For local development and testing
- **Make**: For build automation

### Getting Started

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

4. **Run tests:**
   ```bash
   make test
   ```

5. **Build the manager:**
   ```bash
   make build
   ```

## Kubebuilder Commands Used

The project was created and maintained using these Kubebuilder commands:

### Initial Project Creation
```bash
# Initialize the project
kubebuilder init --domain sharedvolume.io --repo github.com/sharedvolume/nfs-server-controller

# Create the API and controller
kubebuilder create api --group nfs --version v1alpha1 --kind NfsServer --resource --controller
```

### Adding Features
```bash
# Generate manifests and code
make manifests generate

# Update CRDs in the cluster
make install

# Run the controller locally
make run
```

### Common Development Tasks

#### Code Generation

Kubebuilder automatically generates code based on markers and annotations:

```bash
# Generate deepcopy methods, client code, and CRDs
make generate manifests
```

Key markers used in the project:

- `//+kubebuilder:object:root=true` - Marks the root object
- `//+kubebuilder:subresource:status` - Enables status subresource
- `//+kubebuilder:printcolumn` - Defines kubectl output columns
- `//+kubebuilder:rbac` - Generates RBAC permissions

#### Testing

The project uses Ginkgo and Gomega for testing:

```bash
# Run unit tests
make test

# Run with coverage
go test ./... -coverprofile cover.out
go tool cover -html=cover.out

# Run end-to-end tests
make test-e2e
```

#### Local Development

1. **Install CRDs:**
   ```bash
   make install
   ```

2. **Run controller locally:**
   ```bash
   make run
   ```

3. **Apply sample configurations:**
   ```bash
   kubectl apply -f config/samples/
   ```

## Controller Development Patterns

### Reconciliation Logic

The controller follows the standard Kubernetes controller pattern:

```go
func (r *NfsServerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
    // 1. Fetch the resource
    // 2. Handle deletion (if needed)
    // 3. Validate the specification
    // 4. Create/Update owned resources
    // 5. Update status
    // 6. Return result
}
```

### Owner References

All managed resources use owner references for automatic cleanup:

```go
ctrl.SetControllerReference(nfsServer, ownedResource, r.Scheme)
```

### Status Management

Status updates are handled separately from spec changes:

```go
err = retry.RetryOnConflict(retry.DefaultRetry, func() error {
    latest := &nfsv1alpha1.NfsServer{}
    if err := r.Get(ctx, req.NamespacedName, latest); err != nil {
        return err
    }
    latest.Status = updatedStatus
    return r.Status().Update(ctx, latest)
})
```

## Building and Deployment

### Local Build

```bash
# Build the binary
make build

# Build Docker image
make docker-build IMG=nfs-server-controller:dev

# Load into Kind cluster
kind load docker-image nfs-server-controller:dev
```

### Deployment

```bash
# Deploy to cluster
make deploy IMG=sharedvolume/nfs-server-controller:latest

# Undeploy
make undeploy
```

### Release Process

The project uses GitHub Actions for automated releases:

1. **Tag the release:**
   ```bash
   git tag v0.1.0
   git push origin v0.1.0
   ```

2. **GitHub Actions will:**
   - Build multi-architecture images
   - Push to Docker Hub
   - Generate installation manifests
   - Create GitHub release

## API Design Principles

### Following Kubernetes Conventions

The API design follows Kubernetes API conventions:

- **Spec/Status separation**: Clear distinction between desired and observed state
- **Declarative**: Users declare what they want, controller figures out how
- **Idempotent**: Multiple applications of the same spec produce the same result
- **Extensible**: API can evolve without breaking existing users

### Custom Resource Definition

```yaml
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: nfsservers.nfs.sharedvolume.io
spec:
  group: nfs.sharedvolume.io
  versions:
  - name: v1alpha1
    served: true
    storage: true
    schema:
      # OpenAPI schema definition
```

### Validation

API validation is enforced through:
- OpenAPI schema in CRD
- Webhook validation (future enhancement)
- Controller-side validation

## Testing Strategy

### Unit Tests

- Controller logic testing with fake clients
- API validation testing
- Utility function testing

### Integration Tests

- Full controller testing with envtest
- CRD creation and management
- Resource reconciliation testing

### End-to-End Tests

- Real cluster testing with Kind
- Complete workflow validation
- Performance and stress testing

## Debugging

### Controller Logs

```bash
# View controller logs
kubectl logs -n nfs-server-controller-system deployment/nfs-server-controller-manager -f

# Enable debug logging
kubectl patch deployment -n nfs-server-controller-system nfs-server-controller-manager \
  -p '{"spec":{"template":{"spec":{"containers":[{"name":"manager","args":["--log-level=debug"]}]}}}}'
```

### Resource Investigation

```bash
# Check resource status
kubectl describe nfsserver <name>

# Check events
kubectl get events --sort-by='.lastTimestamp'

# Check owned resources
kubectl get pods,svc,pvc -l app=<nfs-server-name>
```

## Contributing to Development

### Code Style

- Follow Go conventions and best practices
- Use meaningful variable and function names
- Add comments for exported functions
- Keep functions focused and small

### Commit Messages

Follow conventional commits:
```
feat: add support for custom NFS export options
fix: resolve PVC binding issue in multi-zone clusters
docs: update API reference documentation
```

### Pull Request Process

1. Fork and create feature branch
2. Make changes with tests
3. Update documentation
4. Submit pull request
5. Address review feedback

## Useful Resources

### Kubebuilder Documentation
- [Kubebuilder Book](https://book.kubebuilder.io/)
- [Controller Runtime](https://pkg.go.dev/sigs.k8s.io/controller-runtime)
- [Kubebuilder Tutorial](https://book.kubebuilder.io/cronjob-tutorial/cronjob-tutorial.html)

### Kubernetes Development
- [Kubernetes API Conventions](https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md)
- [Writing Controllers](https://github.com/kubernetes/community/blob/master/contributors/devel/sig-api-machinery/controllers.md)
- [Operator Best Practices](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/)

### Testing and Quality
- [Ginkgo Testing Framework](https://onsi.github.io/ginkgo/)
- [Gomega Matcher Library](https://onsi.github.io/gomega/)
- [Go Testing Best Practices](https://github.com/golang/go/wiki/TestComments)

This development guide provides a foundation for contributing to the NFS Server Controller project using Kubebuilder and following Kubernetes best practices.
