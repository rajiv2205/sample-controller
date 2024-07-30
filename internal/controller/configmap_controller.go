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

/*
import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	apiv1alpha1 "github.com/sample-controller/api/v1alpha1"
)
*/
import (
	"context"
	"fmt"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// ConfigMapReconciler reconciles a ConfigMap object
type ConfigMapReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	lastSeen map[client.ObjectKey]corev1.ConfigMap
}

// +kubebuilder:rbac:groups=api.operatortest.io,resources=configmaps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=api.operatortest.io,resources=configmaps/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=api.operatortest.io,resources=configmaps/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ConfigMap object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.18.4/pkg/reconcile
func (r *ConfigMapReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	_ = log.FromContext(ctx)

	namespace := "default"
	configMapName := "webapp-config"
	deploymentName := "webapp"

	key := types.NamespacedName{
		Namespace: namespace,
		Name:      configMapName,
	}

	cmData := &corev1.ConfigMap{}
	err := r.Get(ctx, key, cmData)
	if err != nil {
		// Handle the error (e.g., ConfigMap not found)
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	lastSeenConfigMap := r.lastSeen[key]

	if lastSeenConfigMap.ResourceVersion == cmData.ResourceVersion {
		// No update, just return
		return ctrl.Result{}, nil

	}

	if len(lastSeenConfigMap.ResourceVersion) == 0 {
		// add the cmData, just return
		r.lastSeen[key] = *cmData
		return ctrl.Result{}, nil

	}

	r.lastSeen[key] = *cmData

	log.Log.Info(fmt.Sprintf("ConfigMap Updated: %s\n", cmData.Name))
	log.Log.Info(fmt.Sprintf("ConfigMap Data: %v\n", cmData.Data))

	err = r.triggerRolloutRestart(ctx, namespace, deploymentName)
	if err != nil {
		return ctrl.Result{}, err
	}
	// Log the rolling restart
	log.Log.Info(fmt.Sprintf("Triggered rolling restart for Deployment: %s\n", deploymentName))

	return ctrl.Result{}, nil
}

func (r *ConfigMapReconciler) triggerRolloutRestart(ctx context.Context, namespace string, deploymentNameToRestart string) error {

	deploymentName := deploymentNameToRestart
	deploy := &appsv1.Deployment{}
	err := r.Get(ctx, types.NamespacedName{
		Namespace: namespace,
		Name:      deploymentName,
	}, deploy)

	if err != nil {
		// Handle the error (e.g., deployment not found)
		return err
	}

	// Add or update the annotation to trigger the restart
	deploy.Spec.Template.Annotations["webapp/restart"] = fmt.Sprintf("%d", time.Now().Unix())

	// Update the Deployment with the new annotation
	return r.Update(ctx, deploy)

}

// +kubebuilder:rbac:groups=api.operatortest.io,resources=configmaps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=api.operatortest.io,resources=configmaps/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=api.operatortest.io,resources=configmaps/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ConfigMap object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.18.4/pkg/reconcile

// SetupWithManager sets up the controller with the Manager.
func (r *ConfigMapReconciler) SetupWithManager(mgr ctrl.Manager) error {
	r.lastSeen = make(map[client.ObjectKey]corev1.ConfigMap)
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1.ConfigMap{}).
		Complete(r)
}
