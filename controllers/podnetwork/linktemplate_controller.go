/*
Copyright 2021.

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

package podnetwork

import (
	"context"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	podnetworkv1alpha1 "github.com/opdev/pod-network-operator/apis/podnetwork/v1alpha1"
)

// LinkTemplateReconciler reconciles a LinkTemplate object
type LinkTemplateReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
	linkTemplateList *podnetworkv1alpha1.LinkTemplateList
}

// +kubebuilder:rbac:groups=podnetwork.opdev.io,resources=linktemplates,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=podnetwork.opdev.io,resources=linktemplates/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=podnetwork.opdev.io,resources=linktemplates/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the LinkTemplate object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.0/pkg/reconcile
func (r *LinkTemplateReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	
	reqLogger = r.Log.WithValues("linktemplate", req.NamespacedName)
	
	r.linkTemplateList = &podnetworkv1alpha1.LinkTemplateList{}
	err := r.Client.List(context.TODO(), r.linkTemplateList)
	if err != nil {
		return ctrl.Result{}, err
	}

	if len(r.linkTemplateList.Items) <= 0 {
		return ctrl.Result{}, nil
	}

	// TODO : UPDATE CONDITIONS/PHASE ON STATUS HERE !!!!!



	// FOR EACH linkTemplate created identify the pods in the namespace
	// that will have labels/annotations asking for that particular template
	// apply the template to each one of them.

	
	for _, linkTemplate := range r.linkTemplateList.Items {










	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *LinkTemplateReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&podnetworkv1alpha1.LinkTemplate{}).
		Complete(r)
}
