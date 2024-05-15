/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logger "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
	"time"

	hehev1alpha1 "kubebuilderTest/api/v1alpha1"
)

// TestKindReconciler reconciles a TestKind object
type TestKindReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

//+kubebuilder:rbac:groups=hehe.black.cat,resources=testkinds,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=hehe.black.cat,resources=testkinds/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=hehe.black.cat,resources=testkinds/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the TestKind object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.17.3/pkg/reconcile
func (r *TestKindReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logger.FromContext(ctx)
	testKindObj := &hehev1alpha1.TestKind{}
	if err := r.Client.Get(ctx, req.NamespacedName, testKindObj); err != nil {
		log.Error(err, "unable to fetch TestKind resource")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	log.Info("request name and namespace", "Name", req.NamespacedName.Name, "Namespace", req.NamespacedName.Namespace)
	deployment := &appsv1.Deployment{}

	if err := r.Client.Get(ctx, types.NamespacedName{Namespace: req.Namespace, Name: testKindObj.Spec.DeploymentName}, deployment); errors.IsNotFound(err) {
		if err := r.Client.Create(ctx, getDeployment(testKindObj)); err != nil {
			log.Error(err, "unable to create deployment")
			return ctrl.Result{}, err
		}
	} else if deploymentUpdateRequired(testKindObj, deployment) {
		if err := r.Client.Update(ctx, getDeployment(testKindObj)); err != nil {
			log.Error(err, "unable to update deployment")
			return ctrl.Result{}, err
		}
	} else if err != nil {
		log.Error(err, "error getting deployment")
		return ctrl.Result{}, err
	}

	service := &corev1.Service{}
	if err := r.Client.Get(ctx, types.NamespacedName{Namespace: req.Namespace, Name: testKindObj.Spec.ServiceName}, service); errors.IsNotFound(err) {
		if err := r.Client.Create(ctx, getService(testKindObj)); err != nil {
			log.Error(err, "unable to create service")
			return ctrl.Result{}, err
		}
	} else if serviceUpdateRequired(testKindObj, service) {
		if err := r.Client.Update(ctx, getService(testKindObj)); err != nil {
			log.Error(err, "unable to update service")
			return ctrl.Result{}, err
		}
	}

	if err := r.updateTestKindStatus(ctx, testKindObj, deployment); err != nil {
		log.Error(err, "unable to update TestKindStatus")
		return ctrl.Result{}, err
	}

	r.Recorder.Event(testKindObj, corev1.EventTypeNormal, "Reconciled", "Reconciliation successful")

	return ctrl.Result{}, nil
}

func (r *TestKindReconciler) updateTestKindStatus(ctx context.Context, testKindObj *hehev1alpha1.TestKind, deployment *appsv1.Deployment) error {
	testKindObjCopy := testKindObj.DeepCopy()
	testKindObjCopy.Status.ReplicaCount = deployment.Status.AvailableReplicas

	if err := r.Status().Update(ctx, testKindObjCopy); err != nil {
		return err
	}
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *TestKindReconciler) SetupWithManager(mgr ctrl.Manager) error {
	reconciliationSourceChannel := make(chan event.GenericEvent)

	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				reconciliationSourceChannel <- event.GenericEvent{
					Object: &hehev1alpha1.TestKind{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "bookstore-controller-test-kubebuilder",
							Namespace: "default",
						},
					},
				}
			}
		}
	}()

	return ctrl.NewControllerManagedBy(mgr).
		For(&hehev1alpha1.TestKind{}).
		//Owns(&corev1.Pod{}).
		///Owns(&appsv1.Deployment{}).
		//Owns(&corev1.Service{}).
		Watches(&source.Channel{Source: reconciliationSourceChannel}, &handler.EnqueueRequestForObject{}).
		Complete(r)
}
