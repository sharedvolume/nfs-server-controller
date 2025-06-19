/*
Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type StorageSpec struct {
	Capacity         string `json:"capacity"`
	StorageClassName string `json:"storageClassName,omitempty"`
	PersistentVolume string `json:"persistentVolume,omitempty"`
}

type NfsServerSpec struct {
	Storage  StorageSpec `json:"storage"`
	Replicas *int32      `json:"replicas,omitempty"`
	Path     string      `json:"path,omitempty"`
	Image    string      `json:"image,omitempty"`
	Address  string      `json:"address,omitempty"`
}

type NfsServerStatus struct {
	Ready   bool   `json:"ready"`
	Phase   string `json:"phase,omitempty"`
	Message string `json:"message,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=nfs,scope=Namespaced
// +kubebuilder:printcolumn:name="Ready",type=boolean,JSONPath=".status.ready",description="NFS Server Ready"
// +kubebuilder:printcolumn:name="Address",type=string,JSONPath=".spec.address",description="NFS Server Address"
// +kubebuilder:printcolumn:name="Capacity",type=string,JSONPath=".spec.storage.capacity",description="NFS Storage Capacity"

// NfsServer is the Schema for the nfsservers API.
type NfsServer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NfsServerSpec   `json:"spec,omitempty"`
	Status NfsServerStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// NfsServerList contains a list of NfsServer.
type NfsServerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NfsServer `json:"items"`
}

func init() {
	SchemeBuilder.Register(&NfsServer{}, &NfsServerList{})
}
