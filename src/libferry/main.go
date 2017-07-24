//
// Copyright © 2017 Ikey Doherty <ikey@solus-project.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package libferry

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"
)

const (
	// Version of the ferry client library
	Version = "0.0.1"
)

// A Client is used to communicate with the system ferryd
type Client struct {
	client *http.Client
}

// NewClient will return a new Client for the local unix socket, suitable
// for communicating with the daemon.
func NewClient(address string) *Client {
	return &Client{
		client: &http.Client{
			Transport: &http.Transport{
				DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
					return net.Dial("unix", address)
				},
				DisableKeepAlives:     false,
				IdleConnTimeout:       30 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
			},
			Timeout: 20 * time.Second,
		},
	}
}

// Close will kill any idle connections still in "keep-alive" and ensure we're
// not leaking file descriptors.
func (c *Client) Close() {
	trans := c.client.Transport.(*http.Transport)
	trans.CloseIdleConnections()
}

func (c *Client) formURI(part string) string {
	return fmt.Sprintf("http://localhost.localdomain:0/%s", part)
}

// GetVersion will return the version of the remote daemon
func (c *Client) GetVersion() (string, error) {
	var vq VersionRequest
	resp, e := c.client.Get(c.formURI("api/v1/version"))
	if e != nil {
		return "", e
	}
	defer resp.Body.Close()
	if e = json.NewDecoder(resp.Body).Decode(&vq); e != nil {
		return "", e
	}
	return vq.Version, nil
}

// A helper to wrap the trivial functionality, chaining off
// the appropriate errors, etc.
func (c *Client) getBasicResponse(url string, outT interface{}) error {
	resp, e := c.client.Get(url)
	if e != nil {
		return e
	}
	defer resp.Body.Close()
	if resp.ContentLength > 0 {
		if e = json.NewDecoder(resp.Body).Decode(outT); e != nil {
			return e
		}
	}
	fc := outT.(*Response)
	if !fc.Error {
		return nil
	}
	return errors.New(fc.ErrorString)
}

// CreateRepo will attempt to create a repository in the daemon
func (c *Client) CreateRepo(id string) error {
	uri := c.formURI("/api/v1/create_repo/" + id)
	return c.getBasicResponse(uri, &Response{})
}

// IndexRepo will attempt to index a repository in the daemon
func (c *Client) IndexRepo(id string) error {
	uri := c.formURI("/api/v1/index_repo/" + id)
	return c.getBasicResponse(uri, &Response{})
}
