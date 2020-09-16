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
	"bytes"
	"github.com/onosproject/helmit/pkg/helm/context"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/cli/values"
	"helm.sh/helm/v3/pkg/downloader"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/release"
	"log"
	"os"
)

var settings = cli.New()

// Installer is a release installer
type Installer interface {
	Install(opts ...InstallOption) (Client, error)
}

// Uninstaller is a release uninstaller
type Uninstaller interface {
	Uninstall(opts ...UninstallOption) error
}

// New creates a new release
func New(name string, options Options) (*Release, error) {
	return &Release{
		Name:    name,
		Options: options,
	}, nil
}

// Release is a Helm chart release
type Release struct {
	Options
	Name          string
	configuration *action.Configuration
	release       *release.Release
}

// getConfiguration gets the Helm configuration for the release
func (r *Release) getConfiguration() (*action.Configuration, error) {
	if r.configuration == nil {
		config := &action.Configuration{}
		if err := config.Init(settings.RESTClientGetter(), r.Namespace, "memory", log.Printf); err != nil {
			return nil, err
		}
		r.configuration = config
	}
	return r.configuration, nil
}

// getRelease loads the release from the Helm client
func (r *Release) getRelease() (*release.Release, error) {
	if r.release == nil {
		configuration, err := r.getConfiguration()
		if err != nil {
			return nil, err
		}
		release, err := configuration.Releases.Get(r.Name, 1)
		if err != nil {
			return nil, err
		}
		r.release = release
	}
	return r.release, nil
}

// Install installs the release
func (r *Release) Install(opts ...InstallOption) (Client, error) {
	options := InstallOptions{}
	for _, opt := range opts {
		opt(&options)
	}

	configuration, err := r.getConfiguration()
	if err != nil {
		return nil, err
	}

	install := action.NewInstall(configuration)

	// Setup the repo options
	install.RepoURL = options.RepoURL
	install.Username = options.Username
	install.Password = options.Password
	install.CaFile = options.CaFile
	install.KeyFile = options.KeyFile
	install.CertFile = options.CertFile

	// Setup the chart options
	install.Version = r.Version

	// Setup the release options
	install.ReleaseName = r.Name
	install.Namespace = r.Namespace
	install.Atomic = options.Atomic
	install.Replace = options.Replace
	install.DryRun = options.DryRun
	install.DisableHooks = options.DisableHooks
	install.DisableOpenAPIValidation = options.DisableOpenAPIValidation
	install.SkipCRDs = options.SkipCRDs
	install.IncludeCRDs = options.IncludeCRDs
	install.Wait = options.Wait
	install.Timeout = options.Timeout

	// Locate the chart path
	path, err := install.ChartPathOptions.LocateChart(r.Chart, settings)
	if err != nil {
		return nil, err
	}

	// Check chart dependencies to make sure all are present in /charts
	chart, err := loader.Load(path)
	if err != nil {
		return nil, err
	}

	if req := chart.Metadata.Dependencies; req != nil {
		// If CheckDependencies returns an error, we have unfulfilled dependencies.
		// As of Helm 2.4.0, this is treated as a stopping condition:
		// https://github.com/helm/helm/issues/2209
		if err := action.CheckDependencies(chart, req); err != nil {
			if install.DependencyUpdate {
				man := &downloader.Manager{
					Out:              os.Stdout,
					ChartPath:        path,
					Keyring:          install.ChartPathOptions.Keyring,
					SkipUpdate:       false,
					Getters:          getter.All(cli.New()),
					RepositoryConfig: settings.RepositoryConfig,
					RepositoryCache:  settings.RepositoryCache,
				}
				if err := man.Update(); err != nil {
					return nil, err
				}
			} else {
				return nil, err
			}
		}
	}

	ctx := context.GetContext().Release(r.Name)
	valuesOptions := &values.Options{
		ValueFiles: ctx.ValueFiles,
		Values:     ctx.Values,
	}
	overrides, err := valuesOptions.MergeValues(getter.All(settings))
	if err != nil {
		return nil, err
	}

	values := mergeValues(overrides, normalizeValues(r.Values))
	release, err := install.Run(chart, values)
	if err != nil {
		return nil, err
	}

	resources, err := configuration.KubeClient.Build(bytes.NewBufferString(release.Manifest), true)
	if err != nil {
		return nil, err
	}
	return newClient(release.Namespace, resources)
}

// Uninstall uninstalls the release
func (r *Release) Uninstall(opts ...UninstallOption) error {
	configuration, err := r.getConfiguration()
	if err != nil {
		return err
	}

	uninstall := action.NewUninstall(configuration)
	_, err = uninstall.Run(r.Name)
	return err
}

var _ Installer = &Release{}

var _ Uninstaller = &Release{}
