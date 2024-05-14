package controller

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"kubebuilderTest/api/v1alpha1"
	hehev1alpha1 "kubebuilderTest/api/v1alpha1"
)

func getDeployment(testKindObj *v1alpha1.TestKind) *appsv1.Deployment {
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testKindObj.Spec.DeploymentName,
			Namespace: testKindObj.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(testKindObj, v1alpha1.GroupVersion.WithKind("TestKind")),
			},
		},

		Spec: appsv1.DeploymentSpec{
			Replicas: testKindObj.Spec.Replicas,

			Selector: &metav1.LabelSelector{
				MatchLabels: testKindObj.GetSelectorLabels(),
			},

			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: testKindObj.GetSelectorLabels(),
				},

				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            testKindObj.Spec.DeploymentName,
							Image:           testKindObj.Spec.DeploymentImageName + ":" + testKindObj.Spec.DeploymentImageTag,
							ImagePullPolicy: corev1.PullPolicy(testKindObj.Spec.ImagePullPolicy),

							Ports: []corev1.ContainerPort{
								{
									ContainerPort: testKindObj.Spec.ContainerPort,
								},
							},

							Env: []corev1.EnvVar{
								{
									Name: "AdminUsername",
									//Value: "admin22",
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: &corev1.SecretKeySelector{
											LocalObjectReference: corev1.LocalObjectReference{
												Name: "env-secrets",
											},
											Key: testKindObj.Spec.EnvAdminUsername,
											Optional: func() *bool {
												var flag = false
												return &flag
											}(),
										},
									},
								},

								{
									Name: "AdminPassword",
									//Value: "admin72",
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: &corev1.SecretKeySelector{
											LocalObjectReference: corev1.LocalObjectReference{
												Name: "env-secrets",
											},
											Key: testKindObj.Spec.EnvAdminPassword,
											Optional: func() *bool {
												var flag = false
												return &flag
											}(),
										},
									},
								},

								{
									Name: "JWTSECRET",
									//Value: "orangeCat",
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: &corev1.SecretKeySelector{
											LocalObjectReference: corev1.LocalObjectReference{
												Name: "env-secrets",
											},
											Key: testKindObj.Spec.EnvJWTSECRET,
											Optional: func() *bool {
												var flag = false
												return &flag
											}(),
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func getService(testKindObj *v1alpha1.TestKind) *corev1.Service {
	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "core/v1",
		},

		ObjectMeta: metav1.ObjectMeta{
			Name:      testKindObj.Spec.ServiceName,
			Namespace: testKindObj.Namespace,
			Labels:    testKindObj.GetSelectorLabels(),
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(testKindObj, v1alpha1.GroupVersion.WithKind("TestKind")),
			},
		},

		Spec: corev1.ServiceSpec{
			Selector: testKindObj.GetSelectorLabels(),
			Type:     corev1.ServiceType(testKindObj.Spec.ServiceType),
			Ports: []corev1.ServicePort{
				{
					Protocol:   corev1.ProtocolTCP,
					Port:       testKindObj.Spec.ContainerPort,
					TargetPort: intstr.FromInt32(testKindObj.Spec.TargetPort),
					NodePort:   testKindObj.Spec.NodePort,
				},
			},
		},
	}
}

func deploymentUpdateRequired(testKindObj *hehev1alpha1.TestKind, deployment *appsv1.Deployment) bool {
	if testKindObj.Spec.DeploymentName != deployment.Name {
		return true
	}
	if testKindObj.Spec.Replicas != nil && testKindObj.Spec.Replicas != deployment.Spec.Replicas {
		return true
	}
	return false
}

func serviceUpdateRequired(testKindObj *hehev1alpha1.TestKind, service *corev1.Service) bool {
	if testKindObj.Spec.ServiceName != service.Name || testKindObj.Spec.ServiceType != string(service.Spec.Type) {
		return true
	}
	return false
}
