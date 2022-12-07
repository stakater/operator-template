/*
Resource watcher will trigger reconcile everytime a watched resource changed
*/

package controllers

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	watcherstakatercomv1alpha1 "github.com/stakater/operator-boilerplate/api/v1alpha1"
)

// ResourceWatcherReconciler reconciles a ResourceWatcher object
type ResourceWatcherReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=watcher.stakater.com.stakater.com,resources=resourcewatchers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=watcher.stakater.com.stakater.com,resources=resourcewatchers/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=watcher.stakater.com.stakater.com,resources=resourcewatchers/finalizers,verbs=update
func (r *ResourceWatcherReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ResourceWatcherReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&watcherstakatercomv1alpha1.ResourceWatcher{}).
		Complete(r)
}
