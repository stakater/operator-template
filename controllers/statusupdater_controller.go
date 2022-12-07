/*
Status updater will update status with condition and events
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

type StatusUpdaterReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=watcher.stakater.com.stakater.com,resources=statusupdaters,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=watcher.stakater.com.stakater.com,resources=statusupdaters/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=watcher.stakater.com.stakater.com,resources=statusupdaters/finalizers,verbs=update

func (r *StatusUpdaterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *StatusUpdaterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&watcherstakatercomv1alpha1.StatusUpdater{}).
		Complete(r)
}
