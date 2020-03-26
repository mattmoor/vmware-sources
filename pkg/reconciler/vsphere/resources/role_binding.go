/*
Copyright 2020 The Knative Authors

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

package resources

import (
	"context"

	"github.com/mattmoor/vmware-sources/pkg/apis/sources/v1alpha1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/pkg/kmeta"
)

// MakeRoleBinding creates a RoleBinding object for the receive adapter
// service account 'sa' in the Namespace 'ns'. This is necessary for
// the receive adapter to be able to store state in configmaps.
// TOOD: if you create serviceaccounts instead of using default, pass it in here and use that instead.
//func MakeRoleBinding(ctx context.Context, vms *v1alpha1.VSphereSource, name string, nsName string, sa *corev1.ServiceAccount, clusterRoleName string) *rbacv1.RoleBinding {
func MakeRoleBinding(ctx context.Context, vms *v1alpha1.VSphereSource, name string, nsName string) *rbacv1.RoleBinding {
	return &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			OwnerReferences: []metav1.OwnerReference{*kmeta.NewControllerRef(vms)},
			Name:            name,
			Namespace:       nsName,
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     "receive-adapter-cm-reader",
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Namespace: nsName,
				Name:      "default",
			},
		},
	}
}
