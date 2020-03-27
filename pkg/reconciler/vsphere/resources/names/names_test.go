/*
Copyright 2018 The Knative Authors

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

package names

import (
	"strings"
	"testing"

	"github.com/mattmoor/vmware-sources/pkg/apis/sources/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestNames(t *testing.T) {
	tests := []struct {
		name string
		vss  *v1alpha1.VSphereSource
		f    func(*v1alpha1.VSphereSource) string
		want string
	}{{
		name: "Deployment too long",
		vss: &v1alpha1.VSphereSource{
			ObjectMeta: metav1.ObjectMeta{
				Name: strings.Repeat("f", 63),
			},
		},
		f:    Deployment,
		want: "ffffffffffffffffffff105d7597f637e83cc711605ac3ea4957-deployment",
	}, {
		name: "Deployment long enough",
		vss: &v1alpha1.VSphereSource{
			ObjectMeta: metav1.ObjectMeta{
				Name: strings.Repeat("f", 52),
			},
		},
		f:    Deployment,
		want: strings.Repeat("f", 52) + "-deployment",
	}, {
		name: "Deployment",
		vss: &v1alpha1.VSphereSource{
			ObjectMeta: metav1.ObjectMeta{
				Name: "foo",
			},
		},
		f:    Deployment,
		want: "foo-deployment",
	}, {
		name: "SinkBinding, barely fits",
		vss: &v1alpha1.VSphereSource{
			ObjectMeta: metav1.ObjectMeta{
				Name: strings.Repeat("u", 50),
			},
		},
		f:    SinkBinding,
		want: strings.Repeat("u", 50) + "-sinkbinding",
	}, {
		name: "SinkBinding, already too long",
		vss: &v1alpha1.VSphereSource{
			ObjectMeta: metav1.ObjectMeta{
				Name: strings.Repeat("u", 63),
			},
		},
		f:    SinkBinding,
		want: "uuuuuuuuuuuuuuuuuuuca47ad1ce8479df271ec0d23653ce256-sinkbinding",
	}, {
		name: "SinkBinding",
		vss: &v1alpha1.VSphereSource{
			ObjectMeta: metav1.ObjectMeta{
				Name: "foo",
			},
		},
		f:    SinkBinding,
		want: "foo-sinkbinding",
	}, {
		name: "vspherebinding",
		vss: &v1alpha1.VSphereSource{
			ObjectMeta: metav1.ObjectMeta{
				Name: "baz",
			},
		},
		f:    VSphereBinding,
		want: "baz-vspherebinding",
	}, {
		name: "configmap",
		vss: &v1alpha1.VSphereSource{
			ObjectMeta: metav1.ObjectMeta{
				Name: "baz",
			},
		},
		f:    ConfigMap,
		want: "baz-configmap",
	}, {
		name: "rolebinding",
		vss: &v1alpha1.VSphereSource{
			ObjectMeta: metav1.ObjectMeta{
				Name: "baz",
			},
		},
		f:    RoleBinding,
		want: "baz-rolebinding",
	}, {
		name: "serviceaccount",
		vss: &v1alpha1.VSphereSource{
			ObjectMeta: metav1.ObjectMeta{
				Name: "baz",
			},
		},
		f:    ServiceAccount,
		want: "baz-serviceaccount",
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.f(test.vss)
			if got != test.want {
				t.Errorf("%s() = %v, wanted %v", test.name, got, test.want)
			}
		})
	}
}
