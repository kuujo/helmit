// Copyright 2019-present Open Networking Foundation.
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

package helm

import (
	"github.com/onosproject/helmit/pkg/helm/release"
	"github.com/onosproject/helmit/pkg/helm/repo"
	k8sconfig "github.com/onosproject/helmit/pkg/kubernetes/config"
)

// Namespace returns the Helm namespace
func Namespace() string {
	return k8sconfig.GetNamespaceFromEnv()
}

// NewRepo creates a new Helm chart repository
func NewRepo(name string, opts ...repo.Option) (*repo.Repository, error) {
	return repo.New(name, repo.NewOptions(opts...))
}

// NewRelease creates a new Helm chart release
func NewRelease(name string, chart string, opts ...release.Option) (*release.Release, error) {
	return release.New(name, release.NewOptions(chart, opts...))
}
