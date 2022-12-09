package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ResourceWatcherSpec struct {
	// +optional
	WatchSecrets bool `json:"watchSecrets,omitempty"`

	// +optional
	WatchConfigmaps bool `json:"watchConfigmaps,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

type ResourceWatcher struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ResourceWatcherSpec   `json:"spec,omitempty"`
	Status ResourceWatcherStatus `json:"status,omitempty"`
}

type ResourceWatcherStatus struct {
	// +patchMergeKey=type
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=type
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`

	Secrets map[string]string `json:"secrets,omitempty"`
}

func (m *ResourceWatcher) GetConditions() []metav1.Condition {
	return m.Status.Conditions
}

func (m *ResourceWatcher) SetConditions(conditions []metav1.Condition) {
	m.Status.Conditions = conditions
}

// +kubebuilder:object:root=true
type ResourceWatcherList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ResourceWatcher `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ResourceWatcher{}, &ResourceWatcherList{})
}
