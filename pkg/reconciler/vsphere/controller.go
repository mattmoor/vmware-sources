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

	"knative.dev/pkg/tracker"

	"knative.dev/pkg/configmap"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/logging"

	vsphereinformer "github.com/mattmoor/vmware-sources/pkg/client/injection/informers/sources/v1alpha1/vspheresource"
	vspherereconciler "github.com/mattmoor/vmware-sources/pkg/client/injection/reconciler/sources/v1alpha1/vspheresource"
)

// NewController creates a Reconciler and returns the result of NewImpl.
func NewController(
	ctx context.Context,
	cmw configmap.Watcher,
) *controller.Impl {
	logger := logging.FromContext(ctx)

	vsphereInformer := vsphereinformer.Get(ctx)

	r := &Reconciler{}
	impl := vspherereconciler.NewImpl(ctx, r)
	r.Tracker = tracker.New(impl.EnqueueKey, controller.GetTrackerLease(ctx))

	logger.Info("Setting up event handlers.")

	vsphereInformer.Informer().AddEventHandler(controller.HandleAll(impl.Enqueue))

	return impl
}
