/*
Copyright 2023.

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

package controllers

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	ctrllog "sigs.k8s.io/controller-runtime/pkg/log"

	appv1alpha1 "github.com/kursad-yildirim/kubernetes-mongodb-operator/api/v1alpha1"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// TuffMongoDBReconciler reconciles a TuffMongoDB object
type TuffMongoDBReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=app.1k.local,resources=tuffmongodbs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=app.1k.local,resources=tuffmongodbs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=app.1k.local,resources=tuffmongodbs/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the TuffMongoDB object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *TuffMongoDBReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := ctrllog.FromContext(ctx)

	// Fetch the Mongodb.db instance
	instance := &appv1alpha1.TuffMongoDB{}
	err := r.Get(context.TODO(), req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return ctrl.Result{}, err
	}
	// List all pods owned by this MicroMongo instance
	tuffMongoDB := instance
	mongoPodList := &corev1.PodList{}
	tmLabels := map[string]string{
		"app":     tuffMongoDB.Name,
		"version": "v0.1",
	}
	tmLabelSelector := labels.SelectorFromSet(tmLabels)
	tmListOptions := &client.ListOptions{Namespace: tuffMongoDB.Namespace, LabelSelector: tmLabelSelector}
	if err = r.List(context.TODO(), mongoPodList, tmListOptions); err != nil {
		return ctrl.Result{}, err
	}
	// Count the pods that are pending or running as available
	var availableMongoPods []corev1.Pod
	for _, pod := range mongoPodList.Items {
		if pod.ObjectMeta.DeletionTimestamp != nil {
			continue
		}
		if pod.Status.Phase == corev1.PodRunning || pod.Status.Phase == corev1.PodPending {
			availableMongoPods = append(availableMongoPods, pod)
		}
	}
	numAvailableMongoPods := int32(len(availableMongoPods))
	availableMongoPodNames := []string{}
	for _, pod := range availableMongoPods {
		availableMongoPodNames = append(availableMongoPodNames, pod.ObjectMeta.Name)
	}
	// Update the status if necessary
	status := appv1alpha1.TuffMongoDBStatus{
		MongoPodNames:          availableMongoPodNames,
		MongoAvailableReplicas: numAvailableMongoPods,
	}
	if !reflect.DeepEqual(tuffMongoDB.Status, status) {
		tuffMongoDB.Status = status
		err = r.Status().Update(context.TODO(), tuffMongoDB)
		if err != nil {
			log.Error(err, "Failed  to update tuffMongoDB status.")
			return ctrl.Result{}, err
		}
	}

	if numAvailableMongoPods > tuffMongoDB.Spec.MongoReplicas {
		log.Info("Scaling down mongodb pods", "Currently available", numAvailableMongoPods, "Required replicas", tuffMongoDB.Spec.MongoReplicas)
		diff := numAvailableMongoPods - tuffMongoDB.Spec.MongoReplicas
		dpods := availableMongoPods[:diff]
		for _, podToDelete := range dpods {
			err = r.Delete(context.TODO(), &podToDelete)
			if err != nil {
				log.Error(err, "Failed to delete mongodb pod", "mongoPod.name", podToDelete.Name)
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{Requeue: true}, nil
	}

	if numAvailableMongoPods < tuffMongoDB.Spec.MongoReplicas {
		log.Info("Scaling up mongodb pods", "Currently available", numAvailableMongoPods, "Required replicas", tuffMongoDB.Spec.MongoReplicas)
		// Define a new mongo pod object
		mongoPod := newMongoPodForCR(tuffMongoDB)
		// Set TuffMongoDB instance as the owner and controller
		if err := controllerutil.SetControllerReference(tuffMongoDB, mongoPod, r.Scheme); err != nil {
			return ctrl.Result{}, err
		}
		err = r.Create(context.TODO(), mongoPod)
		if err != nil {
			log.Error(err, "Failed to create mongodb pod", "MongoPod.name", mongoPod.Name)
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	}

	return ctrl.Result{}, nil
}

func newMongoPodForCR(cr *appv1alpha1.TuffMongoDB) (mongoPod *corev1.Pod) {
	tmLabels := map[string]string{
		"app":     cr.Name,
		"version": "v0.1",
	}
	mongoPod = &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: cr.Name + "-pod-",
			Namespace:    cr.Namespace,
			Labels:       tmLabels,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:         cr.Spec.MongoContainerName,
					Image:        cr.Spec.MongoImage,
					Ports:        cr.Spec.MongoPorts,
					VolumeMounts: cr.Spec.MongoVolumeMounts,
					Command:      []string{"sleep", "16000"},
				},
			},
			Volumes: cr.Spec.MongoVolumes,
		},
	}

	return
}

// SetupWithManager sets up the controller with the Manager.
func (r *TuffMongoDBReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appv1alpha1.TuffMongoDB{}).
		Complete(r)
}
