/*
Copyright 2019 The Knative Authors

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

package vsphere

import (
	"context"

	sourcesv1alpha1 "github.com/mattmoor/vmware-sources/pkg/apis/sources/v1alpha1"
	vspherereconciler "github.com/mattmoor/vmware-sources/pkg/client/injection/reconciler/sources/v1alpha1/vspheresource"
	"knative.dev/pkg/reconciler"
	"knative.dev/pkg/tracker"
)

// Reconciler implements vspherereconciler.Interface for
// VSphereSource resources.
type Reconciler struct {
	// Tracker builds an index of what resources are watching other resources
	// so that we can immediately react to changes tracked resources.
	Tracker tracker.Interface
}

// Check that our Reconciler implements Interface
var _ vspherereconciler.Interface = (*Reconciler)(nil)

// ReconcileKind implements Interface.ReconcileKind.
func (r *Reconciler) ReconcileKind(ctx context.Context, o *sourcesv1alpha1.VSphereSource) reconciler.Event {
	o.Status.InitializeConditions()

	// TODO(mattmoor): This

	o.Status.ObservedGeneration = o.Generation
	return nil
}
