// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

// Code generated by lister-gen. DO NOT EDIT.

package v2

import (
	ciliumiov2 "github.com/cilium/cilium/pkg/k8s/apis/cilium.io/v2"
	labels "k8s.io/apimachinery/pkg/labels"
	listers "k8s.io/client-go/listers"
	cache "k8s.io/client-go/tools/cache"
)

// CiliumBGPClusterConfigLister helps list CiliumBGPClusterConfigs.
// All objects returned here must be treated as read-only.
type CiliumBGPClusterConfigLister interface {
	// List lists all CiliumBGPClusterConfigs in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*ciliumiov2.CiliumBGPClusterConfig, err error)
	// Get retrieves the CiliumBGPClusterConfig from the index for a given name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*ciliumiov2.CiliumBGPClusterConfig, error)
	CiliumBGPClusterConfigListerExpansion
}

// ciliumBGPClusterConfigLister implements the CiliumBGPClusterConfigLister interface.
type ciliumBGPClusterConfigLister struct {
	listers.ResourceIndexer[*ciliumiov2.CiliumBGPClusterConfig]
}

// NewCiliumBGPClusterConfigLister returns a new CiliumBGPClusterConfigLister.
func NewCiliumBGPClusterConfigLister(indexer cache.Indexer) CiliumBGPClusterConfigLister {
	return &ciliumBGPClusterConfigLister{listers.New[*ciliumiov2.CiliumBGPClusterConfig](indexer, ciliumiov2.Resource("ciliumbgpclusterconfig"))}
}
