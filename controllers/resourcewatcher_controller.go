/*
Resource watcher will trigger reconcile everytime a watched resource changed in the same namespace the CR is deployed
*/

package controllers

import (
	"context"
	"fmt"
	alpha1 "github.com/stakater/operator-boilerplate/api/v1alpha1"
	"github.com/stakater/operator-boilerplate/utils"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/types"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// ResourceWatcherReconciler reconciles a ResourceWatcher object
type ResourceWatcherReconciler struct {
	utils.ReconcilerBase
}

// +kubebuilder:rbac:groups=watcher.stakater.com.stakater.com,resources=resourcewatchers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=watcher.stakater.com.stakater.com,resources=resourcewatchers/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=watcher.stakater.com.stakater.com,resources=resourcewatchers/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=configmaps,verbs=watch;list
// +kubebuilder:rbac:groups="",resources=secrets,verbs=watch;list
func (r *ResourceWatcherReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	instance := &alpha1.ResourceWatcher{}

	err := r.GetResource(ctx, req.NamespacedName, instance)
	if err != nil {
		return ctrl.Result{}, r.ManageError(ctx, instance, err, "reconciliation failed!")
	}

	if reflect.DeepEqual(&alpha1.ResourceWatcher{}, instance) {
		return ctrl.Result{}, nil
	}

	r.Logger().Info("change detected!, start reconciling ...")
	r.Logger().Info(fmt.Sprintf("found resource with name %s ...", instance.GetName()))

	if instance.Spec.WatchSecrets {
		secrets, err := r.ListResourceMetadata(ctx, v1.SchemeGroupVersion.WithKind("SecretList"), client.InNamespace(req.Namespace))
		if err != nil {
			r.Logger().Error(err, "error fetching secrets")
		}

		instance.Status.Secrets = make(map[string]string)
		for _, secret := range secrets.Items {
			if instance.Status.Secrets[secret.Name] == secret.ResourceVersion {
				continue
			}

			r.Logger().Info(fmt.Sprintf("secret %s changed!", secret.Name))
			instance.Status.Secrets[secret.Name] = secret.ResourceVersion
		}
	}

	return ctrl.Result{}, r.UpdateStatus(ctx, instance)
}

// SetupWithManager sets up the controller with the Manager.
func (r *ResourceWatcherReconciler) SetupWithManager(mgr ctrl.Manager) error {
	// Add all watchers to the cache index with name & namespace for lookup
	IndexField := "resource-watchers"
	if err := mgr.GetFieldIndexer().IndexField(context.Background(), &alpha1.ResourceWatcher{}, IndexField, func(rawObj client.Object) []string {
		configDeployment := rawObj.(*alpha1.ResourceWatcher)
		return []string{configDeployment.Name, configDeployment.Namespace}
	}); err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&alpha1.ResourceWatcher{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})).
		Watches(&source.Kind{Type: &v1.Secret{}}, handler.EnqueueRequestsFromMapFunc(func(object client.Object) []reconcile.Request {
			// For all secret changes call their respective watcher to reconcile
			watchers := &alpha1.ResourceWatcherList{}
			listOps := &client.ListOptions{
				// Look up the watchers we stored above for significantly faster lookup
				FieldSelector: fields.OneTermEqualSelector(IndexField, object.GetNamespace()),
				Namespace:     object.GetNamespace(),
			}
			err := r.GetClient().List(context.TODO(), watchers, listOps)
			if err != nil {
				return []reconcile.Request{}
			}

			requests := make([]reconcile.Request, len(watchers.Items))
			for i, item := range watchers.Items {
				requests[i] = reconcile.Request{
					NamespacedName: types.NamespacedName{
						Name:      item.GetName(),
						Namespace: item.GetNamespace(),
					},
				}
			}
			return requests
		})).
		Complete(r)
}
