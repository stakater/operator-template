/*
Status updater will update status with condition and events
*/

package controllers

import (
	"context"
	"fmt"
	"github.com/stakater/operator-boilerplate/api/v1alpha1"
	"github.com/stakater/operator-boilerplate/utils"
	ctrl "sigs.k8s.io/controller-runtime"
)

type StatusUpdaterReconciler struct {
	utils.ReconcilerBase
}

//+kubebuilder:rbac:groups=watcher.stakater.com.stakater.com,resources=statusupdaters,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=watcher.stakater.com.stakater.com,resources=statusupdaters/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=watcher.stakater.com.stakater.com,resources=statusupdaters/finalizers,verbs=update

func (r *StatusUpdaterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	instance := &v1alpha1.StatusUpdater{}
	err := r.GetCtrResource(ctx, req.NamespacedName, instance)
	if err != nil {
		r.Logger().Info("no StatusUpdater resource found ")
		return ctrl.Result{}, err
	}

	r.Logger().Info(fmt.Sprintf("found StatusUpdater [%s], start reconciling ...", instance.GetName()))

	// Make sure all incidents are reported
	for _, incident := range instance.Spec.Incidents {
		found := false
		for _, condition := range instance.GetConditions() {
			if incident.Type == condition.Type {
				found = true
				break
			}
		}

		// If incident has not been reported then report it
		if !found {
			if incident.Reason == v1alpha1.FailedReason {
				r.UpdateFailedCondition(instance, incident.Type, incident.Message, fmt.Errorf(incident.Message))
			} else {
				r.UpdateSuccessCondition(instance, incident.Type, incident.Message)
			}
		}
	}

	// Update conditions will not trigger status update automatically so we do it here
	err = r.UpdateStatus(ctx, instance)
	if err != nil {
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *StatusUpdaterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.StatusUpdater{}).
		Complete(r)
}
