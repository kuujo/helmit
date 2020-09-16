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

package repo

import (
	"context"
	"github.com/gofrs/flock"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/helmpath"
	"helm.sh/helm/v3/pkg/repo"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var settings = cli.New()

// Client is a Helm repository client interface
type Client interface {
	Repo(name string, opts ...Option) *Repository
	Repos() []*Repository
}

func New(name string, options Options) (*Repository, error) {
	return &Repository{
		Name:      name,
		Options:   options,
		repoFile:  settings.RepositoryConfig,
		repoCache: settings.RepositoryCache,
	}, nil
}

// Repository is a Helm chart repository
type Repository struct {
	Options
	Name      string
	repoFile  string
	repoCache string
	mu        sync.RWMutex
}

func (r *Repository) Add(opts ...AddOption) error {
	config := AddOptions{}
	for _, opt := range opts {
		opt(&config)
	}

	err := os.MkdirAll(filepath.Dir(r.repoFile), os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return err
	}

	// Acquire a file lock for process synchronization
	fileLock := flock.New(strings.Replace(r.repoFile, filepath.Ext(r.repoFile), ".lock", 1))
	lockCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	locked, err := fileLock.TryLockContext(lockCtx, time.Second)
	if err == nil && locked {
		defer fileLock.Unlock()
	}
	if err != nil {
		return err
	}

	b, err := ioutil.ReadFile(r.repoFile)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	var f repo.File
	if err := yaml.Unmarshal(b, &f); err != nil {
		return err
	}

	if f.Has(r.Name) {
		return errors.Errorf("repository %q already exists", r.Name)
	}

	e := repo.Entry{
		Name:     r.Name,
		URL:      config.URL,
		Username: config.Username,
		Password: config.Password,
		CertFile: config.CertFile,
		KeyFile:  config.KeyFile,
		CAFile:   config.CaFile,
	}

	cr, err := repo.NewChartRepository(&e, getter.All(settings))
	if err != nil {
		return err
	}

	if _, err := cr.DownloadIndexFile(); err != nil {
		return errors.Wrapf(err, "%q is not a valid chart repository or cannot be reached", config.URL)
	}

	f.Update(&e)

	if err := f.WriteFile(r.repoFile, 0644); err != nil {
		return err
	}
	return nil
}

func (r *Repository) Remove(opts ...RemoveOption) error {
	cr, err := repo.LoadFile(r.repoFile)
	if os.IsNotExist(errors.Cause(err)) || len(cr.Repositories) == 0 {
		return errors.New("no repositories configured")
	}

	if !cr.Remove(r.Name) {
		return errors.Errorf("no repo named %q found", r.Name)
	}
	if err := cr.WriteFile(r.repoFile, 0644); err != nil {
		return err
	}

	idx := filepath.Join(r.repoCache, helmpath.CacheChartsFile(r.Name))
	if _, err := os.Stat(idx); err == nil {
		os.Remove(idx)
	}

	idx = filepath.Join(r.repoCache, helmpath.CacheIndexFile(r.Name))
	if _, err := os.Stat(idx); os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return errors.Wrapf(err, "cannot remove index file %s", idx)
	}
	return os.Remove(idx)
}
