# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial open source release
- Comprehensive documentation and examples
- GitHub Actions for CI/CD
- Security policy and contributing guidelines

## [0.0.11] - 2025-08-03

### Added
- NFS Server Controller with custom resource definition
- Support for dynamic storage provisioning via StorageClass
- Support for pre-existing PersistentVolume binding
- Configurable replica count for high availability
- Automatic service creation for NFS server discovery
- Real-time status monitoring and health checks
- Comprehensive API validation
- Multi-platform Docker image support (amd64, arm64)

### Features
- **Custom Resource Definition (CRD)**: Define NFS servers declaratively
- **Dynamic Provisioning**: Automatic PVC creation with specified storage
- **Lifecycle Management**: Complete resource lifecycle handling
- **Storage Flexibility**: Support for both StorageClass and PV approaches
- **High Availability**: Configurable replica count (default: 2)
- **Service Discovery**: Automatic Kubernetes service creation
- **Status Monitoring**: Real-time ready status and phase tracking

### Technical Details
- Built with Kubebuilder framework
- Go 1.24+ support
- Kubernetes 1.20+ compatibility
- MIT License
- Multi-architecture container images
- Comprehensive test suite
- Security-focused implementation

### Documentation
- Complete API reference documentation
- Installation guide with multiple methods
- Usage examples and best practices
- Security guidelines and considerations
- Contributing guidelines for developers
- Troubleshooting guide

### Container Images
- `sharedvolume/nfs-server-controller:0.0.11`
- `sharedvolume/nfs-server-controller:latest`
- Multi-platform support: linux/amd64, linux/arm64

### Dependencies
- kubernetes v0.33.0
- controller-runtime v0.21.0
- Go 1.24+

---

## Release Notes Template

### [Version] - YYYY-MM-DD

### Added
- New features

### Changed
- Changes in existing functionality

### Deprecated
- Soon-to-be removed features

### Removed
- Removed features

### Fixed
- Bug fixes

### Security
- Security improvements
