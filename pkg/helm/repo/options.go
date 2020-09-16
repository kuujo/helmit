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

func NewOptions(opts ...Option) Options {
	options := Options{}
	for _, opt := range opts {
		opt(&options)
	}
	return options
}

type Options struct {
	Name string
}

type Option func(options *Options)

func WithOptions(options Options) Option {
	return func(o *Options) {
		o = &options
	}
}

type AddOptions struct {
	URL      string
	CaFile   string
	KeyFile  string
	CertFile string
	Username string
	Password string
}

type AddOption func(options *AddOptions)

func WithAddOptions(options AddOptions) AddOption {
	return func(o *AddOptions) {
		o = &options
	}
}

func WithURL(url string) AddOption {
	return func(o *AddOptions) {
		o.URL = url
	}
}

func WithCaFile(caFile string) AddOption {
	return func(o *AddOptions) {
		o.CaFile = caFile
	}
}

func WithCertFile(certFile string) AddOption {
	return func(o *AddOptions) {
		o.CertFile = certFile
	}
}

func WithKeyFile(keyFile string) AddOption {
	return func(o *AddOptions) {
		o.KeyFile = keyFile
	}
}

func WithUsername(username string) AddOption {
	return func(o *AddOptions) {
		o.Username = username
	}
}

func WithPassword(password string) AddOption {
	return func(o *AddOptions) {
		o.Password = password
	}
}

type RemoveOptions struct{}

type RemoveOption func(options *RemoveOptions)

func WithRemoveOptions(options RemoveOptions) RemoveOption {
	return func(o *RemoveOptions) {
		o = &options
	}
}
