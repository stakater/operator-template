package utils

import (
	"context"
	"fmt"
	"github.com/go-logr/logr"
	"github.com/stakater/operator-boilerplate/utils/api"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type ReconcilerFunc func(r *ReconcilerBase) (reconcile.Result, error)

type ReconcilerBase struct {
	Name     string
	client   client.Client
	scheme   *runtime.Scheme
	recorder record.EventRecorder
}

// Create reconciler
func NewBaseReconciler(client client.Client, scheme *runtime.Scheme, recorder record.EventRecorder) ReconcilerBase {
	return ReconcilerBase{
		client:   client,
		scheme:   scheme,
		recorder: recorder,
	}
}

// Create reconciler from manager instance
func NewFromManager(mgr manager.Manager, ctrName string) ReconcilerBase {
	return NewBaseReconciler(mgr.GetClient(), mgr.GetScheme(), mgr.GetEventRecorderFor(ctrName))
}

func (r *ReconcilerBase) GetClient() client.Client {
	return r.client
}

func (r *ReconcilerBase) GetRecorder() record.EventRecorder {
	return r.recorder
}

func (r *ReconcilerBase) GetScheme() *runtime.Scheme {
	return r.scheme
}

func (r *ReconcilerBase) Logger() logr.Logger {
	return logf.Log.WithName(r.Name)
}

// Update Status conditions to inform users of changes
func (r *ReconcilerBase) UpdateCondition(cr client.Object, condition metav1.Condition) {
	if conditionsAware, updateStatus := cr.(api.IConditionsAware); updateStatus {
		conditionsAware.SetConditions(api.AddOrReplaceCondition(condition, conditionsAware.GetConditions()))
	}
}

// Update failed conditions which will also send out error event
func (r *ReconcilerBase) UpdateFailedCondition(cr client.Object, conditionType string, msg string, issue error) {
	r.GetRecorder().Event(cr, "Warning", "ProcessingError", issue.Error())
	r.UpdateCondition(cr, metav1.Condition{
		Type:               conditionType,
		LastTransitionTime: metav1.Now(),
		ObservedGeneration: cr.GetGeneration(),
		Message:            msg,
		Reason:             api.ReconcileErrorReason,
		Status:             metav1.ConditionFalse,
	})
}

// Update success conditions only
func (r *ReconcilerBase) UpdateSuccessCondition(cr client.Object, conditionType string, msg string) {
	r.UpdateCondition(cr, metav1.Condition{
		Type:               conditionType,
		LastTransitionTime: metav1.Now(),
		ObservedGeneration: cr.GetGeneration(),
		Reason:             api.ReconcileSuccessReason,
		Status:             metav1.ConditionTrue,
		Message:            msg,
	})
}

// Create resource with owner reference to make sure when owner seize to exists all secondary resources are garbage collected
func (r *ReconcilerBase) CreateOwnedResource(context context.Context, owner client.Object, obj client.Object) error {
	_ = controllerutil.SetOwnerReference(owner, obj, r.GetScheme())
	obj.SetNamespace(owner.GetNamespace())

	err := r.GetClient().Create(context, obj)

	conditionType := fmt.Sprintf("Create-%s", obj.GetName())
	if err != nil && !errors.IsAlreadyExists(err) {
		r.UpdateFailedCondition(owner, conditionType, fmt.Sprintf("failed to create %s", obj.GetName()), err)
		return err
	}

	r.UpdateSuccessCondition(owner, conditionType, fmt.Sprintf("%s created", obj.GetName()))
	return r.GetClient().Status().Update(context, owner)
}

// Create a simple resource
func (r *ReconcilerBase) CreateResource(context context.Context, obj client.Object) error {
	err := r.GetClient().Create(context, obj)
	if err != nil && !errors.IsAlreadyExists(err) {
		return err
	}

	return nil
}

// Update resource status
func (r *ReconcilerBase) UpdateStatus(context context.Context, cr client.Object) error {
	err := r.GetClient().Status().Update(context, cr)
	if err != nil {
		r.Logger().Error(err, "unable to update status")
		return err
	}
	return nil
}

/*
Error handling will do 2 things
1. Add failed condition
2. Send error event
*/
func (r *ReconcilerBase) ManageError(context context.Context, cr client.Object, issue error, msg string) error {
	r.GetRecorder().Event(cr, "Warning", "ProcessingError", issue.Error())
	if conditionsAware, updateStatus := any(cr).(api.IConditionsAware); updateStatus {
		condition := metav1.Condition{
			Type:               api.ReconcileError,
			LastTransitionTime: metav1.Now(),
			ObservedGeneration: cr.GetGeneration(),
			Message:            msg,
			Reason:             api.ReconcileErrorReason,
			Status:             metav1.ConditionTrue,
		}
		conditionsAware.SetConditions(api.AddOrReplaceCondition(condition, conditionsAware.GetConditions()))
		return r.UpdateStatus(context, cr)
	} else {
		r.Logger().Info("object is not IConditionsAware, not setting status")
	}
	return issue
}

/*
Error handling will add success condition to status
*/
func (r *ReconcilerBase) ManageSuccess(context context.Context, cr client.Object, msg string) error {
	if conditionsAware, updateStatus := cr.(api.IConditionsAware); updateStatus {
		condition := metav1.Condition{
			Type:               api.ReconcileSuccess,
			LastTransitionTime: metav1.Now(),
			ObservedGeneration: cr.GetGeneration(),
			Reason:             api.ReconcileSuccessReason,
			Status:             metav1.ConditionTrue,
			Message:            msg,
		}
		conditionsAware.SetConditions(api.AddOrReplaceCondition(condition, conditionsAware.GetConditions()))
		return r.UpdateStatus(context, cr)
	} else {
		r.Logger().Info("object is not IConditionsAware, not setting status")
	}
	return nil
}

// Fetch the controller resource
func (r *ReconcilerBase) GetCtrResource(ctx context.Context, ns types.NamespacedName, instance client.Object) error {
	err := r.GetClient().Get(ctx, ns, instance)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return nil
}
