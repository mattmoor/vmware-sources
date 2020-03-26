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

package vsphere

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	cloudevents "github.com/cloudevents/sdk-go/v1"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/event"
	"github.com/vmware/govmomi/vim25/types"
	"go.uber.org/zap"

	"k8s.io/client-go/kubernetes"
	"knative.dev/eventing/pkg/adapter"
	"knative.dev/pkg/kvstore"
	"knative.dev/pkg/logging"
	"knative.dev/pkg/source"

	sourcesv1alpha1 "github.com/mattmoor/vmware-sources/pkg/apis/sources/v1alpha1"
	"github.com/mattmoor/vmware-sources/pkg/client/injection/client"
)

var groupResource = sourcesv1alpha1.Resource("vspheresources")

type envConfig struct {
	adapter.EnvConfig
}

func NewEnvConfig() adapter.EnvConfigAccessor {
	return &envConfig{}
}

// vAdapter implements the vSphereSource adapter to trigger a Sink.
type vAdapter struct {
	Logger    *zap.SugaredLogger
	Namespace string
	Source    string
	VClient   *govmomi.Client
	CEClient  cloudevents.Client
	Reporter  source.StatsReporter
	KVStore   kvstore.Interface
}

func NewAdapter(ctx context.Context, processed adapter.EnvConfigAccessor, ceClient cloudevents.Client, reporter source.StatsReporter) adapter.Adapter {
	env := processed.(*envConfig)

	logger := logging.FromContext(ctx)

	vClient, err := New(ctx)
	if err != nil {
		logger.Fatalf("Unable to create vSphere client: %v", err)
	}

	source, err := Address(ctx)
	if err != nil {
		logger.Fatalf("Unable to determine source: %v", err)
	}

	kubeclient := ctx.Value(client.Key{})
	if kubeclient == nil {
		logger.Fatalf("NIL kubeclient")
	}

	// TODO: configmap name needs to be passed in the env, or the name of the source needs to
	// plumbed in so we can derive it here.
	// https://github.com/mattmoor/vmware-sources/issues/8
	kvstore := kvstore.NewConfigMapKVStore(ctx, "test", env.Namespace, kubeclient.(*kubernetes.Clientset).CoreV1())
	err = kvstore.Init(ctx)
	if err != nil {
		logger.Fatalf("couldn't initialize kv store: %v", err)
	}

	return &vAdapter{
		Logger:    logger,
		Namespace: env.Namespace,
		Source:    source,
		Reporter:  reporter,
		VClient:   vClient,
		CEClient:  ceClient,
		KVStore:   kvstore,
	}
}

// Start implements adapter.Adapter
func (a *vAdapter) Start(stopCh <-chan struct{}) error {
	ctx, cancel := context.WithCancel(context.Background())
	// Cancel the context when the stop channel closes.
	go func() {
		<-stopCh
		cancel()
	}()
	// Below here use ctx.Done() instead of stopCh.

	manager := event.NewManager(a.VClient.Client)

	managedTypes := []types.ManagedObjectReference{a.VClient.ServiceContent.RootFolder}
	return manager.Events(ctx, managedTypes, 1, true /* tail */, false /* force */, a.sendEvents(ctx))
}

func (a *vAdapter) sendEvents(ctx context.Context) func(moref types.ManagedObjectReference, baseEvents []types.BaseEvent) error {
	return func(moref types.ManagedObjectReference, baseEvents []types.BaseEvent) error {
		for _, be := range baseEvents {
			event := cloudevents.NewEvent(cloudevents.VersionV1)

			event.SetType("com.vmware.vsphere." + strings.ToLower(reflect.TypeOf(be).Elem().Name()))
			event.SetTime(be.GetEvent().CreatedTime)
			event.SetID(fmt.Sprintf("%d", be.GetEvent().Key))
			event.SetSource(a.Source)

			switch e := be.(type) {
			case *types.EventEx:
				event.SetExtension("EventEx", e)
			case *types.ExtendedEvent:
				event.SetExtension("ExtendedEvent", e)
			}

			// TODO(mattmoor): Consider setting the subject

			// TODO(mattmoor): Switch to XML when sockeye stops sucking at it.
			event.SetDataContentType(cloudevents.ApplicationJSON)
			event.SetData(be)

			rctx, _, err := a.CEClient.Send(ctx, event)
			rtctx := cloudevents.HTTPTransportContextFrom(rctx)
			if err != nil {
				a.Logger.Error("failed to send cloudevent", zap.Error(err))
				return err
			}

			a.Reporter.ReportEventCount(&source.ReportArgs{
				Namespace:     a.Namespace,
				EventSource:   event.Source(),
				EventType:     event.Type(),
				ResourceGroup: groupResource.String(),
			}, rtctx.StatusCode)
		}

		return nil
	}
}
