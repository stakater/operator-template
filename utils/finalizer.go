package utils

import (
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// IsBeingDeleted returns whether this object has been requested to be deleted
func IsBeingDeleted(obj client.Object) bool {
	return !obj.GetDeletionTimestamp().IsZero()
}

// HasFinalizer returns whether this object has the passed finalizer
// Deprecated use controllerutil.ContainsFinalizer
func HasFinalizer(obj client.Object, finalizer string) bool {
	return controllerutil.ContainsFinalizer(obj, finalizer)
}

// AddFinalizer adds the passed finalizer this object
// Deprecated use controllerutil.AddFinalizer
func AddFinalizer(obj client.Object, finalizer string) {
	controllerutil.AddFinalizer(obj, finalizer)
}

// RemoveFinalizer removes the passed finalizer from object
// Deprecated use controllerutil.RemoveFinalizer
func RemoveFinalizer(obj client.Object, finalizer string) {
	controllerutil.RemoveFinalizer(obj, finalizer)
}
