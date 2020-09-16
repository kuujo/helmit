// Copyright 2020-present Open Networking Foundation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package test

import (
	"testing"
	"time"

	"github.com/onosproject/helmit/pkg/helm"
	"github.com/onosproject/helmit/pkg/test"
	"github.com/stretchr/testify/assert"
)

// ChartTestSuite is a test for chart deployment
type ChartTestSuite struct {
	test.Suite
}

// TestLocalInstall tests a local chart installation
func (s *ChartTestSuite) TestLocalInstall(t *testing.T) {
	_, err := helm.Install("atomix-controller", "atomix-controller").
		Set("scope", "Namespace").
		Wait().
		Do()
	assert.NoError(t, err)

	_, err = helm.Install("raft-storage-controller", "raft-storage-controller").
		Set("scope", "Namespace").
		Wait().
		Do()
	assert.NoError(t, err)

	topo, err := helm.Install("onos-topo", "onos-topo").
		Set("store.controller", "atomix-controller-kubernetes-controller:5679").
		Wait().
		Do()
	assert.NoError(t, err)

	pods, err := topo.Client().CoreV1().Pods().List()
	assert.NoError(t, err)
	assert.Len(t, pods, 2)

	deployment, err := topo.Client().AppsV1().
		Deployments().
		Get("onos-topo")
	assert.NoError(t, err)

	pods, err = deployment.Pods().List()
	assert.NoError(t, err)
	assert.Len(t, pods, 1)
	pod := pods[0]
	err = pod.Delete()
	assert.NoError(t, err)

	err = deployment.Wait(1 * time.Minute)
	assert.NoError(t, err)

	pods, err = deployment.Pods().List()
	assert.NoError(t, err)
	assert.Len(t, pods, 1)
	assert.NotEqual(t, pod.Name, pods[0].Name)

	services, err := topo.Client().CoreV1().Services().List()
	assert.NoError(t, err)
	assert.Len(t, services, 2)

	err = helm.Uninstall("atomix-controller").Do()
	assert.NoError(t, err)

	err = helm.Uninstall("raft-storage-controller").Do()
	assert.NoError(t, err)

	err = helm.Uninstall("onos-topo").Do()
	assert.NoError(t, err)
}

// TestRemoteInstall tests a remote chart installation
func (s *ChartTestSuite) TestRemoteInstall(t *testing.T) {
	err := helm.Repos().
		Add("atomix").
		URL("https://charts.atomix.io").
		Do()
	assert.NoError(t, err)

	_, err = helm.Install("atomix-controller", "atomix/atomix-controller").
		Set("scope", "Namespace").
		Wait().
		Do()
	assert.NoError(t, err)

	_, err = helm.Install("raft-storage-controller", "raft-storage-controller").
		Set("scope", "Namespace").
		Wait().
		Do()
	assert.NoError(t, err)

	topo, err := helm.Install("onos-topo", "onos-topo").
		Set("store.controller", "atomix-controller-kubernetes-controller:5679").
		Wait().
		Do()
	assert.NoError(t, err)

	pods, err := topo.Client().CoreV1().Pods().List()
	assert.NoError(t, err)
	assert.Len(t, pods, 2)

	deployment, err := topo.Client().AppsV1().
		Deployments().
		Get("onos-topo")
	assert.NoError(t, err)

	pods, err = deployment.Pods().List()
	assert.NoError(t, err)
	assert.Len(t, pods, 1)
	pod := pods[0]
	err = pod.Delete()
	assert.NoError(t, err)

	err = deployment.Wait(1 * time.Minute)
	assert.NoError(t, err)

	pods, err = deployment.Pods().List()
	assert.NoError(t, err)
	assert.Len(t, pods, 1)
	assert.NotEqual(t, pod.Name, pods[0].Name)

	services, err := topo.Client().CoreV1().Services().List()
	assert.NoError(t, err)
	assert.Len(t, services, 2)

	err = helm.Uninstall("atomix-controller").Do()
	assert.NoError(t, err)

	err = helm.Uninstall("raft-storage-controller").Do()
	assert.NoError(t, err)

	err = helm.Uninstall("onos-topo").Do()
	assert.NoError(t, err)
}
