/*
Copyright The Kubernetes Authors.

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

package poddefaultnetworknone

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/version"
	"k8s.io/component-helpers/nodedeclaredfeatures/types"
)

const (
	FeatureGate = "PodDefaultNetwork"
	FeatureName = "PodDefaultNetworkNone"
)

var Feature = &feature{}

var _ types.Feature = (*feature)(nil)

type feature struct{}

func (*feature) Name() string { return FeatureName }

func (*feature) Requirements() *types.FeatureRequirements {
	return &types.FeatureRequirements{
		EnabledFeatureGates: []string{FeatureGate},
		RequiredRuntimeFeatures: &types.RuntimeFeatures{
			DefaultNetworkNone: true,
		},
	}
}

func (*feature) Discover(cfg *types.NodeConfiguration) bool {
	return cfg.FeatureGates.Enabled(FeatureGate) && cfg.RuntimeFeatures.DefaultNetworkNone
}

func (*feature) InferForScheduling(podInfo *types.PodInfo) bool {
	return podInfo != nil && podInfo.Spec != nil && podInfo.Spec.DefaultNetwork != nil &&
		*podInfo.Spec.DefaultNetwork == v1.PodDefaultNetworkNone
}

func (*feature) InferForUpdate(_, _ *types.PodInfo) bool { return false }

func (*feature) MaxVersion() *version.Version { return nil }
