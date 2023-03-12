package controllers

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	//	ctrllog "sigs.k8s.io/controller-runtime/pkg/log"

	appv1alpha1 "github.com/kursad-yildirim/kubernetes-mongodb-operator/api/v1alpha1"
	"github.com/kursad-yildirim/kubernetes-mongodb-operator/manifests/deployment"
	//	"github.com/kursad-yildirim/kubernetes-mongodb-operator/manifests/pod"

	//	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	//	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//	"k8s.io/apimachinery/pkg/labels"
	//
	// "reflect"
	// "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
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

	instance := &appv1alpha1.TuffMongoDB{}
	err := r.Get(context.TODO(), req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}
	createMongoDeployment(*instance)
	return ctrl.Result{}, nil
}

func createMongoDeployment(tmCRD appv1alpha1.TuffMongoDB) {

	tuffLabels := map[string]string{
		"app-name":  tmCRD.Spec.Name,
		"component": tmCRD.Spec.Component,
	}

	// template := newPodTemplateSpec(context.TODO(), tmCRD)

	dpl := deployment.New(tmCRD.Spec.Name, tmCRD.Spec.Namespace, tuffLabels, tmCRD.Spec.Replicas).
		WithPaused(false).
		Build()
	fmt.Print(dpl)
}

/*
func newPodTemplateSpec(ctx context.Context, tmCRD appv1alpha1.TuffMongoDB) corev1.PodTemplateSpec {

		containers := []corev1.Container{
			newMongoDBContainer(tmCRD),
		}

		podSpec := pod.NewSpec(tmCRD.Spec.Name, containers, tmCRD.Spec.Volumes).
			Build()

		return corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: *labels,
			},
			Spec: *podSpec,
		}
	}
*/
func newMongoDBContainer(tmCRD appv1alpha1.TuffMongoDB) corev1.Container {
	return corev1.Container{
		Name:            tmCRD.Spec.Name,
		Image:           tmCRD.Spec.Image,
		ImagePullPolicy: "IfNotPresent",
		Ports:           tmCRD.Spec.Ports,
		VolumeMounts:    tmCRD.Spec.VolumeMounts,
	}
}

func (r *TuffMongoDBReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appv1alpha1.TuffMongoDB{}).
		Complete(r)
}
