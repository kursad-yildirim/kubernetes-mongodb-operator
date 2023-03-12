package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// TuffMongoDBSpec defines the desired state of TuffMongoDB
type TuffMongoDBSpec struct {
	Replicas     int32                  `json:"replicas,omitempty"`
	Image        string                 `json:"image,omitempty"`
	Name         string                 `json:"name,omitEmpty"`
	Namespace    string                 `json:"namespace,omitEmpty"`
	Component    string                 `json:"component,omitEmpty"`
	Ports        []corev1.ContainerPort `json:"ports,omitEmpty"`
	VolumeMounts []corev1.VolumeMount   `json:"volumeMounts,omitEmpty"`
	Volumes      []corev1.Volume        `json:"volumes,omitEmpty"`
}

// TuffMongoDBStatus defines the observed state of TuffMongoDB
type TuffMongoDBStatus struct {
	MongoPodNames          []string `json:"mongoPodNames"`
	MongoAvailableReplicas int32    `json:"mongoAvailableReplicas"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// TuffMongoDB is the Schema for the tuffmongodbs API
type TuffMongoDB struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TuffMongoDBSpec   `json:"spec,omitempty"`
	Status TuffMongoDBStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// TuffMongoDBList contains a list of TuffMongoDB
type TuffMongoDBList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TuffMongoDB `json:"items"`
}

func init() {
	SchemeBuilder.Register(&TuffMongoDB{}, &TuffMongoDBList{})
}
