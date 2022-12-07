package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ResourceWatcherSpec struct {
	Foo string `json:"foo,omitempty"`
}

type ResourceWatcherStatus struct {
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

type ResourceWatcher struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ResourceWatcherSpec   `json:"spec,omitempty"`
	Status ResourceWatcherStatus `json:"status,omitempty"`
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
