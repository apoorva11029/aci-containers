/***
Copyright 2019 Cisco Systems Inc. All rights reserved.

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

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	acisnatv1 "github.com/noironetworks/aci-containers/pkg/rdconfig/apis/aci.snat/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeRdConfigs implements RdConfigInterface
type FakeRdConfigs struct {
	Fake *FakeAciV1
	ns   string
}

var rdconfigsResource = schema.GroupVersionResource{Group: "aci.snat", Version: "v1", Resource: "rdconfigs"}

var rdconfigsKind = schema.GroupVersionKind{Group: "aci.snat", Version: "v1", Kind: "RdConfig"}

// Get takes name of the rdConfig, and returns the corresponding rdConfig object, and an error if there is any.
func (c *FakeRdConfigs) Get(name string, options v1.GetOptions) (result *acisnatv1.RdConfig, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(rdconfigsResource, c.ns, name), &acisnatv1.RdConfig{})

	if obj == nil {
		return nil, err
	}
	return obj.(*acisnatv1.RdConfig), err
}

// List takes label and field selectors, and returns the list of RdConfigs that match those selectors.
func (c *FakeRdConfigs) List(opts v1.ListOptions) (result *acisnatv1.RdConfigList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(rdconfigsResource, rdconfigsKind, c.ns, opts), &acisnatv1.RdConfigList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &acisnatv1.RdConfigList{ListMeta: obj.(*acisnatv1.RdConfigList).ListMeta}
	for _, item := range obj.(*acisnatv1.RdConfigList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested rdConfigs.
func (c *FakeRdConfigs) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(rdconfigsResource, c.ns, opts))

}

// Create takes the representation of a rdConfig and creates it.  Returns the server's representation of the rdConfig, and an error, if there is any.
func (c *FakeRdConfigs) Create(rdConfig *acisnatv1.RdConfig) (result *acisnatv1.RdConfig, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(rdconfigsResource, c.ns, rdConfig), &acisnatv1.RdConfig{})

	if obj == nil {
		return nil, err
	}
	return obj.(*acisnatv1.RdConfig), err
}

// Update takes the representation of a rdConfig and updates it. Returns the server's representation of the rdConfig, and an error, if there is any.
func (c *FakeRdConfigs) Update(rdConfig *acisnatv1.RdConfig) (result *acisnatv1.RdConfig, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(rdconfigsResource, c.ns, rdConfig), &acisnatv1.RdConfig{})

	if obj == nil {
		return nil, err
	}
	return obj.(*acisnatv1.RdConfig), err
}

// Delete takes name of the rdConfig and deletes it. Returns an error if one occurs.
func (c *FakeRdConfigs) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(rdconfigsResource, c.ns, name), &acisnatv1.RdConfig{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeRdConfigs) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(rdconfigsResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &acisnatv1.RdConfigList{})
	return err
}

// Patch applies the patch and returns the patched rdConfig.
func (c *FakeRdConfigs) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *acisnatv1.RdConfig, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(rdconfigsResource, c.ns, name, pt, data, subresources...), &acisnatv1.RdConfig{})

	if obj == nil {
		return nil, err
	}
	return obj.(*acisnatv1.RdConfig), err
}