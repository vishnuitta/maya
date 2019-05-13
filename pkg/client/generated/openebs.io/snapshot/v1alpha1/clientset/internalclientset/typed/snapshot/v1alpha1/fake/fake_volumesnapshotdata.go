/*
Copyright 2019 The OpenEBS Authors

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
	v1alpha1 "github.com/openebs/maya/pkg/apis/openebs.io/snapshot/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeVolumeSnapshotDatas implements VolumeSnapshotDataInterface
type FakeVolumeSnapshotDatas struct {
	Fake *FakeOpenebsV1alpha1
}

var volumesnapshotdatasResource = schema.GroupVersionResource{Group: "openebs.io", Version: "v1alpha1", Resource: "volumesnapshotdatas"}

var volumesnapshotdatasKind = schema.GroupVersionKind{Group: "openebs.io", Version: "v1alpha1", Kind: "VolumeSnapshotData"}

// Get takes name of the volumeSnapshotData, and returns the corresponding volumeSnapshotData object, and an error if there is any.
func (c *FakeVolumeSnapshotDatas) Get(name string, options v1.GetOptions) (result *v1alpha1.VolumeSnapshotData, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(volumesnapshotdatasResource, name), &v1alpha1.VolumeSnapshotData{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.VolumeSnapshotData), err
}

// List takes label and field selectors, and returns the list of VolumeSnapshotDatas that match those selectors.
func (c *FakeVolumeSnapshotDatas) List(opts v1.ListOptions) (result *v1alpha1.VolumeSnapshotDataList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(volumesnapshotdatasResource, volumesnapshotdatasKind, opts), &v1alpha1.VolumeSnapshotDataList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.VolumeSnapshotDataList{ListMeta: obj.(*v1alpha1.VolumeSnapshotDataList).ListMeta}
	for _, item := range obj.(*v1alpha1.VolumeSnapshotDataList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested volumeSnapshotDatas.
func (c *FakeVolumeSnapshotDatas) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(volumesnapshotdatasResource, opts))
}

// Create takes the representation of a volumeSnapshotData and creates it.  Returns the server's representation of the volumeSnapshotData, and an error, if there is any.
func (c *FakeVolumeSnapshotDatas) Create(volumeSnapshotData *v1alpha1.VolumeSnapshotData) (result *v1alpha1.VolumeSnapshotData, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(volumesnapshotdatasResource, volumeSnapshotData), &v1alpha1.VolumeSnapshotData{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.VolumeSnapshotData), err
}

// Update takes the representation of a volumeSnapshotData and updates it. Returns the server's representation of the volumeSnapshotData, and an error, if there is any.
func (c *FakeVolumeSnapshotDatas) Update(volumeSnapshotData *v1alpha1.VolumeSnapshotData) (result *v1alpha1.VolumeSnapshotData, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(volumesnapshotdatasResource, volumeSnapshotData), &v1alpha1.VolumeSnapshotData{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.VolumeSnapshotData), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeVolumeSnapshotDatas) UpdateStatus(volumeSnapshotData *v1alpha1.VolumeSnapshotData) (*v1alpha1.VolumeSnapshotData, error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateSubresourceAction(volumesnapshotdatasResource, "status", volumeSnapshotData), &v1alpha1.VolumeSnapshotData{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.VolumeSnapshotData), err
}

// Delete takes name of the volumeSnapshotData and deletes it. Returns an error if one occurs.
func (c *FakeVolumeSnapshotDatas) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(volumesnapshotdatasResource, name), &v1alpha1.VolumeSnapshotData{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeVolumeSnapshotDatas) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(volumesnapshotdatasResource, listOptions)

	_, err := c.Fake.Invokes(action, &v1alpha1.VolumeSnapshotDataList{})
	return err
}

// Patch applies the patch and returns the patched volumeSnapshotData.
func (c *FakeVolumeSnapshotDatas) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.VolumeSnapshotData, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(volumesnapshotdatasResource, name, pt, data, subresources...), &v1alpha1.VolumeSnapshotData{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.VolumeSnapshotData), err
}
