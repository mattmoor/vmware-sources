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

// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/mattmoor/vmware-sources/pkg/apis/sources/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// VSphereBindingLister helps list VSphereBindings.
type VSphereBindingLister interface {
	// List lists all VSphereBindings in the indexer.
	List(selector labels.Selector) (ret []*v1alpha1.VSphereBinding, err error)
	// VSphereBindings returns an object that can list and get VSphereBindings.
	VSphereBindings(namespace string) VSphereBindingNamespaceLister
	VSphereBindingListerExpansion
}

// vSphereBindingLister implements the VSphereBindingLister interface.
type vSphereBindingLister struct {
	indexer cache.Indexer
}

// NewVSphereBindingLister returns a new VSphereBindingLister.
func NewVSphereBindingLister(indexer cache.Indexer) VSphereBindingLister {
	return &vSphereBindingLister{indexer: indexer}
}

// List lists all VSphereBindings in the indexer.
func (s *vSphereBindingLister) List(selector labels.Selector) (ret []*v1alpha1.VSphereBinding, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.VSphereBinding))
	})
	return ret, err
}

// VSphereBindings returns an object that can list and get VSphereBindings.
func (s *vSphereBindingLister) VSphereBindings(namespace string) VSphereBindingNamespaceLister {
	return vSphereBindingNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// VSphereBindingNamespaceLister helps list and get VSphereBindings.
type VSphereBindingNamespaceLister interface {
	// List lists all VSphereBindings in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1alpha1.VSphereBinding, err error)
	// Get retrieves the VSphereBinding from the indexer for a given namespace and name.
	Get(name string) (*v1alpha1.VSphereBinding, error)
	VSphereBindingNamespaceListerExpansion
}

// vSphereBindingNamespaceLister implements the VSphereBindingNamespaceLister
// interface.
type vSphereBindingNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all VSphereBindings in the indexer for a given namespace.
func (s vSphereBindingNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.VSphereBinding, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.VSphereBinding))
	})
	return ret, err
}

// Get retrieves the VSphereBinding from the indexer for a given namespace and name.
func (s vSphereBindingNamespaceLister) Get(name string) (*v1alpha1.VSphereBinding, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("vspherebinding"), name)
	}
	return obj.(*v1alpha1.VSphereBinding), nil
}
