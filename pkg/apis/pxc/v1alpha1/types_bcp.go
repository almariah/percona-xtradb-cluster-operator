package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type PerconaXtraDBBackupList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []PerconaXtraDBBackup `json:"items"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type PerconaXtraDBBackup struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              PXCBackupSpec   `json:"spec"`
	Status            PXCBackupStatus `json:"status,omitempty"`
}

type PXCBackupSpec struct {
	PXCCluster string `json:"pxcCluster"`
	Storage    string `json:"storage,omitempty"`
}

type PXCBackupStatus struct {
}

// OwnerRef returns OwnerReference to object
func (cr *PerconaXtraDBBackup) OwnerRef() metav1.OwnerReference {
	trueVar := true

	return metav1.OwnerReference{
		APIVersion: SchemeGroupVersion.String(),
		Kind:       cr.Kind,
		Name:       cr.Name,
		UID:        cr.UID,
		Controller: &trueVar,
	}
}
