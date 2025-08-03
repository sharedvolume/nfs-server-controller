# API Reference

## NfsServer

NfsServer is a custom resource that defines an NFS server instance.

### Spec

#### NfsServerSpec

| Field | Type | Description | Required | Default |
|-------|------|-------------|----------|---------|
| `storage` | [StorageSpec](#storagespec) | Storage configuration for the NFS server | Yes | - |
| `replicas` | *int32 | Number of NFS server replicas | No | 2 |
| `path` | string | NFS export path | No | "/nfs" |
| `image` | string | Container image for NFS server | No | "sharedvolume/nfs-server:0.0.3-alpine-3.22.0" |
| `address` | string | NFS server address (auto-generated) | No | - |

#### StorageSpec

| Field | Type | Description | Required |
|-------|------|-------------|----------|
| `capacity` | string | Storage capacity (e.g., "10Gi") | Yes |
| `storageClassName` | string | StorageClass name for dynamic provisioning | No* |
| `persistentVolume` | string | Pre-existing PersistentVolume name | No* |

*Either `storageClassName` or `persistentVolume` must be specified, but not both.

### Status

#### NfsServerStatus

| Field | Type | Description |
|-------|------|-------------|
| `ready` | bool | Whether the NFS server is ready to serve requests |
| `phase` | string | Current phase of the NFS server (Running, Pending, Error) |
| `message` | string | Human-readable message about the current status |

### Examples

#### Basic NFS Server with StorageClass

```yaml
apiVersion: sharedvolume.io/v1alpha1
kind: NfsServer
metadata:
  name: basic-nfs
  namespace: default
spec:
  storage:
    capacity: "10Gi"
    storageClassName: "standard"
```

#### NFS Server with Custom Configuration

```yaml
apiVersion: sharedvolume.io/v1alpha1
kind: NfsServer
metadata:
  name: custom-nfs
  namespace: storage
spec:
  storage:
    capacity: "50Gi"
    storageClassName: "fast-ssd"
  replicas: 3
  path: "/exports"
  image: "sharedvolume/nfs-server:custom"
```

#### NFS Server with Pre-existing PV

```yaml
apiVersion: sharedvolume.io/v1alpha1
kind: NfsServer
metadata:
  name: pv-nfs
  namespace: default
spec:
  storage:
    capacity: "100Gi"
    persistentVolume: "my-existing-pv"
  replicas: 1
```

## Status Values

### Phase

- `Running`: NFS server is running and ready
- `Pending`: NFS server is being created or waiting for resources
- `Error`: An error occurred during NFS server creation or operation

### Ready

- `true`: NFS server is ready to serve requests
- `false`: NFS server is not ready (see `message` for details)

## Controller Behavior

The NFS Server Controller watches for `NfsServer` resources and:

1. **Validates** the resource specification
2. **Creates** a PersistentVolumeClaim for storage
3. **Deploys** a ReplicaSet with NFS server pods
4. **Exposes** the NFS server via a Kubernetes Service
5. **Updates** the status based on resource readiness
6. **Manages** the lifecycle of all owned resources

### Owned Resources

For each `NfsServer`, the controller creates:

- **PersistentVolumeClaim**: For NFS data storage
- **ReplicaSet**: For NFS server pod management
- **Service**: For NFS server network access

All owned resources are automatically cleaned up when the `NfsServer` is deleted.

## RBAC Permissions

The controller requires the following permissions:

```yaml
# NfsServer resources
- apiGroups: ["sharedvolume.io"]
  resources: ["nfsservers", "nfsservers/status", "nfsservers/finalizers"]
  verbs: ["get", "list", "watch", "create", "update", "patch", "delete"]

# Owned resources
- apiGroups: [""]
  resources: ["persistentvolumeclaims", "services"]
  verbs: ["get", "list", "watch", "create", "update", "patch", "delete"]

- apiGroups: ["apps"]
  resources: ["replicasets"]
  verbs: ["get", "list", "watch", "create", "update", "patch", "delete"]

# Read-only for validation
- apiGroups: [""]
  resources: ["persistentvolumes"]
  verbs: ["list", "watch"]

- apiGroups: ["storage.k8s.io"]
  resources: ["storageclasses"]
  verbs: ["list", "watch"]
```
