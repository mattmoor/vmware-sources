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

package v1alpha1

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	duckv1 "knative.dev/pkg/apis/duck/v1"
)

func TestVSphereSourceDefaulting(t *testing.T) {
	tests := []struct {
		name string
		c    *VSphereSource
		want *VSphereSource
	}{{
		name: "no change",
		c: &VSphereSource{
			ObjectMeta: metav1.ObjectMeta{
				Name: "valid",
			},
			Spec: VSphereSourceSpec{
				SourceSpec: validSourceSpec,
				VAuthSpec:  validVAuthSpec,
			},
		},
		want: &VSphereSource{
			ObjectMeta: metav1.ObjectMeta{
				Name: "valid",
			},
			Spec: VSphereSourceSpec{
				SourceSpec: validSourceSpec,
				VAuthSpec:  validVAuthSpec,
			},
		},
	}, {
		name: "ref gets namespace",
		c: &VSphereSource{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "valid",
				Namespace: "with-namespace",
			},
			Spec: VSphereSourceSpec{
				SourceSpec: duckv1.SourceSpec{
					Sink: duckv1.Destination{
						Ref: &duckv1.KReference{
							APIVersion: "serving.knative.dev",
							Kind:       "Service",
							Name:       "no-namespace",
						},
					},
				},
				VAuthSpec: validVAuthSpec,
			},
		},
		want: &VSphereSource{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "valid",
				Namespace: "with-namespace",
			},
			Spec: VSphereSourceSpec{
				SourceSpec: duckv1.SourceSpec{
					Sink: duckv1.Destination{
						Ref: &duckv1.KReference{
							APIVersion: "serving.knative.dev",
							Kind:       "Service",
							Namespace:  "with-namespace",
							Name:       "no-namespace",
						},
					},
				},
				VAuthSpec: validVAuthSpec,
			},
		},
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.c.DeepCopy()
			got.SetDefaults(context.Background())
			if !cmp.Equal(test.want, got) {
				t.Errorf("SetDefaults (-want, +got) = %v", cmp.Diff(test.want, got))
			}
		})
	}
}
