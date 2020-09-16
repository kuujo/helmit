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

package helm

import (
	"github.com/onosproject/helmit/pkg/helm/release"
	"github.com/onosproject/helmit/pkg/helm/repo"
	"github.com/onosproject/helmit/pkg/kubernetes/config"
)

// Namespace returns the Helm namespace
func Namespace() string {
	return config.GetNamespaceFromEnv()
}

// Repo returns the repository client
func Repo() *repo.Client {
	return &repo.Client{}
}

// Install installs a chart
func Install(name string, chart string) *release.InstallRequest {
	client, err := release.NewClient(Namespace())
	if err != nil {
		panic(err)
	}
	return client.Install(name, chart)
}

// Uninstall uninstalls a chart
func Uninstall(name string) *release.UninstallRequest {
	client, err := release.NewClient(Namespace())
	if err != nil {
		panic(err)
	}
	return client.Uninstall(name)
}

// Upgrade upgrades a release
func Upgrade(name string, chart string) *release.UpgradeRequest {
	client, err := release.NewClient(Namespace())
	if err != nil {
		panic(err)
	}
	return client.Upgrade(name, chart)
}

// Rollback rolls back a release
func Rollback(name string) *release.RollbackRequest {
	client, err := release.NewClient(Namespace())
	if err != nil {
		panic(err)
	}
	return client.Rollback(name)
}

// request is an interface for Helm requests
type request interface {
	// Do executes the request
	Do() error
}

var _ request = &repo.AddRequest{}

var _ request = &repo.RemoveRequest{}

var _ request = &release.InstallRequest{}

var _ request = &release.UninstallRequest{}

var _ request = &release.UpgradeRequest{}

var _ request = &release.RollbackRequest{}
