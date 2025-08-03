# Troubleshooting Guide

This guide helps you diagnose and resolve common issues with the NFS Server Controller.

## Quick Diagnosis

### Check Controller Status

```bash
# Check if the controller is running
kubectl get deployment -n nfs-server-controller-system

# Check controller logs
kubectl logs -n nfs-server-controller-system deployment/nfs-server-controller-manager -f

# Check if CRDs are installed
kubectl get crd nfsservers.sharedvolume.io
```

### Check NFS Server Status

```bash
# List all NFS servers
kubectl get nfsservers -A

# Get detailed status
kubectl describe nfsserver <name> -n <namespace>

# Check related resources
kubectl get pods,svc,pvc -l app=<nfs-server-name>
```

## Common Issues

### 1. Controller Not Starting

#### Symptoms
- Controller deployment shows 0/1 ready
- Error events in controller deployment
- No logs from controller pods

#### Diagnosis
```bash
# Check deployment status
kubectl describe deployment -n nfs-server-controller-system nfs-server-controller-manager

# Check pod events
kubectl get events -n nfs-server-controller-system --sort-by='.lastTimestamp'

# Check RBAC permissions
kubectl auth can-i create nfsservers --as=system:serviceaccount:nfs-server-controller-system:nfs-server-controller-manager
```

#### Solutions

**Missing CRDs:**
```bash
# Reinstall CRDs
kubectl apply -f https://github.com/sharedvolume/nfs-server-controller/releases/latest/download/install.yaml
```

**RBAC Issues:**
```bash
# Check and fix RBAC
kubectl describe clusterrole nfs-server-controller-manager-role
kubectl describe clusterrolebinding nfs-server-controller-manager-rolebinding
```

**Image Pull Issues:**
```bash
# Check image pull secrets
kubectl describe pod -n nfs-server-controller-system -l control-plane=controller-manager
```

### 2. NFS Server Not Ready

#### Symptoms
- NfsServer resource shows `Ready: false`
- Phase stuck in "Pending" or "Error"
- NFS server pods not running

#### Diagnosis
```bash
# Check NFS server status
kubectl describe nfsserver <name>

# Check related pods
kubectl get pods -l app=<nfs-server-name>
kubectl describe pod -l app=<nfs-server-name>

# Check PVC status
kubectl get pvc <nfs-server-name>
kubectl describe pvc <nfs-server-name>
```

#### Solutions

**PVC Not Bound:**
```bash
# Check storage class
kubectl get storageclass
kubectl describe storageclass <storage-class-name>

# Check PV availability (if using specific PV)
kubectl get pv
kubectl describe pv <pv-name>
```

**Pod Scheduling Issues:**
```bash
# Check node resources
kubectl describe nodes

# Check if privileged containers are allowed
kubectl get psp  # If using Pod Security Policies
```

**Container Issues:**
```bash
# Check container logs
kubectl logs -l app=<nfs-server-name> -c nfs-server

# Check container security context
kubectl get pod -l app=<nfs-server-name> -o yaml | grep -A 10 securityContext
```

### 3. NFS Mount Issues

#### Symptoms
- Clients cannot mount NFS share
- Mount hangs or times out
- Permission denied errors

#### Diagnosis
```bash
# Test NFS server connectivity
kubectl run nfs-test --image=alpine:latest --rm -it -- /bin/sh
# Inside the pod:
apk add nfs-utils
showmount -e <nfs-server-service>.<namespace>.svc.cluster.local

# Check NFS server service
kubectl get svc <nfs-server-name>
kubectl describe svc <nfs-server-name>

# Check endpoints
kubectl get endpoints <nfs-server-name>
```

#### Solutions

**Service Discovery Issues:**
```bash
# Verify service exists and has endpoints
kubectl get svc <nfs-server-name>
kubectl get endpoints <nfs-server-name>

# Test DNS resolution
kubectl run dns-test --image=busybox:1.28 --rm -it -- nslookup <nfs-server-name>.<namespace>.svc.cluster.local
```

**Network Policy Issues:**
```bash
# Check network policies
kubectl get networkpolicy -A

# Temporarily disable network policies for testing
kubectl annotate networkpolicy <policy-name> policy.networking.kubernetes.io/disabled=true
```

**NFS Client Issues:**
```bash
# Ensure NFS utils are installed in client pods
# Add to Dockerfile: RUN apt-get update && apt-get install -y nfs-common
```

### 4. Storage Issues

#### Symptoms
- PVC stuck in Pending state
- Storage capacity errors
- Data not persisting

#### Diagnosis
```bash
# Check PVC status
kubectl describe pvc <nfs-server-name>

# Check storage class
kubectl describe storageclass <storage-class-name>

# Check available storage
kubectl get pv
```

#### Solutions

**Storage Class Issues:**
```bash
# List available storage classes
kubectl get storageclass

# Use existing storage class
kubectl patch nfsserver <name> -p '{"spec":{"storage":{"storageClassName":"<valid-class>"}}}'
```

**Insufficient Resources:**
```bash
# Check cluster capacity
kubectl describe nodes | grep -A 5 "Allocated resources"

# Reduce requested capacity
kubectl patch nfsserver <name> -p '{"spec":{"storage":{"capacity":"<smaller-size>"}}}'
```

### 5. Performance Issues

#### Symptoms
- Slow NFS operations
- High latency
- Timeouts

#### Diagnosis
```bash
# Check resource usage
kubectl top pods -l app=<nfs-server-name>
kubectl describe pod -l app=<nfs-server-name>

# Check storage performance
kubectl exec -it <nfs-server-pod> -- df -h
kubectl exec -it <nfs-server-pod> -- iostat -x 1 5
```

#### Solutions

**Resource Limits:**
```bash
# Increase resource limits (requires controller update)
# Edit the controller to allow resource customization
```

**Storage Performance:**
```bash
# Use faster storage class
kubectl patch nfsserver <name> -p '{"spec":{"storage":{"storageClassName":"fast-ssd"}}}'
```

**Network Optimization:**
```bash
# Use higher performance NFS options in client mounts
# Add mount options: rsize=1048576,wsize=1048576,hard,intr,timeo=600
```

## Advanced Troubleshooting

### Debug Mode

Enable debug logging in the controller:

```bash
kubectl patch deployment -n nfs-server-controller-system nfs-server-controller-manager -p '{"spec":{"template":{"spec":{"containers":[{"name":"manager","args":["--log-level=debug"]}]}}}}'
```

### Controller Metrics

If metrics are enabled:

```bash
# Port forward to metrics endpoint
kubectl port-forward -n nfs-server-controller-system deployment/nfs-server-controller-manager 8080:8080

# Query metrics
curl http://localhost:8080/metrics | grep nfs_server
```

### Resource Debugging

```bash
# Get full resource definition
kubectl get nfsserver <name> -o yaml

# Check owner references
kubectl get replicaset,pvc,service -o wide | grep <nfs-server-name>

# Trace events
kubectl get events --sort-by='.lastTimestamp' | grep <nfs-server-name>
```

### NFS Server Debugging

```bash
# Access NFS server pod
kubectl exec -it <nfs-server-pod> -- /bin/sh

# Check NFS processes
ps aux | grep nfs

# Check NFS exports
cat /etc/exports
exportfs -v

# Check mount points
mount | grep nfs
df -h

# Check NFS logs
dmesg | grep -i nfs
journalctl -u nfs-server
```

## Recovery Procedures

### Controller Recovery

```bash
# Restart controller
kubectl rollout restart deployment -n nfs-server-controller-system nfs-server-controller-manager

# Reinstall controller
kubectl delete -f https://github.com/sharedvolume/nfs-server-controller/releases/latest/download/install.yaml
kubectl apply -f https://github.com/sharedvolume/nfs-server-controller/releases/latest/download/install.yaml
```

### NFS Server Recovery

```bash
# Force recreation of NFS server
kubectl delete nfsserver <name>
kubectl apply -f <nfs-server-manifest.yaml>

# Restart NFS server pods
kubectl delete pod -l app=<nfs-server-name>
```

### Data Recovery

```bash
# If data is lost, check PV reclaim policy
kubectl get pv -o custom-columns=NAME:.metadata.name,RECLAIM:.spec.persistentVolumeReclaimPolicy

# Recover from backup (if backup was configured)
# This depends on your backup solution
```

## Prevention

### Monitoring

Set up monitoring for:
- Controller health
- NFS server resource usage
- Storage capacity
- Network connectivity

### Backup

Implement backup strategies for:
- NFS server configurations
- Data stored on NFS shares
- Controller configuration

### Testing

Regular testing should include:
- NFS server creation/deletion
- Client mount/unmount operations
- Failover scenarios
- Performance benchmarks

## Getting Help

If you can't resolve the issue:

1. **Search existing issues**: [GitHub Issues](https://github.com/sharedvolume/nfs-server-controller/issues)
2. **Create a new issue**: Include logs, resource definitions, and troubleshooting steps tried
3. **Community discussion**: [GitHub Discussions](https://github.com/sharedvolume/nfs-server-controller/discussions)

### Information to Include

When reporting issues, include:

```bash
# System information
kubectl version
kubectl get nodes -o wide

# Controller information
kubectl get deployment -n nfs-server-controller-system -o yaml
kubectl logs -n nfs-server-controller-system deployment/nfs-server-controller-manager --tail=100

# NFS server information
kubectl get nfsserver <name> -o yaml
kubectl describe nfsserver <name>
kubectl get pods,svc,pvc -l app=<nfs-server-name> -o wide

# Events
kubectl get events --sort-by='.lastTimestamp' | tail -20
```
