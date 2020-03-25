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
	"fmt"

	sourcesv1alpha1 "github.com/mattmoor/vmware-sources/pkg/apis/sources/v1alpha1"
	vspherereconciler "github.com/mattmoor/vmware-sources/pkg/client/injection/reconciler/sources/v1alpha1/vspheresource"
	"github.com/mattmoor/vmware-sources/pkg/reconciler/vsphere/resources"
	resourcenames "github.com/mattmoor/vmware-sources/pkg/reconciler/vsphere/resources/names"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes"
	appsv1listers "k8s.io/client-go/listers/apps/v1"
	"knative.dev/eventing/pkg/apis/duck"
	clientset "knative.dev/eventing/pkg/client/clientset/versioned"
	sourcesv1alpha1lister "knative.dev/eventing/pkg/client/listers/sources/v1alpha1"
	"knative.dev/pkg/logging"
	"knative.dev/pkg/reconciler"
)

// Reconciler implements vspherereconciler.Interface for
// VSphereSource resources.
type Reconciler struct {
	adapterImage string

	kubeclient     kubernetes.Interface
	eventingclient clientset.Interface

	deploymentLister  appsv1listers.DeploymentLister
	sinkbindingLister sourcesv1alpha1lister.SinkBindingLister
}

// Check that our Reconciler implements Interface
var _ vspherereconciler.Interface = (*Reconciler)(nil)

// ReconcileKind implements Interface.ReconcileKind.
func (r *Reconciler) ReconcileKind(ctx context.Context, vms *sourcesv1alpha1.VSphereSource) reconciler.Event {
	vms.Status.InitializeConditions()

	if err := r.reconcileSinkBinding(ctx, vms); err != nil {
		return err
	}
	if err := r.reconcileDeployment(ctx, vms); err != nil {
		return err
	}

	vms.Status.ObservedGeneration = vms.Generation
	return nil
}

func (r *Reconciler) reconcileSinkBinding(ctx context.Context, vms *sourcesv1alpha1.VSphereSource) error {
	ns := vms.Namespace
	sinkbindingName := resourcenames.SinkBindingName(vms)

	sinkbinding, err := r.sinkbindingLister.SinkBindings(ns).Get(sinkbindingName)
	if apierrs.IsNotFound(err) {
		sinkbinding := resources.MakeSinkBinding(ctx, vms)
		sinkbinding, err = r.eventingclient.SourcesV1alpha1().SinkBindings(ns).Create(sinkbinding)
		if err != nil {
			return fmt.Errorf("failed to create sinkbinding %q: %w", sinkbindingName, err)
		}
		logging.FromContext(ctx).Infof("Created sinkbinding %q", sinkbindingName)
	} else if err != nil {
		return fmt.Errorf("failed to get sinkbinding %q: %w", sinkbindingName, err)
	} else {
		// The sinkbinding exists, but make sure that it has the shape that we expect.
		desiredSinkBinding := resources.MakeSinkBinding(ctx, vms)
		sinkbinding = sinkbinding.DeepCopy()
		sinkbinding.Spec = desiredSinkBinding.Spec
		sinkbinding, err = r.eventingclient.SourcesV1alpha1().SinkBindings(ns).Update(sinkbinding)
		if err != nil {
			return fmt.Errorf("failed to create sinkbinding %q: %w", sinkbindingName, err)
		}
	}

	// TODO(mattmoor): Check IsReady

	return nil
}

func (r *Reconciler) reconcileDeployment(ctx context.Context, vms *sourcesv1alpha1.VSphereSource) error {
	ns := vms.Namespace
	deploymentName := resourcenames.DeploymentName(vms)

	deployment, err := r.deploymentLister.Deployments(ns).Get(deploymentName)
	if apierrs.IsNotFound(err) {
		deployment := resources.MakeDeployment(ctx, vms, r.adapterImage)
		deployment, err = r.kubeclient.AppsV1().Deployments(ns).Create(deployment)
		if err != nil {
			return fmt.Errorf("failed to create deployment %q: %w", deploymentName, err)
		}
		logging.FromContext(ctx).Infof("Created deployment %q", deploymentName)
	} else if err != nil {
		return fmt.Errorf("failed to get deployment %q: %w", deploymentName, err)
	} else {
		// The deployment exists, but make sure that it has the shape that we expect.
		desiredDeployment := resources.MakeDeployment(ctx, vms, r.adapterImage)
		deployment = deployment.DeepCopy()
		deployment.Spec = desiredDeployment.Spec
		deployment, err = r.kubeclient.AppsV1().Deployments(ns).Update(deployment)
		if err != nil {
			return fmt.Errorf("failed to create deployment %q: %w", deploymentName, err)
		}
	}

	if duck.DeploymentIsAvailable(&deployment.Status, false) {
		logging.FromContext(ctx).Infof("TODO: Update status")
	} else {
		logging.FromContext(ctx).Infof("TODO: Update status")
	}

	return nil
}
