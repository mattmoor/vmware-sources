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
	"os"

	"k8s.io/client-go/tools/cache"
	"knative.dev/pkg/configmap"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/logging"

	"github.com/mattmoor/vmware-sources/pkg/apis/sources/v1alpha1"
	vsphereinformer "github.com/mattmoor/vmware-sources/pkg/client/injection/informers/sources/v1alpha1/vspheresource"
	vspherereconciler "github.com/mattmoor/vmware-sources/pkg/client/injection/reconciler/sources/v1alpha1/vspheresource"
	eventingclient "knative.dev/eventing/pkg/client/injection/client"
	sinkbindinginformer "knative.dev/eventing/pkg/client/injection/informers/sources/v1alpha1/sinkbinding"
	kubeclient "knative.dev/pkg/client/injection/kube/client"
	deploymentinformer "knative.dev/pkg/client/injection/kube/informers/apps/v1/deployment"
	cminformer "knative.dev/pkg/client/injection/kube/informers/core/v1/configmap"
	rbacinformer "knative.dev/pkg/client/injection/kube/informers/rbac/v1/rolebinding"
)

// NewController creates a Reconciler and returns the result of NewImpl.
func NewController(
	ctx context.Context,
	cmw configmap.Watcher,
) *controller.Impl {
	logger := logging.FromContext(ctx)

	vsphereInformer := vsphereinformer.Get(ctx)
	deploymentInformer := deploymentinformer.Get(ctx)
	sinkbindingInformer := sinkbindinginformer.Get(ctx)
	rbacInformer := rbacinformer.Get(ctx)
	cmInformer := cminformer.Get(ctx)

	r := &Reconciler{
		adapterImage:      os.Getenv("VSPHERE_ADAPTER"),
		kubeclient:        kubeclient.Get(ctx),
		eventingclient:    eventingclient.Get(ctx),
		deploymentLister:  deploymentInformer.Lister(),
		sinkbindingLister: sinkbindingInformer.Lister(),
		cmLister:          cmInformer.Lister(),
		rbacLister:        rbacInformer.Lister(),
	}
	impl := vspherereconciler.NewImpl(ctx, r)

	logger.Info("Setting up event handlers.")

	vsphereInformer.Informer().AddEventHandler(controller.HandleAll(impl.Enqueue))

	deploymentInformer.Informer().AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: controller.FilterGroupKind(v1alpha1.Kind("VSphereSource")),
		Handler:    controller.HandleAll(impl.EnqueueControllerOf),
	})

	sinkbindingInformer.Informer().AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: controller.FilterGroupKind(v1alpha1.Kind("VSphereSource")),
		Handler:    controller.HandleAll(impl.EnqueueControllerOf),
	})

	return impl
}
