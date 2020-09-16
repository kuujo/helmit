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

package release

import (
	"github.com/onosproject/helmit/pkg/kubernetes"
	"helm.sh/helm/v3/pkg/kube"
)

// Client is a release client
type Client interface {
	kubernetes.Client
}

// newClient creates a new release client
func newClient(namespace string, resources kube.ResourceList) (Client, error) {
	client, err := kubernetes.NewForResources(namespace, resources)
	if err != nil {
		return nil, err
	}
	return &kubernetesClient{
		Client: client,
	}, nil
}

// kubernetesClient is the default release client
type kubernetesClient struct {
	Client
}
