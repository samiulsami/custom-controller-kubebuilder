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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// TestKindSpec defines the desired state of TestKind
type TestKindSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	EnvAdminUsername string `json:"envAdminUsername"`
	EnvAdminPassword string `json:"envAdminPassword"`
	EnvJWTSECRET     string `json:"envJWTSECRET"`

	DeploymentImageName string `json:"deploymentImageName"`
	DeploymentImageTag  string `json:"deploymentImageTag"`
	// LoadBalancer, ClusterIP, NodePort
	ImagePullPolicy string `json:"imagePullPolicy"`

	DeploymentName string `json:"deploymentName"`
	Replicas       *int32 `json:"replicas"`
	ServiceName    string `json:"serviceName"`
	ServiceType    string `json:"serviceType"`

	//3000
	ContainerPort int32 `json:"containerPort"`
	//3000
	NodePort int32 `json:"nodePort"`
	//3000
	TargetPort int32 `json:"targetPort"`

	TestMap2 map[string]string `json:"testMap2,omitempty"`
	// Huge nested map test
	UnholyAbomination map[string]map[string]map[string]map[string]map[string]map[string]map[string]map[string]map[string]map[string]map[string]map[string]map[string]map[string]map[string]map[string]map[string]map[string]map[string]map[string]map[string]map[string]map[string]map[string]map[string]map[string]map[string]map[string]map[string]map[string]map[string]map[string]map[string]map[string]map[string]map[string]map[string]map[string]map[string]map[string]map[string]map[string]int32 `json:"unholyAbomination,omitempty"`
}

// TestKindStatus defines the observed state of TestKind
type TestKindStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	ReplicaCount int32 `json:"replicaCount"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// TestKind is the Schema for the testkinds API
type TestKind struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TestKindSpec   `json:"spec,omitempty"`
	Status TestKindStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// TestKindList contains a list of TestKind
type TestKindList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TestKind `json:"items"`
}

func (testKind *TestKind) GetSelectorLabels() map[string]string {
	return map[string]string{
		"app":        testKind.Name + "-app",
		"controller": testKind.Name + "-customController1",
	}
}

func init() {
	SchemeBuilder.Register(&TestKind{}, &TestKindList{})
}
