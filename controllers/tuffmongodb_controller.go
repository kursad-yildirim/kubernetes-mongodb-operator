package controllers

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	ctrllog "sigs.k8s.io/controller-runtime/pkg/log"

	appv1alpha1 "github.com/kursad-yildirim/kubernetes-mongodb-operator/api/v1alpha1"
	"github.com/kursad-yildirim/kubernetes-mongodb-operator/manifests/deployment"

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

func (r *TuffMongoDBReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := ctrllog.FromContext(ctx)

	instance := &appv1alpha1.TuffMongoDB{}
	err := r.Get(context.TODO(), req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}
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
		mongoPod := newMongoPodForCR(tuffMongoDB)
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

func (r *TuffMongoDBReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appv1alpha1.TuffMongoDB{}).
		Complete(r)
}
