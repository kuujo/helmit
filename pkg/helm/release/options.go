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
	"time"
)

func NewOptions(chart string, opts ...Option) Options {
	options := Options{
		Chart: chart,
	}
	for _, opt := range opts {
		opt(&options)
	}
	return options
}

type Options struct {
	Chart     string
	Version   string
	Namespace string
	Values    map[string]interface{}
}

type Option func(*Options)

func WithNamespace(namespace string) Option {
	return func(o *Options) {
		o.Namespace = namespace
	}
}

func WithVersion(version string) Option {
	return func(options *Options) {
		options.Version = version
	}
}

func WithValue(path string, value interface{}) Option {
	return func(o *Options) {
		setValue(o.Values, path, value)
	}
}

func WithValues(values map[string]interface{}) Option {
	return func(o *Options) {
		o.Values = values
	}
}

type InstallOptions struct {
	RepoURL                  string
	CaFile                   string
	KeyFile                  string
	CertFile                 string
	Username                 string
	Password                 string
	SkipCRDs                 bool
	IncludeCRDs              bool
	DisableHooks             bool
	DisableOpenAPIValidation bool
	DryRun                   bool
	Replace                  bool
	Atomic                   bool
	Wait                     bool
	Timeout                  time.Duration
}

type InstallOption func(*InstallOptions)

func WithInstallOptions(options InstallOptions) InstallOption {
	return func(o *InstallOptions) {
		o = &options
	}
}

func WithRepoURL(url string) InstallOption {
	return func(o *InstallOptions) {
		o.RepoURL = url
	}
}

func WithCaFile(caFile string) InstallOption {
	return func(o *InstallOptions) {
		o.CaFile = caFile
	}
}

func WithCertFile(certFile string) InstallOption {
	return func(o *InstallOptions) {
		o.CertFile = certFile
	}
}

func WithKeyFile(keyFile string) InstallOption {
	return func(o *InstallOptions) {
		o.KeyFile = keyFile
	}
}

func WithUsername(username string) InstallOption {
	return func(o *InstallOptions) {
		o.Username = username
	}
}

func WithPassword(password string) InstallOption {
	return func(o *InstallOptions) {
		o.Password = password
	}
}

func WithSkipCRDs() InstallOption {
	return func(o *InstallOptions) {
		o.SkipCRDs = true
	}
}

func WithIncludeCRDs() InstallOption {
	return func(o *InstallOptions) {
		o.IncludeCRDs = true
	}
}

func WithDisableHooks() InstallOption {
	return func(o *InstallOptions) {
		o.DisableHooks = true
	}
}

func WithDisableOpenAPIValidation() InstallOption {
	return func(o *InstallOptions) {
		o.DisableOpenAPIValidation = true
	}
}

func WithDryRun() InstallOption {
	return func(o *InstallOptions) {
		o.DryRun = true
	}
}

func WithReplace() InstallOption {
	return func(o *InstallOptions) {
		o.Replace = true
	}
}

func WithAtomic() InstallOption {
	return func(o *InstallOptions) {
		o.Atomic = true
	}
}

func WithWait() InstallOption {
	return func(o *InstallOptions) {
		o.Wait = true
	}
}

func WithTimeout(timeout time.Duration) InstallOption {
	return func(o *InstallOptions) {
		o.Timeout = timeout
	}
}

type UninstallOptions struct{}

type UninstallOption func(*UninstallOptions)

func WithUninstallOptions(options UninstallOptions) UninstallOption {
	return func(o *UninstallOptions) {
		o = &options
	}
}
