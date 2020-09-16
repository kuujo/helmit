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
	"errors"
	"github.com/onosproject/helmit/pkg/helm/context"
	"github.com/onosproject/helmit/pkg/kubernetes"
	"github.com/onosproject/helmit/pkg/kubernetes/config"
	"github.com/onosproject/helmit/pkg/kubernetes/filter"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/cli/values"
	"helm.sh/helm/v3/pkg/downloader"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/storage/driver"
	"log"
	"os"
	"sync"
	"time"
)

var settings = cli.New()

var conf = &configs{
	configs: make(map[string]*action.Configuration),
}

type configs struct {
	configs map[string]*action.Configuration
	mu      sync.Mutex
}

func (c *configs) get(namespace string) (*action.Configuration, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	config, ok := c.configs[namespace]
	if !ok {
		config = &action.Configuration{}
		if err := config.Init(settings.RESTClientGetter(), namespace, "memory", log.Printf); err != nil {
			return nil, err
		}
		c.configs[namespace] = config
	}
	return config, nil
}

// NewClient returns a new release client
func NewClient(namespace string) (*Client, error) {
	config, err := conf.get(namespace)
	if err != nil {
		return nil, err
	}
	return &Client{
		namespace: namespace,
		config:    config,
	}, nil
}

// Client is the Helm release client
type Client struct {
	namespace string
	config    *action.Configuration
}

// Get gets a release
func (c *Client) Get(name string) (*Release, error) {
	list, err := c.config.Releases.List(func(r *release.Release) bool {
		return r.Namespace == c.namespace && r.Name == name
	})
	if err != nil {
		return nil, err
	} else if len(list) == 0 {
		return nil, errors.New("release not found")
	} else if len(list) > 1 {
		return nil, errors.New("release is ambiguous")
	}
	return c.getRelease(list[0])
}

// List lists releases
func (c *Client) List() ([]*Release, error) {
	list, err := c.config.Releases.List(func(r *release.Release) bool {
		return r.Namespace == c.namespace
	})
	if err != nil {
		return nil, err
	}

	releases := make([]*Release, len(list))
	for i, release := range list {
		r, err := c.getRelease(release)
		if err != nil {
			return nil, err
		}
		releases[i] = r
	}
	return releases, nil
}

func (c *Client) getRelease(release *release.Release) (*Release, error) {
	resources, err := c.config.KubeClient.Build(bytes.NewBufferString(release.Manifest), true)
	if err != nil {
		return nil, err
	}

	parent, err := kubernetes.NewForNamespace(c.namespace)
	if err != nil {
		return nil, err
	}

	client, err := kubernetes.NewFiltered(c.namespace, filter.Resources(parent, resources))
	if err != nil {
		return nil, err
	}

	return &Release{
		release: release,
		client:  client,
	}, nil
}

// Install installs a release
func (c *Client) Install(release string, chart string) *InstallRequest {
	return &InstallRequest{
		name:   release,
		chart:  chart,
		values: make(map[string]interface{}),
	}
}

// Uninstall uninstalls a release
func (c *Client) Uninstall(release string) *UninstallRequest {
	return &UninstallRequest{
		name: release,
	}
}

// Upgrade upgrades a release
func (c *Client) Upgrade(release string, chart string) *UpgradeRequest {
	return &UpgradeRequest{
		name:   release,
		chart:  chart,
		values: make(map[string]interface{}),
	}
}

// Rollback rolls back a release
func (c *Client) Rollback(release string) *RollbackRequest {
	return &RollbackRequest{
		name: release,
	}
}

// Release is a Helm release
type Release struct {
	release *release.Release
	client  kubernetes.Client
}

// Client returns the release client
func (r *Release) Client() kubernetes.Client {
	return r.client
}

// InstallRequest is a release install request
type InstallRequest struct {
	name                     string
	namespace                string
	chart                    string
	repo                     string
	caFile                   string
	keyFile                  string
	certFile                 string
	username                 string
	password                 string
	version                  string
	values                   map[string]interface{}
	skipCRDs                 bool
	includeCRDs              bool
	disableHooks             bool
	disableOpenAPIValidation bool
	dryRun                   bool
	replace                  bool
	atomic                   bool
	wait                     bool
	timeout                  time.Duration
}

func (r *InstallRequest) Namespace(namespace string) *InstallRequest {
	r.namespace = namespace
	return r
}

func (r *InstallRequest) CaFile(caFile string) *InstallRequest {
	r.caFile = caFile
	return r
}

func (r *InstallRequest) KeyFile(keyFile string) *InstallRequest {
	r.keyFile = keyFile
	return r
}

func (r *InstallRequest) CertFile(certFile string) *InstallRequest {
	r.certFile = certFile
	return r
}

func (r *InstallRequest) Username(username string) *InstallRequest {
	r.username = username
	return r
}

func (r *InstallRequest) Password(password string) *InstallRequest {
	r.password = password
	return r
}

func (r *InstallRequest) Repo(url string) *InstallRequest {
	r.repo = url
	return r
}

func (r *InstallRequest) Version(version string) *InstallRequest {
	r.version = version
	return r
}

func (r *InstallRequest) Set(path string, value interface{}) *InstallRequest {
	setValue(r.values, path, value)
	return r
}

func (r *InstallRequest) SkipCRDs() *InstallRequest {
	r.skipCRDs = true
	return r
}

func (r *InstallRequest) IncludeCRDs() *InstallRequest {
	r.includeCRDs = true
	return r
}

func (r *InstallRequest) DisableHooks() *InstallRequest {
	r.disableHooks = true
	return r
}

func (r *InstallRequest) DisableOpenAPIValidation() *InstallRequest {
	r.disableOpenAPIValidation = true
	return r
}

func (r *InstallRequest) DryRun() *InstallRequest {
	r.dryRun = true
	return r
}

func (r *InstallRequest) Replace() *InstallRequest {
	r.replace = true
	return r
}

func (r *InstallRequest) Atomic() *InstallRequest {
	r.atomic = true
	return r
}

func (r *InstallRequest) Wait() *InstallRequest {
	r.wait = true
	return r
}

func (r *InstallRequest) Timeout(timeout time.Duration) *InstallRequest {
	r.timeout = timeout
	return r
}

func (r *InstallRequest) Do() error {
	namespace := r.namespace
	if namespace == "" {
		namespace = config.GetNamespaceFromEnv()
	}

	configuration, err := conf.get(namespace)
	if err != nil {
		return err
	}

	install := action.NewInstall(configuration)

	// Setup the repo options
	install.RepoURL = r.repo
	install.Username = r.username
	install.Password = r.password
	install.CaFile = r.caFile
	install.KeyFile = r.keyFile
	install.CertFile = r.certFile

	// Setup the chart options
	install.Version = r.version

	// Setup the release options
	install.ReleaseName = r.name
	install.Namespace = namespace
	install.Atomic = r.atomic
	install.Replace = r.replace
	install.DryRun = r.dryRun
	install.DisableHooks = r.disableHooks
	install.DisableOpenAPIValidation = r.disableOpenAPIValidation
	install.SkipCRDs = r.skipCRDs
	install.IncludeCRDs = r.includeCRDs
	install.Wait = r.wait
	install.Timeout = r.timeout

	// Locate the chart path
	path, err := install.ChartPathOptions.LocateChart(r.chart, settings)
	if err != nil {
		return err
	}

	// Check chart dependencies to make sure all are present in /charts
	chart, err := loader.Load(path)
	if err != nil {
		return err
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
					return err
				}
			} else {
				return err
			}
		}
	}

	ctx := context.GetContext().Release(r.name)
	valuesOptions := &values.Options{
		ValueFiles: ctx.ValueFiles,
		Values:     ctx.Values,
	}
	overrides, err := valuesOptions.MergeValues(getter.All(settings))
	if err != nil {
		return err
	}

	values := mergeValues(overrides, normalizeValues(r.values))
	_, err = install.Run(chart, values)
	return err
}

// UninstallRequest is a release uninstall request
type UninstallRequest struct {
	name      string
	namespace string
}

func (r *UninstallRequest) Namespace(namespace string) *UninstallRequest {
	r.namespace = namespace
	return r
}

func (r *UninstallRequest) Do() error {
	namespace := r.namespace
	if namespace == "" {
		namespace = config.GetNamespaceFromEnv()
	}

	configuration, err := conf.get(namespace)
	if err != nil {
		return err
	}

	uninstall := action.NewUninstall(configuration)
	_, err = uninstall.Run(r.name)
	return err
}

// UpgradeRequest is a release upgrade request
type UpgradeRequest struct {
	name         string
	namespace    string
	chart        string
	repo         string
	caFile       string
	keyFile      string
	certFile     string
	username     string
	password     string
	version      string
	values       map[string]interface{}
	disableHooks bool
	dryRun       bool
	atomic       bool
	wait         bool
	timeout      time.Duration
}

func (r *UpgradeRequest) Namespace(namespace string) *UpgradeRequest {
	r.namespace = namespace
	return r
}

func (r *UpgradeRequest) CaFile(caFile string) *UpgradeRequest {
	r.caFile = caFile
	return r
}

func (r *UpgradeRequest) KeyFile(keyFile string) *UpgradeRequest {
	r.keyFile = keyFile
	return r
}

func (r *UpgradeRequest) CertFile(certFile string) *UpgradeRequest {
	r.certFile = certFile
	return r
}

func (r *UpgradeRequest) Username(username string) *UpgradeRequest {
	r.username = username
	return r
}

func (r *UpgradeRequest) Password(password string) *UpgradeRequest {
	r.password = password
	return r
}

func (r *UpgradeRequest) Repo(url string) *UpgradeRequest {
	r.repo = url
	return r
}

func (r *UpgradeRequest) Version(version string) *UpgradeRequest {
	r.version = version
	return r
}

func (r *UpgradeRequest) Set(path string, value interface{}) *UpgradeRequest {
	setValue(r.values, path, value)
	return r
}

func (r *UpgradeRequest) DisableHooks() *UpgradeRequest {
	r.disableHooks = true
	return r
}

func (r *UpgradeRequest) DryRun() *UpgradeRequest {
	r.dryRun = true
	return r
}

func (r *UpgradeRequest) Atomic() *UpgradeRequest {
	r.atomic = true
	return r
}

func (r *UpgradeRequest) Wait() *UpgradeRequest {
	r.wait = true
	return r
}

func (r *UpgradeRequest) Timeout(timeout time.Duration) *UpgradeRequest {
	r.timeout = timeout
	return r
}

func (r *UpgradeRequest) Do() error {
	namespace := r.namespace
	if namespace == "" {
		namespace = config.GetNamespaceFromEnv()
	}

	configuration, err := conf.get(namespace)
	if err != nil {
		return err
	}

	upgrade := action.NewUpgrade(configuration)

	// Setup the repo options
	upgrade.RepoURL = r.repo
	upgrade.Username = r.username
	upgrade.Password = r.password
	upgrade.CaFile = r.caFile
	upgrade.KeyFile = r.keyFile
	upgrade.CertFile = r.certFile

	// Setup the chart options
	upgrade.Version = r.version

	// Setup the release options
	upgrade.Namespace = namespace
	upgrade.Atomic = r.atomic
	upgrade.DryRun = r.dryRun
	upgrade.DisableHooks = r.disableHooks
	upgrade.Wait = r.wait
	upgrade.Timeout = r.timeout

	// Locate the chart path
	path, err := upgrade.ChartPathOptions.LocateChart(r.chart, settings)
	if err != nil {
		return err
	}

	// Check chart dependencies to make sure all are present in /charts
	chart, err := loader.Load(path)
	if err != nil {
		return err
	}

	ctx := context.GetContext().Release(r.name)
	valuesOptions := &values.Options{
		ValueFiles: ctx.ValueFiles,
		Values:     ctx.Values,
	}
	overrides, err := valuesOptions.MergeValues(getter.All(settings))
	if err != nil {
		return err
	}

	values := mergeValues(overrides, normalizeValues(r.values))

	if upgrade.Install {
		// If a release does not exist, install it. If another error occurs during
		// the check, ignore the error and continue with the upgrade.
		histClient := action.NewHistory(configuration)
		histClient.Max = 1
		if _, err := histClient.Run(r.name); err == driver.ErrReleaseNotFound {
			install := action.NewInstall(configuration)
			install.ChartPathOptions = upgrade.ChartPathOptions
			install.DryRun = upgrade.DryRun
			install.DisableHooks = upgrade.DisableHooks
			install.Timeout = upgrade.Timeout
			install.Wait = upgrade.Wait
			install.Devel = upgrade.Devel
			install.Namespace = upgrade.Namespace
			install.Atomic = upgrade.Atomic
			install.PostRenderer = upgrade.PostRenderer

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
							return err
						}
					} else {
						return err
					}
				}
			}

			_, err = install.Run(chart, values)
			return err
		}
	}

	_, err = upgrade.Run(r.name, chart, values)
	return err
}

// RollbackRequest is a release rollback request
type RollbackRequest struct {
	name      string
	namespace string
}

func (r *RollbackRequest) Namespace(namespace string) *RollbackRequest {
	r.namespace = namespace
	return r
}

func (r *RollbackRequest) Do() error {
	namespace := r.namespace
	if namespace == "" {
		namespace = config.GetNamespaceFromEnv()
	}

	configuration, err := conf.get(namespace)
	if err != nil {
		return err
	}

	rollback := action.NewRollback(configuration)
	return rollback.Run(r.name)
}
