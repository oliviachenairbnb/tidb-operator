// Copyright PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"
	time "time"

	pingcapv1alpha1 "github.com/pingcap/tidb-operator/pkg/apis/federation/pingcap/v1alpha1"
	versioned "github.com/pingcap/tidb-operator/pkg/client/federation/clientset/versioned"
	internalinterfaces "github.com/pingcap/tidb-operator/pkg/client/federation/informers/externalversions/internalinterfaces"
	v1alpha1 "github.com/pingcap/tidb-operator/pkg/client/federation/listers/pingcap/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// VolumeBackupInformer provides access to a shared informer and lister for
// VolumeBackups.
type VolumeBackupInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.VolumeBackupLister
}

type volumeBackupInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewVolumeBackupInformer constructs a new informer for VolumeBackup type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewVolumeBackupInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredVolumeBackupInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredVolumeBackupInformer constructs a new informer for VolumeBackup type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredVolumeBackupInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.FederationV1alpha1().VolumeBackups(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.FederationV1alpha1().VolumeBackups(namespace).Watch(context.TODO(), options)
			},
		},
		&pingcapv1alpha1.VolumeBackup{},
		resyncPeriod,
		indexers,
	)
}

func (f *volumeBackupInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredVolumeBackupInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *volumeBackupInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&pingcapv1alpha1.VolumeBackup{}, f.defaultInformer)
}

func (f *volumeBackupInformer) Lister() v1alpha1.VolumeBackupLister {
	return v1alpha1.NewVolumeBackupLister(f.Informer().GetIndexer())
}