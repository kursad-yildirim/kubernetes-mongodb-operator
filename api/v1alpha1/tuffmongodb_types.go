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

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// TuffMongoDBSpec defines the desired state of TuffMongoDB
type TuffMongoDBSpec struct {
	MongoReplicas      int32                  `json:"mongoReplicas,omitempty"`
	MongoImage         string                 `json:"mongoImage,omitempty"`
	MongoContainerName string                 `json:"mongoContainerName,omitEmpty"`
	MongoPorts         []corev1.ContainerPort `json:"mongoPorts,omitEmpty"`
	MongoVolumeMounts  []corev1.VolumeMount   `json:"mongoVolumeMounts,omitEmpty"`
	MongoVolumes       []corev1.Volume        `json:"mongoVolumes,omitEmpty"`
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
