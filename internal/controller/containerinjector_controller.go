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
	"fmt"
	applev1 "github.com/Nishant-Ingle/apple-operator/api/v1"
	apps "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// ContainerInjectorReconciler reconciles a ContainerInjector object
type ContainerInjectorReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=apple.nishant.ingle,resources=containerinjectors,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apple.nishant.ingle,resources=containerinjectors/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apple.nishant.ingle,resources=containerinjectors/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ContainerInjector object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.0/pkg/reconcile
func (r *ContainerInjectorReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	injector := &applev1.ContainerInjector{}
	err := r.Client.Get(ctx, req.NamespacedName, injector)

	if err != nil {
		return ctrl.Result{}, err
	}
	deployments := &apps.DeploymentList{}
	err = r.Client.List(ctx, deployments, &client.ListOptions{Namespace: req.NamespacedName.Namespace})

	for _, v := range deployments.Items {
		fmt.Println(v.Name)
		fmt.Println("Labels: ", v.ObjectMeta.Labels)
	}
	// TODO(user): your logic here

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ContainerInjectorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&applev1.ContainerInjector{}).
		Watches(
			&apps.Deployment{},
			handler.EnqueueRequestsFromMapFunc(r.getAll),
		).
		Complete(r)
}

func (r *ContainerInjectorReconciler) getAll(ctx context.Context, obj client.Object) []reconcile.Request {
	result := []reconcile.Request{}
	injectorList := applev1.ContainerInjectorList{}
	r.Client.List(context.Background(), &injectorList)

	for _, injector := range injectorList.Items {
		result = append(result, ctrl.Request{NamespacedName: client.ObjectKey{Namespace: injector.Namespace, Name: injector.Name}})
	}

	return result

	//logger := log.FromContext(ctx)
	//logger.Info("Mapping object to requests", "objectName", obj.GetName())
	//
	//// Example: Create a request to reconcile the same object
	//req := reconcile.Request{NamespacedName: client.ObjectKeyFromObject(obj)}
	//return []reconcile.Request{req}
}