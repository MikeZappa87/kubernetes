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
	"testing"

	v1 "k8s.io/api/core/v1"
	"k8s.io/component-helpers/nodedeclaredfeatures/types"
)

func TestFeature(t *testing.T) {
	for _, tc := range []struct {
		name        string
		gateEnabled bool
		runtime     bool
		wantFeature bool
		podMode     *v1.PodDefaultNetwork
		wantRequire bool
	}{
		{
			name:        "gate and runtime support",
			gateEnabled: true,
			runtime:     true,
			wantFeature: true,
			podMode:     networkMode(v1.PodDefaultNetworkNone),
			wantRequire: true,
		},
		{
			name:        "runtime support without gate",
			runtime:     true,
			podMode:     networkMode(v1.PodDefaultNetworkNone),
			wantRequire: true,
		},
		{
			name:        "gate without runtime support",
			gateEnabled: true,
			podMode:     networkMode(v1.PodDefaultNetworkNone),
			wantRequire: true,
		},
		{
			name:    "default pod network does not require feature",
			runtime: true,
			podMode: networkMode(v1.PodDefaultNetworkPod),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cfg := &types.NodeConfiguration{
				FeatureGates: types.FeatureGateMap{FeatureGate: tc.gateEnabled},
				RuntimeFeatures: types.RuntimeFeatures{
					DefaultNetworkNone: tc.runtime,
				},
			}
			if got := Feature.Discover(cfg); got != tc.wantFeature {
				t.Errorf("Discover() = %t, want %t", got, tc.wantFeature)
			}
			pod := &types.PodInfo{Spec: &v1.PodSpec{DefaultNetwork: tc.podMode}}
			if got := Feature.InferForScheduling(pod); got != tc.wantRequire {
				t.Errorf("InferForScheduling() = %t, want %t", got, tc.wantRequire)
			}
		})
	}
}

func networkMode(mode v1.PodDefaultNetwork) *v1.PodDefaultNetwork { return &mode }
