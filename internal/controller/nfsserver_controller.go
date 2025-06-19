/*
Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file exeptn compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"errors"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/util/retry"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	nfsv1alpha1 "github.com/sharedvolume/nfs-server-controller/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NfsServerReconciler reconciles a NfsServer object
type NfsServerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=nfs.sharedvolume.io,resources=nfsservers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=nfs.sharedvolume.io,resources=nfsservers/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=nfs.sharedvolume.io,resources=nfsservers/finalizers,verbs=update
// +kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=replicasets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=persistentvolumeclaims,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=persistentvolumes,verbs=list;watch
// +kubebuilder:rbac:groups=storage.k8s.io,resources=storageclasses,verbs=list;watch
// +kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the NfsServer object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.21.0/pkg/reconcile
func (r *NfsServerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	// 1. Fetch the NfsServer object
	nfsServer := &nfsv1alpha1.NfsServer{}
	err := r.Get(ctx, req.NamespacedName, nfsServer)
	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Handle deletion
	if !nfsServer.ObjectMeta.DeletionTimestamp.IsZero() {
		// If the NfsServer is being deleted, do not attempt to recreate resources
		return ctrl.Result{}, nil
	}

	// 2. Validate storage fields
	storage := nfsServer.Spec.Storage
	if (storage.StorageClassName != "" && storage.PersistentVolume != "") || (storage.StorageClassName == "" && storage.PersistentVolume == "") {
		retry.RetryOnConflict(retry.DefaultRetry, func() error {
			latest := &nfsv1alpha1.NfsServer{}
			if err := r.Get(ctx, req.NamespacedName, latest); err != nil {
				return err
			}
			latest.Status.Ready = false
			latest.Status.Phase = "Error"
			latest.Status.Message = "Exactly one of storageClassName or persistentVolume must be set."
			return r.Status().Update(ctx, latest)
		})
		return ctrl.Result{}, errors.New("Exactly one of storageClassName or persistentVolume must be set.")
	}

	// 3. Generate PVC
	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      nfsServer.Name,
			Namespace: nfsServer.Namespace,
		},
	}
	_, err = ctrl.CreateOrUpdate(ctx, r.Client, pvc, func() error {
		pvc.Spec.AccessModes = []corev1.PersistentVolumeAccessMode{corev1.ReadWriteMany}
		pvc.Spec.Resources.Requests = corev1.ResourceList{
			corev1.ResourceStorage: resourceMustParse(storage.Capacity),
		}
		if storage.StorageClassName != "" {
			pvc.Spec.StorageClassName = &storage.StorageClassName
		} else {
			pvc.Spec.StorageClassName = nil
			pvc.Spec.VolumeName = storage.PersistentVolume
		}
		return ctrl.SetControllerReference(nfsServer, pvc, r.Scheme)
	})
	if err != nil {
		log.Error(err, "Failed to create/update PVC")
		nfsServer.Status.Ready = false
		nfsServer.Status.Phase = "Error"
		nfsServer.Status.Message = fmt.Sprintf("PVC error: %v", err)
		r.Status().Update(ctx, nfsServer)
		return ctrl.Result{}, err
	}

	// 4. Generate ReplicaSet
	rs := makeNfsReplicaSet(nfsServer)
	err = r.CreateOrUpdateReplicaSet(ctx, rs, nfsServer)
	if err != nil {
		log.Error(err, "Failed to create/update ReplicaSet")
		nfsServer.Status.Ready = false
		nfsServer.Status.Phase = "Error"
		nfsServer.Status.Message = fmt.Sprintf("ReplicaSet error: %v", err)
		r.Status().Update(ctx, nfsServer)
		return ctrl.Result{}, err
	}

	// 5. Generate Service
	svc := makeNfsService(nfsServer)
	err = r.CreateOrUpdateService(ctx, svc, nfsServer)
	if err != nil {
		log.Error(err, "Failed to create/update Service")
		nfsServer.Status.Ready = false
		nfsServer.Status.Phase = "Error"
		nfsServer.Status.Message = fmt.Sprintf("Service error: %v", err)
		r.Status().Update(ctx, nfsServer)
		return ctrl.Result{}, err
	}

	// 6. Status: check readiness
	ready, phase, msg := r.CheckNfsServerReady(ctx, nfsServer)
	nfsServer.Status.Ready = ready
	nfsServer.Status.Phase = phase
	nfsServer.Status.Message = msg

	// 7. Address
	address := fmt.Sprintf("%s.%s.svc.cluster.local", nfsServer.Name, nfsServer.Namespace)
	if nfsServer.Spec.Address != address {
		retry.RetryOnConflict(retry.DefaultRetry, func() error {
			latest := &nfsv1alpha1.NfsServer{}
			if err := r.Get(ctx, req.NamespacedName, latest); err != nil {
				return err
			}
			latest.Spec.Address = address
			return r.Update(ctx, latest)
		})
	}

	err = retry.RetryOnConflict(retry.DefaultRetry, func() error {
		latest := &nfsv1alpha1.NfsServer{}
		if err := r.Get(ctx, req.NamespacedName, latest); err != nil {
			return err
		}
		latest.Status = nfsServer.Status
		return r.Status().Update(ctx, latest)
	})
	if err != nil {
		log.Error(err, "Failed to update NfsServer status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// Removed all logic related to adding or removing finalizers. Kubernetes garbage collector will handle owned resources if ownerReferences are set.

// Status phases and messages as constants
const (
	PhaseRunning = "Running"
	PhasePending = "Pending"
	PhaseError   = "Error"

	MsgDeploymentNotReady = "Deployment not ready"
	MsgPVCNotBound        = "PVC not bound"
	MsgServiceNotFound    = "Service not found"
	MsgNfsServerRunning   = "NfsServer is running"
	MsgStorageValidation  = "Exactly one of storageClassName or persistentVolume must be set."
	MsgPVCError           = "PVC error: %v"
	MsgDeploymentError    = "Deployment error: %v"
	MsgServiceError       = "Service error: %v"
)

// SetupWithManager sets up the controller with the Manager.
func (r *NfsServerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&nfsv1alpha1.NfsServer{}).
		Owns(&appsv1.ReplicaSet{}).
		Owns(&corev1.PersistentVolumeClaim{}).
		Owns(&corev1.Service{}).
		Named("nfsserver").
		Complete(r)
}

// Helper: parse resource quantity
func resourceMustParse(q string) resource.Quantity {
	res, err := resource.ParseQuantity(q)
	if err != nil {
		return resource.MustParse("1Gi") // fallback
	}
	return res
}

// Helper: create ReplicaSet spec
func makeNfsReplicaSet(nfsServer *nfsv1alpha1.NfsServer) *appsv1.ReplicaSet {
	replicas := int32(2)
	if nfsServer.Spec.Replicas != nil {
		replicas = *nfsServer.Spec.Replicas
	}
	image := nfsServer.Spec.Image
	if image == "" {
		image = "sharedvolume/nfs-server:alpine-3.22.0"
	}
	mountPath := nfsServer.Spec.Path
	if mountPath == "" {
		mountPath = "/nfs"
	}
	labels := map[string]string{"app": nfsServer.Name}
	return &appsv1.ReplicaSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      nfsServer.Name,
			Namespace: nfsServer.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.ReplicaSetSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{MatchLabels: labels},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: labels},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Name:            "nfs-server",
						Image:           image,
						SecurityContext: &corev1.SecurityContext{Privileged: pointerBool(true)},
						Env:             []corev1.EnvVar{{Name: "SHARED_DIRECTORY", Value: mountPath}},
						Args:            []string{mountPath},
						Ports:           []corev1.ContainerPort{{Name: "nfs", ContainerPort: 2049}, {Name: "mountd", ContainerPort: 20048}, {Name: "rpcbind", ContainerPort: 111}},
						VolumeMounts:    []corev1.VolumeMount{{Name: "nfs-data", MountPath: mountPath}},
						ReadinessProbe: &corev1.Probe{
							ProbeHandler: corev1.ProbeHandler{
								Exec: &corev1.ExecAction{
									Command: []string{"sh", "-c", "showmount -e localhost"},
								},
							},
							InitialDelaySeconds: 5,
							PeriodSeconds:       5,
						},
					}},
					Volumes: []corev1.Volume{{
						Name: "nfs-data",
						VolumeSource: corev1.VolumeSource{
							PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{ClaimName: nfsServer.Name},
						},
					}},
				},
			},
		},
	}
}

// Helper: create Service spec
func makeNfsService(nfsServer *nfsv1alpha1.NfsServer) *corev1.Service {
	labels := map[string]string{"app": nfsServer.Name}
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      nfsServer.Name,
			Namespace: nfsServer.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			ClusterIP: "None",
			Selector:  labels,
			Ports:     []corev1.ServicePort{{Name: "nfs", Port: 2049, TargetPort: intstrFromInt(2049)}, {Name: "mountd", Port: 20048, TargetPort: intstrFromInt(20048)}, {Name: "rpcbind", Port: 111, TargetPort: intstrFromInt(111)}},
		},
	}
}

// Helper: pointer to bool
func pointerBool(b bool) *bool { return &b }

// Helper: intstr from int
func intstrFromInt(i int) intstr.IntOrString {
	return intstr.IntOrString{Type: intstr.Int, IntVal: int32(i)}
}

// Helper: Create or Update ReplicaSet
func (r *NfsServerReconciler) CreateOrUpdateReplicaSet(ctx context.Context, rs *appsv1.ReplicaSet, owner metav1.Object) error {
	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		found := &appsv1.ReplicaSet{}
		err := r.Get(ctx, client.ObjectKey{Name: rs.Name, Namespace: rs.Namespace}, found)
		if err != nil {
			if client.IgnoreNotFound(err) != nil {
				return err
			}
			// Not found, so create
			ctrl.SetControllerReference(owner, rs, r.Scheme)
			return r.Create(ctx, rs)
		}
		ctrl.SetControllerReference(owner, rs, r.Scheme)
		rs.ResourceVersion = found.ResourceVersion
		return r.Update(ctx, rs)
	})
}

// Helper: Create or Update Service
func (r *NfsServerReconciler) CreateOrUpdateService(ctx context.Context, svc *corev1.Service, owner metav1.Object) error {
	found := &corev1.Service{}
	err := r.Get(ctx, client.ObjectKey{Name: svc.Name, Namespace: svc.Namespace}, found)
	if err != nil {
		if client.IgnoreNotFound(err) != nil {
			return err
		}
		// Not found, so create
		ctrl.SetControllerReference(owner, svc, r.Scheme)
		return r.Create(ctx, svc)
	}
	ctrl.SetControllerReference(owner, svc, r.Scheme)
	svc.ResourceVersion = found.ResourceVersion
	return r.Update(ctx, svc)
}

// Helper: Check readiness of ReplicaSet, PVC, Service
func (r *NfsServerReconciler) CheckNfsServerReady(ctx context.Context, nfsServer *nfsv1alpha1.NfsServer) (bool, string, string) {
	rs := &appsv1.ReplicaSet{}
	err := r.Get(ctx, client.ObjectKey{Name: nfsServer.Name, Namespace: nfsServer.Namespace}, rs)
	if err != nil || rs.Status.ReadyReplicas < 1 {
		return false, PhasePending, MsgDeploymentNotReady
	}
	pvc := &corev1.PersistentVolumeClaim{}
	err = r.Get(ctx, client.ObjectKey{Name: nfsServer.Name, Namespace: nfsServer.Namespace}, pvc)
	if err != nil || pvc.Status.Phase != corev1.ClaimBound {
		return false, PhasePending, MsgPVCNotBound
	}
	svc := &corev1.Service{}
	err = r.Get(ctx, client.ObjectKey{Name: nfsServer.Name, Namespace: nfsServer.Namespace}, svc)
	if err != nil {
		return false, PhasePending, MsgServiceNotFound
	}
	return true, PhaseRunning, MsgNfsServerRunning
}
