//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright 2023.

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

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	"k8s.io/api/core/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TuffMongoDB) DeepCopyInto(out *TuffMongoDB) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TuffMongoDB.
func (in *TuffMongoDB) DeepCopy() *TuffMongoDB {
	if in == nil {
		return nil
	}
	out := new(TuffMongoDB)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *TuffMongoDB) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TuffMongoDBList) DeepCopyInto(out *TuffMongoDBList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]TuffMongoDB, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TuffMongoDBList.
func (in *TuffMongoDBList) DeepCopy() *TuffMongoDBList {
	if in == nil {
		return nil
	}
	out := new(TuffMongoDBList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *TuffMongoDBList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TuffMongoDBSpec) DeepCopyInto(out *TuffMongoDBSpec) {
	*out = *in
	if in.MongoPorts != nil {
		in, out := &in.MongoPorts, &out.MongoPorts
		*out = make([]v1.ContainerPort, len(*in))
		copy(*out, *in)
	}
	if in.MongoVolumeMounts != nil {
		in, out := &in.MongoVolumeMounts, &out.MongoVolumeMounts
		*out = make([]v1.VolumeMount, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.MongoVolumes != nil {
		in, out := &in.MongoVolumes, &out.MongoVolumes
		*out = make([]v1.Volume, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TuffMongoDBSpec.
func (in *TuffMongoDBSpec) DeepCopy() *TuffMongoDBSpec {
	if in == nil {
		return nil
	}
	out := new(TuffMongoDBSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TuffMongoDBStatus) DeepCopyInto(out *TuffMongoDBStatus) {
	*out = *in
	if in.MongoPodNames != nil {
		in, out := &in.MongoPodNames, &out.MongoPodNames
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TuffMongoDBStatus.
func (in *TuffMongoDBStatus) DeepCopy() *TuffMongoDBStatus {
	if in == nil {
		return nil
	}
	out := new(TuffMongoDBStatus)
	in.DeepCopyInto(out)
	return out
}
