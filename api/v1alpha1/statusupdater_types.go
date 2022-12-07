package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Reason string

var FailedReason Reason = "Failed"
var SuccessReason Reason = "Success"

type Incident struct {
	// +required
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`^([a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*/)?(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])$`
	// +kubebuilder:validation:MaxLength=316
	Type string `json:"type" protobuf:"bytes,1,opt,name=type"`
	// +required
	// +kubebuilder:validation:Enum:=Failed;Success;
	Reason Reason `json:"reason" protobuf:"bytes,5,opt,name=reason"`
	// message is a human readable message indicating details about the transition.
	// This may be an empty string.
	// +required
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MaxLength=32768
	Message string `json:"message" protobuf:"bytes,6,opt,name=message"`
}
type StatusUpdaterSpec struct {
	// +required
	Incidents []Incident `json:"incidents,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

type StatusUpdater struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   StatusUpdaterSpec   `json:"spec,omitempty"`
	Status StatusUpdaterStatus `json:"status,omitempty"`
}

type StatusUpdaterStatus struct {
	// +patchMergeKey=type
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=type
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`
}

func (m *StatusUpdater) GetConditions() []metav1.Condition {
	return m.Status.Conditions
}

func (m *StatusUpdater) SetConditions(conditions []metav1.Condition) {
	m.Status.Conditions = conditions
}

//+kubebuilder:object:root=true

type StatusUpdaterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []StatusUpdater `json:"items"`
}

func init() {
	SchemeBuilder.Register(&StatusUpdater{}, &StatusUpdaterList{})
}
